package rename

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type empty struct{}

type walker struct {
	depth uint8
    ctx   context.Context
	cErr  chan error
	sem   chan empty
    wg    sync.WaitGroup
}

func (w *walker) walk(root string, depth uint8) (error) {
    if depth > w.depth {
        return nil
    }
    i, err := os.Lstat(root)
    if err != nil {
        return err
    }
    newRoot := strings.ReplaceAll(root, " ", "_")
    if err := os.Rename(root, newRoot); err != nil {
        return err
    }
    if !i.IsDir() {
        return nil
    }
    f, err := os.Open(newRoot)
    for {
        select {
        case <- w.ctx.Done():
            return w.ctx.Err()
        default:
            entries, err := f.Readdir(1024)
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
                case w.sem <- empty{}:
                    w.wg.Add(1)
                    go func() {
                        defer func() {
                            <- w.sem
                            w.wg.Done()
                        }()
                        w.walk(path.Join(newRoot, entry.Name()), depth + 1)
                    }()
                default:
                    w.walk(path.Join(newRoot, entry.Name()), depth + 1)
                }
            }
        }
    }
}

func ReplaceSpace(
    roots []string, depth uint8, workers uint8, deadline time.Duration, logger *slog.Logger,
) (error) {
    if len(roots) <= 0 {
        return nil
    }
    if depth <= 0 {
        return nil
    }
    if workers <= 0 {
        workers = 4
    }
    if deadline <= 0 {
        deadline = time.Second * 60
    }

    bgCtx := context.Background()
    wCtx, wCancel := context.WithDeadline(bgCtx, time.Now().Add(deadline))
    defer wCancel()

    cErr := make(chan error)

    w := walker {
        depth: depth,
        ctx: wCtx,
        cErr: cErr,
        sem: make(chan empty, workers),
    }

    for _, root := range roots {
        w.wg.Add(1)
        go func() {
            defer w.wg.Done()
            cErr <- w.walk(root, 0)
        }()
    }

    wgCtx, wgCancel := context.WithDeadline(bgCtx, time.Now().Add(deadline * 2))
    defer wgCancel()
    go func() {
        w.wg.Wait()
        wgCancel()
    }()

    for {
        select {
        case gErr := <- cErr:
            if gErr != nil {
                logger.Error("Received error", "error", gErr)
            }
        case <- wgCtx.Done():
            wgErr := wgCtx.Err()
            if wgErr != nil && errors.Is(wgErr, context.DeadlineExceeded) {
                panic("Wait timeout. It would be a deadlock")
            }
            return nil
        }
    }
}
