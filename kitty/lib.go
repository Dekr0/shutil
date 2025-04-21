package kitty

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Dekr0/shutil/config"
	"github.com/Dekr0/shutil/fzf"
	"github.com/valyala/fastjson"
)

func SwitchCurrentWindowTab() (error) {
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
	selections := make([]string, len(jTabs))
	/** Might malformed ?*/
	for i, jTab := range jTabs {
		selections[i] = fmt.Sprintf(
			"%d -||- %s\n", jTab.GetInt("id"), jTab.GetStringBytes("title"),
		)
	}

	fzfProc := exec.Command("fzf")
	pFzf, pErr := fzfProc.StdinPipe()
	if pErr != nil {
		return pErr
	}
	cFzfSimpleResponse := make(chan *fzf.FzFSimpleResponse)

	go func() {
		o, e := fzfProc.Output()
		cFzfSimpleResponse <- &fzf.FzFSimpleResponse{ Selection: o, Error: e }
	}()

	for _, selection := range selections {
		_, wErr := pFzf.Write([]byte(selection))
		if wErr != nil {
			return wErr
		}
	}

	for {
		select {
		case r := <- cFzfSimpleResponse:
			if r.Error != nil {
				return r.Error
			}

			splits := bytes.Split(r.Selection, []byte(" -||- "))	
			if len(splits) != 2 {
				return errors.New("Selection is malformed")
			}

			focusErr := exec.Command(
				"kitten", "@", "focus-tab", "-m", fmt.Sprintf("id:%s", splits[0]),
			).Run()
			if focusErr != nil {
				return focusErr
			}

			return nil
		default:
			time.Sleep(time.Microsecond * 100)
		}
	}
}

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
