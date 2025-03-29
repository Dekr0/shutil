package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/Dekr0/shutil/cd"
	"github.com/Dekr0/shutil/pkg"
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

	usePkgAdd := flag.String("pkg-add", "", "Package to be added from the profile")
	usePkgRm := flag.String("pkg-rm", "", "Package to be removed from the profile")
	usePkgCategory := flag.String("pkg-category", "other", "Package category")

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
	
	if len(*usePkgAdd) > 0 {
		err := pkg.AddPkg(*usePkgAdd, *usePkgCategory)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	if len(*usePkgRm) > 0 {
		err := pkg.RmPkg(*usePkgRm, *usePkgCategory)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

    flag.Usage()
}
