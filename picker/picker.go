package picker

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Dekr0/shutil/fzf"
	"golang.org/x/sys/unix"
)

func SearchExecutable(ctx context.Context) ([]byte, error) {
	env := os.Getenv("PATH")
	paths := strings.Split(env, ":")

	pipe, rc, err := fzf.SimplePickAsyncDual(ctx)
	if err != nil {
		return nil, nil
	}

	sem := make(chan struct{}, 100)
	e := make(chan error)
	for i := range paths {
		p := paths[i]
		select {
		case err := <- e:
			return nil, err
		case sem <- struct{}{}:
			go func() {
				if err := searchExecutable(ctx, pipe, p); err != nil {
					e <- err
				}
				<- sem
			}()
		default:
			if err := searchExecutable(ctx, pipe, p); err != nil {
				return nil, err
			}
		}
	}

	for {
		select {
		case <- ctx.Done():
			return nil, ctx.Err()
		case err := <- e:
			if err != nil {
				return nil, err
			}
		case r := <- rc:
			if r.Error != nil {
				return nil, err
			}
			
			selected := bytes.Trim(
				bytes.TrimSpace(r.Selected),
				"\n",
			)

			return []byte(selected), nil
		}
	}
}

func searchExecutable(ctx context.Context, pipe io.WriteCloser, path string) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for {
		select {
		case <- ctx.Done():
			return ctx.Err()
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
				case <- ctx.Done():
					return ctx.Err()
				default:
					if entry.IsDir() {
						continue
					}
					fullPath := filepath.Join(path, entry.Name())
					if runtime.GOOS == "windows" {
						// TODO
					} else {
						if err := unix.Access(fullPath, unix.X_OK); err != nil {
							continue
						}
						if _, err := pipe.Write([]byte(entry.Name() + "\n")); 
						err != nil {
							return err
						}
					}
				}
			}
		}
	}
}
