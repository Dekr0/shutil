package wezterm

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/Dekr0/shutil/fzf"
	"github.com/valyala/fastjson"
)

func SwitchCurrentWindowTab() error {
	weztermQueryProc := exec.Command("wezterm", "cli", "list", "--format", "json")

    blob, queryErr := weztermQueryProc.Output()
    if queryErr != nil {
        return queryErr
    }
    json, parserErr := fastjson.Parse(string(blob))
    if parserErr != nil {
        return parserErr
    }
    jTabs, jsonErr := json.Array()
    if jsonErr != nil {
        return jsonErr
    }

    selections := make([][]byte, len(jTabs))
    for i, jTab := range jTabs {
        tabId := jTab.GetUint("tab_id")
        tabTitle := jTab.GetStringBytes("tab_title")
        if bytes.Compare(tabTitle, []byte("")) == 0 {
            tabTitle = []byte(fmt.Sprintf("Unnamed Tab %d", i))
        }
        selections[i] = []byte(fmt.Sprintf("%d - %s\n", tabId, tabTitle))
    }

    fzfProc := exec.Command("fzf")
    fzfPipe, pipeErr := fzfProc.StdinPipe()
    if pipeErr != nil {
        return pipeErr
    }

    cFzF := make(chan *fzf.FzFSimpleResponse, 1)

    go func() {
        o, e := fzfProc.Output()
        cFzF <- &fzf.FzFSimpleResponse{Selection: o, Error: e}
    }()

    for _, selection := range selections {
        _, writeErr := fzfPipe.Write(selection)
        if writeErr != nil {
            return writeErr
        }
    }

    for {
        select {
        case r := <- cFzF:
            if r.Error != nil {
                return r.Error
            }
            splits := bytes.Split(r.Selection, []byte(" - "))
            if len(splits) != 2 {
                return errors.New("Malformed FzF response")
            }
            weztermProc := exec.Command("wezterm", "cli", "activate-tab", "--tab-id", string(splits[0]))
            return weztermProc.Run()
        default:
            time.Sleep(time.Microsecond * 100)
        }
    }
}
