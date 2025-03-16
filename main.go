package main

import (
	"flag"
	"fmt"
	"os"

	"dekr0.com/shutil/cd"
)

func main() {
    useWalkerRoot := flag.Bool(
        "walker-roots",
        false,
        "Search for descendant (directory only) based on an array of directories",
    )
    useWalkerRootDepth := flag.Uint(
        "walker-root-depth",
        1,
        "Depth of searching descendant",
    )
    
    flag.Parse()

    if *useWalkerRoot {
        out, err := cd.SearchDir(flag.Args(), uint8(*useWalkerRootDepth))
        if err != nil {
            fmt.Println(err.Error())
            os.Exit(1)
        }
        os.Stdout.Write(out)
        os.Exit(0)
    }
}
