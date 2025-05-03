package picker

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"sync"
	"time"

	"github.com/Dekr0/shutil/fzf"
)

type empty struct{}

type walker struct {
    depth uint8
    ctx   context.Context
    d     chan string
    e     chan error
    sem   chan empty
    wg    sync.WaitGroup
}

func (w *walker) walk(root string, depth uint8) (error) {
    if depth > w.depth {
        return nil
    }
    info, err := os.Lstat(root)
    if err != nil {
        return err
    }
    if !info.IsDir() {
        return nil
    }
    f, err := os.Open(root)
    if err != nil {
        return err
    }
    for {
        select {
        case <- w.ctx.Done():
            return w.ctx.Err()
        default:
            entries, err := f.ReadDir(1024)
            if err != nil {
                if err == io.EOF {
                    return nil
                }
                return err
            }
            for _, entry := range entries {
                select {
                case <- w.ctx.Done():
                    return w.ctx.Err()
                default:
                    if !entry.IsDir() {
                        continue
                    }
                    p := path.Join(root, entry.Name())
                    w.d <- p
                    select {
                        case <- w.ctx.Done():
                            return w.ctx.Err()
                        case w.sem <- empty{}:
                            w.wg.Add(1)
                            go func() {
                                defer func() {
                                    <- w.sem
                                    w.wg.Done()
                                }()
                                if err := w.walk(p, depth + 1); err != nil {
                                    w.e <- err
                                }
                            }()
                        default:
                            if err := w.walk(p, depth + 1); err != nil {
                                w.e <- err
                            }
                    }
                }
            }
        }
    }
}

func SearchDir(
    ctx context.Context, roots []string, depth uint8, workers uint8,
) ([]byte, error) {
    if len(roots) == 0 {
        return nil, nil
    }
    if depth == 0 {
        return nil, nil
    }
    if workers == 0 {
        workers = 4
    }

	for i := range roots {
		roots[i] = os.ExpandEnv(roots[i])
	}

    proc := exec.CommandContext(ctx, "fzf")
    pipe, pErr := proc.StdinPipe()
    if pErr != nil {
        return nil, pErr
    }
    
    /*
     Non-buffered channel will block for send and receive until the other side 
     is ready.
     Buffered channel will block only when the buffer is full.
    */
    e := make(chan error)
    d := make(chan string)
	r := make(chan *fzf.FzFSimpleResponse)
    closed := false

    var wg sync.WaitGroup

    w := walker {
        depth: depth,
        ctx: ctx,
        d: d,
        e: e,
        sem: make(chan empty, workers),
    }

    go func() {
        o, e := proc.Output()
		r <- &fzf.FzFSimpleResponse{
			Selected: o, 
			Error: e,
		}
    }()

    for _, root := range roots {
        wg.Add(1)
        go func() {
            defer wg.Done()
            e <- w.walk(root, 0)
        }()
    }

    for {
        select {
        case err := <- e:
            if err != nil {
                slog.Error("Received error from a walker", "error", err)
            }
        case r := <- r:
            ctx, cancel := context.WithTimeout(context.Background(), time.Second * 4)
			defer cancel()
            go func() {
                for {
                    select {
                    case <- e:
                    case <- d:
                    case <- ctx.Done():
                        return
                    }
                }
            }()

            w.wg.Wait()
            wg.Wait()

            return r.Selected, r.Error
        case p := <- d:
            if closed {
                continue
            }
            p = fmt.Sprintf("%s\n", p)
            if _, err := pipe.Write([]byte(p)); err != nil {
                return nil, err
            }
        case <- ctx.Done():
            ctx, cancel := context.WithTimeout(context.Background(), time.Second * 4)
			defer cancel()
            
            /*
            Semaphore pattern:
            - sending an item into the channel => taking a semaphore
            - receiving an item into the channel => releasing a semaphore

            All walk go routine might block at sending either error message or 
            path.
            When deadline is reached, the main thread is no longer polling for 
            error message or path. 
            This will cause some walk go routine to be unable to detect 
            cancellation signal since it's block at send. 
            This go routine is here to discard the remaining error message or path 
            It will run until all go routine is canceled, mark as done for the 
            wait group, and return.

            Main thread should avoid taking / releasing semaphores because this 
            will cause go routines, which finish execution normally without 
            canceling, will block depends on main thread is taking semaphores, or 
            is releasing semaphores.
            
            Let say it's releasing semaphores.
            By releasing semaphores, the main thread is consuming items from the 
            semaphore channel.
            Since parts of the go routines already exits, either due to cancellation 
            or exiting normally, there are less go routines that will take the 
            semaphore.
            Meanwhile, the main thread is actively releasing semaphores.
            Also, there are go routines, which return, and are about to releasing 
            semaphores as well.
            This is will cause the semaphore channel to be emptied asymmetrically 
            given that main thread is actively consuming items from the channel. 
            If majority of the go routines returns, the go routines that are trying 
            to release semaphore will block since the semaphore channel is emptied 
            and no go routine is taking semaphore.
            */
            go func() {
                for {
                    select {
                    case <- e:
                    case <- d:
                    case <- ctx.Done():
                        return
                    }
                }
            }()
            w.wg.Wait()
            wg.Wait()
            return nil, ctx.Err()
        }
    }
}
