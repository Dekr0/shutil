package cd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
)

func SearchDir(roots []string, depth uint8) ([]byte, error) {
    if len(roots) <= 0 {
        return nil, errors.New("shutil --walker-roots: no starting root directory is provided")
    }

    if depth <= 0 {
        return nil, errors.New("shutil --walker-root-depth: depth must be greater than 0")
    }

    type CmdResult struct {
        out []byte
        err error
    }

    type PathResult struct {
        path string
        err error
    }

    fzf := exec.Command("fzf")
    fzfInPipe, err := fzf.StdinPipe()
    if err != nil {
        return nil, err
    }

    cmdResultChan := make(chan *CmdResult)
    pathResultChan := make(chan *PathResult)

    go func() { 
        o, e := fzf.Output()
        cmdResultChan <- &CmdResult{ o, e }
    }()

    /**
     Since this is a deterministic operation, and if FzF finishes before this,
     this function will directly exit, and the application will exit immediately.
     No need to break this loop into units of work and handle clean up and 
     cancellation.
     */
    go func() {
        for _, root := range roots {
            paths := []string{ root }
            pathResultChan <- &PathResult{ root, nil }
            depthCount := depth
            for depthCount > 0 && len(paths) > 0 {
                currNumPaths := len(paths)
                for i := 0; i < currNumPaths; i++ {
                    entries, err := os.ReadDir(paths[i])
                    if err != nil {
                        if !errors.Is(err, os.ErrPermission) {
                            pathResultChan <- &PathResult{ "", err }
                        }
                    }
                    for _, entry := range entries {
                        if entry.IsDir() {
                            newPath := path.Join(paths[i], entry.Name())
                            pathResultChan <- &PathResult{ newPath, nil }
                            paths = append(paths, newPath)
                        }
                    }
                }
                paths = paths[currNumPaths:]
                depthCount--
            }
        }
    }()

    var cmdResult *CmdResult
    for {
        select {
        case cmdResult = <- cmdResultChan:
            return cmdResult.out, cmdResult.err
        case pathResult := <- pathResultChan:
            if pathResult.err != nil {
                return nil, pathResult.err
            }
            path := fmt.Sprintf("%s\n", pathResult.path)
            _, err := fzfInPipe.Write([]byte(path))
            if err != nil {
                return nil, err
            }
        }
    }
}
