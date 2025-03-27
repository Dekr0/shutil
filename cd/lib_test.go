package cd

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"
)


func TestSearchDir(t *testing.T)  {
    /*f, err := os.Create("log.txt")
    if err != nil {
        t.Fatalf(err.Error())
    }*/
    logger := slog.New(slog.NewJSONHandler(
        os.Stdout,
        &slog.HandlerOptions {
            Level: slog.LevelInfo,
        },
    ))

    out, err := SearchDir([]string{"D:/codebase/hd2_asset_db"}, 3, 8, logger)
    if err != nil {
        t.Fatalf(err.Error())
    }

    t.Log(string(out))
}

func _TestCancel(t *testing.T) {
    bg := context.Background()

    ctx, cancel := context.WithCancel(bg)
    defer cancel()

    var wg sync.WaitGroup

    wg.Add(1)
    go func() {
        defer wg.Done()
        for i := 0; i < 100; i++ {
            select {
            case <- ctx.Done():
                t.Log("Aborted operation.")
                return
            default:
                t.Log(i)
                time.Sleep(time.Second * 2)
            }
        }
    }()

    time.Sleep(time.Second * 10)
    cancel()
    t.Log("Sent cancel signal.")
    wg.Wait()
}
