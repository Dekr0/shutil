package kitty

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"time"

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
