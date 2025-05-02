package fzf

import (
	"context"
	"io"
	"os/exec"
)

type FzFSimpleResponse struct {
	Selected []byte
	Error    error
}

func SimplePickAsync(ctx context.Context, selections [][]byte) (
	chan *FzFSimpleResponse,
	error,
) {
    proc := exec.CommandContext(ctx, "fzf")

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

func SimplePickSync(selections []byte) ([]byte, error) {
	proc := exec.Command("fzf")
	pipe, err := proc.StdinPipe()
	if err != nil {
		return nil, err
	}
	_, err = pipe.Write(selections)
	if err != nil {
		return nil, err
	}
	res, err := proc.Output()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func SimplePickAsyncDual(ctx context.Context) (
	io.WriteCloser,
	<- chan *FzFSimpleResponse,
	error,
) {
	proc := exec.CommandContext(ctx, "fzf")
	pipe, err := proc.StdinPipe()
	if err != nil {
		return nil, nil, err
	}

    c := make(chan *FzFSimpleResponse, 1)

	go func() {
		o, e := proc.Output()
		c <- &FzFSimpleResponse{o, e}
	}()

	return pipe, c, nil
}
