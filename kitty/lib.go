package kitty

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Dekr0/shutil/config"
	"github.com/valyala/fastjson"
)

func StartSession(conf *config.Config) error {
	if len(conf.SessionProfiles) <= 0 {
		return nil
	}
	
	var builder strings.Builder
	for i, conf := range conf.SessionProfiles {
		builder.WriteString(fmt.Sprintf("%d - %s\n", i, conf.Name))
	}

	fzf := exec.Command("fzf")
	pipe, err := fzf.StdinPipe()
	if err != nil {
		return err
	}
	_, err = pipe.Write([]byte(builder.String()))
	if err != nil {
		return err
	}
	res, err := fzf.Output()
	if err != nil {
		return err
	}
	idx, err := strconv.ParseUint(
		string(res[0:bytes.IndexFunc(res, func(r rune) bool { return r == ' ' })]),
		10, 8,
	)
	if err != nil {
		return err
	}
	for _, session := range conf.SessionProfiles[idx].Sessions {
		kitten := exec.Command("kitten", "@", "launch",
			"--type", "tab",
			"--tab-title", session.Name,
			"--cwd", os.ExpandEnv(session.Path),
		)
		_, err := kitten.Output()
		if err != nil {
			return err
		}
	}

	return nil
}

func StoreSession(conf *config.Config) (error) {
	kittenQueryProc := exec.Command("kitten", "@", "ls", "-t")

	blob, queryErr := kittenQueryProc.Output()
	if queryErr != nil {
		return queryErr
	}

	json, parseErr := fastjson.Parse(string(blob))
	if parseErr != nil {
		return parseErr
	}

	jWindows, err := json.Array()
	if err != nil {
		return err
	}
	if len(jWindows) == 0 {
		return errors.New("There's no active Kitty windows.")
	}

	jActiveWindow := jWindows[0]

	jTabs := jActiveWindow.GetArray("tabs")

	profile := config.SessionProfile{
		Name: "",
		Sessions: make([]config.Session, 0, 4),
	}
	for _, jTab := range jTabs {
		window := jTab.GetArray("windows")[0]
		profile.Sessions = append(profile.Sessions, config.Session{
			Name: string(jTab.GetStringBytes("title")),
			Path: string(window.GetStringBytes("cwd")),
		})
	}

	fmt.Print("Enter profile name: ")
	_, err = fmt.Scan(&profile.Name)
	if err != nil {
		return err
	}

	conf.SessionProfiles = append(conf.SessionProfiles, &profile)

	return conf.SaveConfig()
}
