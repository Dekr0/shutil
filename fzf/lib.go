package fzf

import (
	"os/exec"
)

type FzFSimpleResponse struct {
	Selected []byte
	Error    error
}

func SimplePick(selections [][]byte) (chan *FzFSimpleResponse, error) {
    proc := exec.Command("fzf")

    pipe, err := proc.StdinPipe()
    if err != nil {
        return nil, err
    }

    c := make(chan *FzFSimpleResponse, 1)

    go func() {
        o, e := proc.Output()
        c <- &FzFSimpleResponse{Selected: o, Error: e}
    }()

    for _, selection := range selections {
        _, err := pipe.Write(selection)
        if err != nil {
            return nil, err
        }
    }

	return c, nil
}
