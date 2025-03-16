package cd

import (
	"os"
	"os/exec"
	"path"
	"strings"
)

func SearchDir(roots []string, depth uint8) ([]byte, error) {
    var builder strings.Builder

    for depth > 0 {
        stack := []string{}
        for _, root := range roots {
            children, err := os.ReadDir(root)
            if err != nil {
                return nil, err
            }
            for _, child := range children {
                if child.Type().IsDir() {
                    full_path := path.Join(root, child.Name())
                    stack = append(stack, full_path)
                    builder.WriteString(full_path)
                    builder.WriteByte('\n')
                }
            }
        }
        depth--
    }

    fzf := exec.Command("fzf")

    in, err := fzf.StdinPipe()
    if err != nil {
        return nil, err
    }

    _, err = in.Write([]byte(builder.String()))
    if err != nil {
        return nil, err
    }

    pick, err := fzf.Output()
    if err != nil {
        return nil, err
    }

    return pick, nil
}
