package mux

import (
	"bytes"
	"time"

	"github.com/Dekr0/shutil/cd"
	"github.com/Dekr0/shutil/fzf"
)

type Mux interface {
	CurrentWindowTabs() ([][]byte, error)
	NewTab([]byte) error
	ActivateTab([]byte) error
}

func ActivateTab(mux Mux) error {
	selections, err := mux.CurrentWindowTabs()
	if err != nil {
		return err
	}

	c, err := fzf.SimplePick(selections)
	if err != nil {
		return err
	}

    for {
        select {
        case r := <- c:
            if r.Error != nil {
                return r.Error
            }
			return mux.ActivateTab(r.Selected)
        default:
            time.Sleep(time.Microsecond * 100)
        }
    }
}

func NewTab(roots []string, depth uint8, workers uint8, mux Mux) error {
	selected, err := cd.SearchDir(roots, depth, workers)
	if err != nil {
		return err
	}
	selected = bytes.TrimSpace(selected)
	selected = bytes.ReplaceAll(selected, []byte{'\n'}, []byte{})
	return mux.NewTab(selected)
}
