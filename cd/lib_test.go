package cd

import "testing"

func TestSearchDir(t *testing.T)  {
    out, err := SearchDir([]string{"D:/codebase/hd2_asset_db"}, 2)
    if err != nil {
        t.Fatalf(err.Error())
    }
    t.Log(out)
}
