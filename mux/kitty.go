package mux

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"

	"github.com/valyala/fastjson"
)

type Kitty struct {
	Mux
}

func (k *Kitty) ActivateTab(blob []byte) error {
	splits := bytes.Split(blob, []byte(" -||- "))	
	if len(splits) != 2 {
	    return errors.New("Malformed FzF response")
	}
	return exec.Command(
		"kitten", "@", "focus-tab",
		"-m", fmt.Sprintf("id:%s", splits[0]),
	).Run()
}

func (k *Kitty) CurrentWindowTabs() ([][]byte, error) {
	blob, err := exec.Command("kitten", "@", "ls", "-t").Output()
	if err != nil {
		return nil, err
	}

	json, err := fastjson.Parse(string(blob))
	if err != nil {
		return nil, err
	}

	jWindows, err := json.Array()
	if err != nil {
		return nil, err
	}
	if len(jWindows) == 0 {
		return nil, errors.New("There's no active Kitty windows.")
	}

	jActiveWindow := jWindows[0]

	jTabs := jActiveWindow.GetArray("tabs")
	selections := make([][]byte, len(jTabs))
	for i, jTab := range jTabs {
		selections[i] = fmt.Appendf(
			nil,
			"%d -||- %s\n",
			jTab.GetInt("id"), jTab.GetStringBytes("title"),
		)
	}

	return selections, nil
}
