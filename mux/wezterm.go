package mux

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"

	"github.com/valyala/fastjson"
)

type Wezterm struct {
	Mux
}

func (w *Wezterm) ActivateTab(blob []byte) error {
	splits := bytes.Split(blob, []byte(" - "))
	if len(splits) != 2 {
	    return errors.New("Malformed FzF response")
	}
	return exec.Command(
		"wezterm", "cli", "activate-tab", 
		"--tab-id", string(splits[0]),
	).Run()
}

func (w *Wezterm) CurrentWindowTabs() ([][]byte, error) {
	p := exec.Command("wezterm", "cli", "list", "--format", "json")

    blob, err := p.Output()
    if err != nil {
        return nil, err
    }
    json, err := fastjson.Parse(string(blob))
    if err != nil {
        return nil, err
    }
    jTabs, err := json.Array()
    if err != nil {
        return nil, err
    }

    selections := make([][]byte, len(jTabs) - 1)
    for i, jTab := range jTabs {
		if i >= len(jTabs) - 1 {
			continue
		}
        tabId := jTab.GetUint("tab_id")
        tabTitle := jTab.GetStringBytes("tab_title")
        if bytes.Compare(tabTitle, []byte("")) == 0 {
            tabTitle = fmt.Appendf(nil, "Unnamed Tab %d", i)
        }
        selections[i] = fmt.Appendf(nil, "%d - %s\n", tabId, tabTitle)
    }

	return selections, nil
}

func (w *Wezterm) NewTab(blob []byte) (error) {
	return exec.Command(
		"wezterm", "cli", "spawn", 
		"--cwd", string(blob),
	).Run()
}
