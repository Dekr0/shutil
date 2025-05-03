package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Dekr0/shutil/picker"
	"github.com/Dekr0/shutil/config"
	"github.com/Dekr0/shutil/mux"
	"github.com/Dekr0/shutil/pkg"
	"github.com/Dekr0/shutil/rename"
)

func main() {
	useFzFExec := flag.Bool(
		"cofi",
		false,
		"Fuzzy find executable from PATH",
	)

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

	useKittyActivateTab := flag.Bool(
		"kitty_activate_tab",
		false,
		"Using FzF to search for the tabs in the current active window you want " +
		"to swap",
	)
    useWeztermActivateTab := flag.Bool(
        "wezterm_activate_tab",
        false,
        "Using FzF to search for the tabs in the current active window you want " +
        "to swap",
    )

	useWeztermNewTab := flag.Bool(
		"wezterm_new_tab",
		false,
		"Launch a new wt tab using directories in book mark",
	)

	useKittyNewSessions := flag.Bool(
		"kitty_new_sessions",
		false,
		"Using FzF to search for session profile you want to use in the active" +
		" windows",
	)
	useWeztermNewSessions := flag.Bool(
		"wezterm_new_sessions",
		false,
		"Using FzF to search for session profile you want to use in the active" +
		" windows",
	)

	useKittyCreateSessionProfile := flag.Bool(
		"kitty_create_session_profile",
		false,
		"Store the current active sessions as a new profile",
	)
	useWeztermCreateSessionProfile := flag.Bool(
		"wezterm_create_session_profile",
		false,
		"Store the current active sessions as a new profile",
	)

	useHealthCheck := flag.Bool(
		"health_check",
		false,
		"Check necessary applications required to run shutil, and check the " + 
		"correctness of config",
	)

	usePkgAdd := flag.String("pkg_add", "", "Package to be added from the profile")
	usePkgRm := flag.String("pkg_rm", "", "Package to be removed from the profile")
	usePkgCategory := flag.String("pkg_category", "other", "Package category")

    flag.Parse()

	f, err := os.OpenFile("/tmp/shutil.log", os.O_CREATE | os.O_WRONLY, 0666)
	if err != nil {
		slog.Error("Failed to create shutil log files", "error", err)
	}
    logger := slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{
        AddSource: false,
        Level: slog.LevelInfo,
    }))
	slog.SetDefault(logger)

	ctx, cancel := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGABRT,
	)
	defer cancel()

	if *useHealthCheck{
		healthCheck()
		os.Exit(0)
	}

	c, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load configuration: %s", err.Error())
	}

	if *useFzFExec {
		out, err := picker.SearchExecutable(ctx)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		os.Stdout.Write(out)
		os.Exit(0)
	}

    if *useWalker {
		roots := flag.Args()
		if len(roots) <= 0 {
			/* Look for bookmarks */
			if c != nil {
				roots = c.BookmarkRoots
			}
		}

        out, err := picker.SearchDir(
			ctx,
            roots,
            uint8(*useWalkerDepth), uint8(*useWalkerWorker),
        )
        if err != nil {
			slog.Error(err.Error())
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
			slog.Error(err.Error())
            os.Exit(1)
        }
        os.Exit(0)
    }
	
	if len(*usePkgAdd) > 0 {
		/* TODO: store this into .shutil */
		err := pkg.AddPkg(*usePkgAdd, *usePkgCategory)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	if len(*usePkgRm) > 0 {
		err := pkg.RmPkg(*usePkgRm, *usePkgCategory)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	if *useKittyActivateTab {
		err := mux.ActivateTab(ctx, &mux.Kitty{})
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}
    if *useWeztermActivateTab {
        err := mux.ActivateTab(ctx, &mux.Wezterm{})
        if err != nil {
			slog.Error(err.Error())
            os.Exit(1)
        }
        os.Exit(0)
    }

	if *useWeztermNewTab {
		err := mux.NewTab(
			ctx,
			c.BookmarkRoots,
			uint8(*useWalkerDepth), uint8(*useWalkerWorker),
			&mux.Wezterm{},
		)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	if *useKittyNewSessions {
		if err := mux.NewSessions(c, &mux.Kitty{}); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}
	if *useWeztermNewSessions {
		if err := mux.NewSessions(c, &mux.Wezterm{}); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	if *useKittyCreateSessionProfile {
		if err := mux.CreateSessionProfile(c, &mux.Kitty{}); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}
	if *useWeztermCreateSessionProfile {
		if err := mux.CreateSessionProfile(c, &mux.Wezterm{}); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

    flag.Usage()
}
