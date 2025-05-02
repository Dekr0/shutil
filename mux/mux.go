package mux

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Dekr0/shutil/cd"
	"github.com/Dekr0/shutil/config"
	"github.com/Dekr0/shutil/fzf"
)

type Mux interface {
	ActivateTab(string) error
	CurrentWindowTabs() ([][]byte, error)
	NewSessions(*config.SessionProfile) error
	NewTab(string) error
	SnapshotSessions() ([]*config.Session, error)
}

func ActivateTab(ctx context.Context, mux Mux) error {
	selections, err := mux.CurrentWindowTabs()
	if err != nil {
		return err
	}
	c, err := fzf.SimplePickAsync(ctx, selections)
	if err != nil {
		return err
	}

    for {
        select {
		case <- ctx.Done():
			return ctx.Err()
        case r := <- c:
            if r.Error != nil {
                return r.Error
            }
			return mux.ActivateTab(string(r.Selected))
        default:
            time.Sleep(time.Microsecond * 100)
        }
    }
}

func NewSessions(c *config.Config, mux Mux) error {
	p, err := PickSessionProfile(c)
	if err != nil {
		return err
	}
	return mux.NewSessions(p)
}

func PickSessionProfile(c *config.Config) (
	*config.SessionProfile,
	error,
) {
	if len(c.SessionProfiles) <= 0 {
		return nil, nil
	}
	
	var builder strings.Builder
	for i, p := range c.SessionProfiles {
		builder.WriteString(fmt.Sprintf("%d - %s\n", i, p.Name))
	}
	
	res, err := fzf.SimplePickSync([]byte(builder.String()))
	if err != nil {
		return nil, err
	}

	idx, err := strconv.ParseUint(
		string(res[0:bytes.IndexFunc(res, func(r rune) bool { return r == ' ' })]),
		10, 8,
	)
	if err != nil {
		return nil, err
	}

	return c.SessionProfiles[idx], nil
}

func NewTab(
	ctx context.Context,
	roots []string,
	depth uint8,
	workers uint8,
	mux Mux,
) error {
	selected, err := cd.SearchDir(ctx, roots, depth, workers)
	if err != nil {
		return err
	}
	selected = bytes.TrimSpace(selected)
	selected = bytes.ReplaceAll(selected, []byte{'\n'}, []byte{})
	return mux.NewTab(string(selected))
}

func CreateSessionProfile(c *config.Config, mux Mux) error {
	sessions, err := mux.SnapshotSessions()
	if err != nil {
		return err
	}
	p := config.SessionProfile{Name: "", Sessions: sessions}
	fmt.Print("Enter profile name: ")
	if _, err = fmt.Scan(&p.Name); err != nil {
		return err
	}
	c.SessionProfiles = append(c.SessionProfiles, &p)
	return c.SaveConfig()
}
