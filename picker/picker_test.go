package picker

import (
	"context"
	"testing"
	"time"
)

func TestSearchExecutable(t *testing.T) {
	timeout, cancel := context.WithTimeout(context.Background(), time.Second * 32)
	defer cancel()
	app, err := SearchExecutable(timeout)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", app)
}
