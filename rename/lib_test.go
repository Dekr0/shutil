package rename

import (
	"log/slog"
	"os"
	"path"
	"testing"
	"time"
)

type tree struct {
    name string
    isDir bool
    nodes []tree
}

func generateTest() error {
    t := tree{
        name: "folder 1",
        isDir: true,
        nodes: []tree{
            { "test 1.txt", false, nil },
            { "test 2.txt", false, nil },
            { 
                "sub folder 1", 
                true,
                []tree{
                    { "sub test 1.txt", false, nil},
                    { "sub test 2.txt", false, nil},
                    { 
                        "sub sub folder 1", 
                        true,
                        []tree{
                            { "sub sub test 1.txt", false, nil },
                            { "sub sub test 2.txt", false, nil },
                        },
                    },
                },
            },
            { 
                "sub folder 2", 
                true,
                []tree{
                    { "sub test 3.txt", false, nil },
                    { "sub test 4.txt", false, nil },
                },
            },
        },
    }
    nodes := []tree{ t }
    for len(nodes) > 0 {
        top := nodes[0]
        nodes = nodes[1:]
        if !top.isDir {
            if err := os.WriteFile(top.name, []byte{'a'}, 0666); err != nil {
                return err
            }
            continue
        }
        if err := os.Mkdir(top.name, 0666); err != nil {
            return err
        }
        for _, node := range top.nodes {
            node.name = path.Join(top.name, node.name)
            nodes = append(nodes, node)
        }
    }
    return nil
}

func TestReplaceWhiteSpace(t *testing.T) {
    logger := slog.New(slog.NewJSONHandler(
        os.Stdout,
        &slog.HandlerOptions {
            AddSource: false,
            Level: slog.LevelInfo,
        },
    ))

    /*err := generateTest()
    if err != nil {
        t.Fatal(err)
    }*/

    err := ReplaceSpace(
        []string{ "D:/codebase/shutil/rename/folder_1" },
        4, 0, time.Second * 5, logger,
    )
    if err != nil {
        t.Fatal(err)
    }
}
