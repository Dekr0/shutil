package wezterm

import (
    "testing"
)

func TestWeztermTabSwitch(t *testing.T) {
    err := SwitchCurrentWindowTab()
    if err != nil {
        t.Fatal(err)
    }
}
