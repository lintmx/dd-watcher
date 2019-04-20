package watcher

import (
	"context"
	"fmt"
	"github.com/gen2brain/beeep"
	"github.com/lintmx/dd-watcher/api"
	"os"
	"sync"
	"time"
)

// Watcher struct
type Watcher struct {
	Wait       *sync.WaitGroup
	LiveAPI    api.LiveAPI
	TimeTicker *time.Ticker
	LiveStatus bool
}

// Run a dd monitor
func (w *Watcher) Run(ctx context.Context) {
	defer w.Wait.Done()
	defer w.TimeTicker.Stop()
	fmt.Fprintf(os.Stdout, "开始监听： %s\n", w.LiveAPI.GetLiveURL())

	w.refresh()
	for {
		select {
		case <-ctx.Done():
			return
		case <-w.TimeTicker.C:
			w.refresh()
		}
	}
}

// Refresh live status
func (w *Watcher) refresh() {
	err := w.LiveAPI.RefreshLiveInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "refersh err - %s : %s\n", w.LiveAPI.GetLiveURL(), err)
		return
	}

	if w.LiveStatus != w.LiveAPI.GetLiveStatus() {
		w.LiveStatus = w.LiveAPI.GetLiveStatus()

		if w.LiveStatus {
			title := fmt.Sprintf("%s 开播了！", w.LiveAPI.GetAuthor())
			body := fmt.Sprintf("%s 在 %s 开播了！", w.LiveAPI.GetAuthor(), w.LiveAPI.GetPlatformName())
			fmt.Fprintf(os.Stdout, "%s: %s - %s\n",
				time.Now().Format("2006-01-02 15:04:05"),
				body,
				w.LiveAPI.GetLiveURL(),
			)
			err := beeep.Alert(title, body, "")
			if err != nil {
				fmt.Fprintf(os.Stderr, "beeep err - %s\n", w.LiveAPI.GetLiveURL())
				return
			}
		}
	}
}
