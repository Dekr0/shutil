package kitty

import "testing"

func TestKitty(t *testing.T) {
	err := SwitchCurrentWindowTab()
	if err != nil {
		t.Fatal(err)
	}
}
