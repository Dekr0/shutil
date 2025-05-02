package mux

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Dekr0/shutil/config"
	"github.com/valyala/fastjson"
)

type Kitty struct {
	Mux
}

func (k *Kitty) ActivateTab(selected string) error {
	splits := strings.SplitN(selected, " - ", 2)	
	if len(splits) != 2 {
	    return errors.New("Malformed FzF response")
	}
	return exec.Command(
		"kitten", "@", "focus-tab",
		"-m", fmt.Sprintf("id:%s", splits[0]),
	).Run()
}

func (k *Kitty) CurrentWindowTabs() ([][]byte, error) {
	tabs, err := k.currentWindowTabs()
	if err != nil {
		return nil, err
	}

	selections := make([][]byte, len(tabs))
	for i, tab := range tabs {
		selections[i] = fmt.Appendf(
			nil,
			"%d - %s\n",
			tab.GetInt("id"), tab.GetStringBytes("title"),
		)
	}

	return selections, nil
}

func (k *Kitty) NewSessions(p *config.SessionProfile) error {
	for _, session := range p.Sessions {
		if err := k.NewTabWithTitle(session.Path, session.Name); err != nil {
			return err
		}
	}
	return nil
}

func (k *Kitty) NewTabWithTitle(pwd string, title string) error {
	if err := exec.Command("kitten", "@", "launch",
		"--type", "tab",
		"--tab-title", title,
		"--cwd", os.ExpandEnv(pwd),
	).Run(); err != nil {
		return err
	}
	return nil
}

func (k *Kitty) NewTab(path string) error {
	return nil
}

func (k *Kitty) SnapshotSessions() ([]*config.Session, error) {
	tabs, err := k.currentWindowTabs()
	if err != nil {
		return nil, err
	}
	sessions := make([]*config.Session, len(tabs))
	for i, tab := range tabs {
		window := tab.GetArray("windows")[0]
		sessions[i] = &config.Session{
			Name: string(tab.GetStringBytes("title")),
			Path: string(window.GetStringBytes("cwd")),
		}
	}
	return sessions, nil
}

func (k *Kitty) currentWindowTabs() ([]*fastjson.Value, error) {
	listing, err := exec.Command("kitten", "@", "ls", "-t").Output()
	if err != nil {
		return nil, err
	}

	json, err := fastjson.Parse(string(listing))
	if err != nil {
		return nil, err
	}

	windows, err := json.Array()
	if err != nil {
		return nil, err
	}
	if len(windows) == 0 {
		return nil, errors.New("There's no active Kitty windows.")
	}

	return windows[0].GetArray("tabs"), nil
}
