package cd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"sync"
	"time"
)

type empty struct{}

type walker struct {
    depth uint8
    ctx   context.Context
    cDir  chan string
    cErr  chan error
    logger *slog.Logger
    sem   chan empty
    wg    sync.WaitGroup
}

func (w *walker) walk(root string, depth uint8) (error) {
    if depth > w.depth {
        return nil
    }
    info, lErr := os.Lstat(root)
    if lErr != nil {
        return lErr
    }
    if !info.IsDir() {
        return nil
    }
    f, oErr := os.Open(root)
    if oErr != nil {
        return oErr
    }
    for {
        select {
        case <- w.ctx.Done():
            return w.ctx.Err()
        default:
            entries, rErr := f.ReadDir(1024)
            if rErr != nil {
                if rErr == io.EOF {
                    return nil
                }
                return rErr
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
                    w.cDir <- p
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
                                if wErr := w.walk(p, depth + 1); wErr != nil {
                                    w.cErr <- wErr
                                }
                            }()
                        default:
                            if wErr := w.walk(p, depth + 1); wErr != nil {
                                w.cErr <- wErr
                            }
                    }
                }
            }
        }
    }
}

type rFzf struct {
    out []byte
    err error
}

func SearchDir(
    roots []string, depth uint8, workers uint8, logger *slog.Logger,
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
    bg := context.Background()

	for i := range roots {
		roots[i] = os.ExpandEnv(roots[i])
	}

    wCtx, wCancel := context.WithCancel(bg)
    defer wCancel()

    fzfCtx, fzfCancel := context.WithCancel(bg)
    defer fzfCancel()

    fzf := exec.CommandContext(fzfCtx, "fzf")
    pFzf, pErr := fzf.StdinPipe()
    if pErr != nil {
        return nil, pErr
    }
    
    /*
     Non-buffered channel will block for send and receive until the other side 
     is ready.
     Buffered channel will block only when the buffer is full.
    */
    cErr := make(chan error)
    cDir := make(chan string)
    cFzf := make(chan *rFzf)
    closed := false

    var wg sync.WaitGroup

    w := walker {
        depth: depth,
        ctx: wCtx,
        cDir: cDir,
        cErr: cErr,
        logger: logger,
        sem: make(chan empty, workers),
    }

    go func() {
        o, e := fzf.Output()
        cFzf <- &rFzf{ o, e }
    }()

    for _, root := range roots {
        wg.Add(1)
        go func() {
            defer wg.Done()
            cErr <- w.walk(root, 0)
        }()
    }

    deadline := time.After(time.Second * 60)
    for {
        select {
        case gErr := <- cErr:
            if gErr != nil {
                logger.Error("Received error", "error", gErr)
            }
        case r := <- cFzf:
            wCancel()

            cCtx, cCanel := context.WithCancel(bg)
            defer cCanel()
            go func() {
                for {
                    select {
                    case <- cErr:
                    case <- cDir:
                    case <- cCtx.Done():
                        return
                    }
                }
            }()

            w.wg.Wait()
            wg.Wait()

            return r.out, r.err
        case p := <- cDir:
            if closed {
                continue
            }
            p = fmt.Sprintf("%s\n", p)
            if _, err := pFzf.Write([]byte(p)); err != nil {
                return nil, err
            }
        case <- deadline:
            fzfCancel()
            wCancel()
            err := errors.New("Exceed deadline.")
            cCtx, cCanel := context.WithCancel(bg)
            defer cCanel()
            
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
                    case <- cErr:
                    case <- cDir:
                    case <- cCtx.Done():
                        return
                    }
                }
            }()
            w.wg.Wait()
            wg.Wait()
            return nil, err
        }
    }
}
