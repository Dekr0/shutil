package main

import (
	"fmt"
	"os/exec"

	"github.com/Dekr0/shutil/config"
)

func healthCheck() {
	_, err := config.LoadConfig()
	fmt.Printf("Configuration file:")
	if err != nil {
		fmt.Printf("\n  Failed to load configuration: %s\n", err.Error())
	} else {
		fmt.Println(" OK")
	}

	which := exec.Command("which", "fzf")
	_, err = which.Output()
	fmt.Printf("FzF:")
	if err != nil {
		fmt.Print("\n  fzf not found\n")
	} else {
		fmt.Printf(" OK\n")
	}

	which = exec.Command("which", "kitty")
	_, err = which.Output()
	fmt.Printf("Kitty (optional):")
	if err != nil {
		fmt.Print("\n  kitty not found\n")
	} else {
		fmt.Print(" OK\n")
	}

	which = exec.Command("which", "wezterm")
	_, err = which.Output()
	fmt.Printf("Wezterm (optional):")
	if err != nil {
		fmt.Print("\n  wezterm not found\n")
	} else {
		fmt.Print("\n  OK")
	}
}
