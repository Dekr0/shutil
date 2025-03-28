package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/Dekr0/shutil/cd"
	"github.com/Dekr0/shutil/rename"
)

func main() {
    useWalker := flag.Bool(
        "walker",
        false,
        "Search for descendant (directory only) based on an array of directories",
    )
    useWalkerDepth := flag.Uint(
        "walker-depth",
        1,
        "Depth of searching descendant",
    )
    useWalkerWorker := flag.Uint(
        "walker-worker",
        0,
        "Worker of walker",
    )
    useReplaceSpace := flag.Bool(
        "replace-space",
        false,
        "Replace white space for an array of files or directories (including " +
        "descendant) with underscore",
    )
    
    flag.Parse()

    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        AddSource: false,
        Level: slog.LevelInfo,
    }))

    if *useWalker {
        out, err := cd.SearchDir(
            flag.Args(),
            uint8(*useWalkerDepth), uint8(*useWalkerWorker), logger,
        )
        if err != nil {
            fmt.Println(err.Error())
            os.Exit(1)
        }
        os.Stdout.Write(out)
        os.Exit(0)
    }
    
    if *useReplaceSpace {
        err := rename.ReplaceSpace(
            flag.Args(),
            uint8(*useWalkerDepth), uint8(*useWalkerWorker), 0,
            logger,
        )
        if err != nil {
            fmt.Println(err.Error())
            os.Exit(1)
        }
        os.Exit(0)
    }

    flag.Usage()
}
