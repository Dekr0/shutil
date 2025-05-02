package mux

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/Dekr0/shutil/config"
	"github.com/valyala/fastjson"
)

type Wezterm struct {
	Mux
}

func (w *Wezterm) ActivateTab(selected string) error {
	splits := strings.SplitN(selected, " - ", 2)
	if len(splits) != 2 {
	    return errors.New("Malformed FzF response")
	}
	slog.Info(splits[0])
	return exec.Command(
		"wezterm", "cli", "activate-tab", 
		"--tab-id", splits[0],
	).Run()
}

func (w *Wezterm) CurrentWindowTabs() ([][]byte, error) {
	tabs, err := w.currentWindowTabs()
	if err != nil {
		return nil, err
	}
    selections := make([][]byte, len(tabs) - 1)
    for i, tab := range tabs {
		if i >= len(tabs) - 1 {
			continue
		}
        tabId := tab.GetUint("tab_id")
        tabTitle := tab.GetStringBytes("tab_title")
        if bytes.Compare(tabTitle, []byte("")) == 0 {
            tabTitle = fmt.Appendf(nil, "Unnamed Tab %d", i)
        }
        selections[i] = fmt.Appendf(nil, "%d - %s\n", tabId, tabTitle)
    }

	return selections, nil
}

func (w *Wezterm) NewSessions(p *config.SessionProfile) error {
	sem := make(chan struct{}, 4)
	var wg sync.WaitGroup
	for _, session := range p.Sessions {
		select {
		case sem <- struct{}{}:
			wg.Add(1)
			go func() {
				if err := w.NewTabWithTitle(session.Path, session.Name); err != nil {
					slog.Error(err.Error())
				}
				<- sem
				wg.Done()
			}()
		default:
			if err := w.NewTabWithTitle(session.Path, session.Name); err != nil {
				return err
			}
		}
	}
	wg.Wait()
	return nil
}

func (w *Wezterm) NewTabWithTitle(pwd string, title string) error {
	if output, err := exec.Command("wezterm", "cli", "spawn",
		"--cwd", os.ExpandEnv(pwd),
		).Output(); err != nil {
		return err
	} else {
		slog.Info(string(output))
		pane_id := strings.TrimSpace(string(output))
		pane_id = strings.Trim(pane_id, "\n")
		if err := exec.Command("wezterm", "cli", "set-tab-title",
			"--pane-id", pane_id, title,
			).Run(); err != nil {
			return err
		}
	}
	return nil
}

func (w *Wezterm) SnapshotSessions() ([] *config.Session, error) {
	tabs, err := w.currentWindowTabs()
	if err != nil {
		return nil, err
	}
	sessions := make([]*config.Session, len(tabs))
	for i, tab := range tabs {
		sessions[i] = &config.Session{
			Name: string(tab.GetStringBytes("tab_title")),
			Path: string(tab.GetStringBytes("cwd")),
		}
	}
	return sessions, nil
}

func (w *Wezterm) NewTab(path string) (error) {
	return exec.Command(
		"wezterm", "cli", "spawn", 
		"--cwd", path,
	).Run()
}

func (w *Wezterm) currentWindowTabs() ([]*fastjson.Value, error) {
	p := exec.Command("wezterm", "cli", "list", "--format", "json")

    listing, err := p.Output()
    if err != nil {
        return nil, err
    }
    json, err := fastjson.Parse(string(listing))
    if err != nil {
        return nil, err
    }
    return json.Array()
}
