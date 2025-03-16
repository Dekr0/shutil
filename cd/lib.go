package cd

import (
	"os"
	"path"
	"strings"
)

func SearchDir(roots []string, depth uint8) (string, error) {
    var builder strings.Builder

    for depth > 0 {
        stack := []string{}
        for _, root := range roots {
            children, err := os.ReadDir(root)
            if err != nil {
                return "", err
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

    return builder.String(), nil
}
