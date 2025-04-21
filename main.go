package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/Dekr0/shutil/cd"
	"github.com/Dekr0/shutil/config"
	"github.com/Dekr0/shutil/kitty"
	"github.com/Dekr0/shutil/pkg"
	"github.com/Dekr0/shutil/rename"
	"github.com/Dekr0/shutil/wezterm"
)

func main() {

    useWalker := flag.Bool(
        "walker",
        false,
        "Search for descendant (directory only) based on an array of directories",
    )
    useWalkerDepth := flag.Uint(
        "walker_depth",
        1,
        "Depth of searching descendant",
    )
    useWalkerWorker := flag.Uint(
        "walker_worker",
        0,
        "Worker of walker",
    )

    useReplaceSpace := flag.Bool(
        "replace_space",
        false,
        "Replace white space for an array of files or directories (including " +
        "descendant) with underscore",
    )

	useKittyFzFTab := flag.Bool(
		"kitty_fzf_tab",
		false,
		"Using FzF to search for the tabs in the current active window you want " +
		"to swap",
	)
    useWeztermFzFTab := flag.Bool(
        "wezterm_fzf_tab",
        false,
        "Using FzF to search for the tabs in the current active window you want " +
        "to swap",
    )

	useKittyStartSession := flag.Bool(
		"kitty_start_session",
		false,
		"Using FzF to search for session profile you want to use in the active windows",
	)
	useKittyStoreSession := flag.Bool(
		"kitty_store_session",
		false,
		"Store the current active sessions as a new profile",
	)

	useHealthCheck := flag.Bool(
		"health_check",
		false,
		"Check necessary applications required to run shutil, and check the correctness of config",
	)

	usePkgAdd := flag.String("pkg_add", "", "Package to be added from the profile")
	usePkgRm := flag.String("pkg_rm", "", "Package to be removed from the profile")
	usePkgCategory := flag.String("pkg_category", "other", "Package category")

    flag.Parse()

    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        AddSource: false,
        Level: slog.LevelInfo,
    }))

	if *useHealthCheck{
		healthCheck()
		os.Exit(0)
	}

	conf, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load configuration: %s", err.Error())
	}

    if *useWalker {
		roots := flag.Args()
		if len(roots) <= 0 {
			/* Look for bookmarks */
			if conf != nil {
				roots = conf.BookmarkRoots
			}
		}

        out, err := cd.SearchDir(
            roots,
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
		/* TODO: store this into .shutil */
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

	if *useKittyFzFTab {
		err := kitty.SwitchCurrentWindowTab()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}
    if *useWeztermFzFTab {
        err := wezterm.SwitchCurrentWindowTab()
        if err != nil {
            fmt.Println(err.Error())
            os.Exit(1)
        }
        os.Exit(0)
    }

	if *useKittyStartSession {
		if err := kitty.StartSession(conf); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	if *useKittyStoreSession {
		if err := kitty.StoreSession(conf); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

    flag.Usage()
}
