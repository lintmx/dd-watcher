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
	TimeTicker time.Duration
	LiveStatus bool
}

// Run a dd monitor
func (w *Watcher) Run(ctx context.Context) {
	timer := time.NewTimer(w.TimeTicker * time.Second)
	defer w.Wait.Done()
	fmt.Fprintf(os.Stdout, "开始监听： %s\n", w.LiveAPI.GetLiveURL())

	w.refresh()
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			w.refresh()
			timer.Reset(w.TimeTicker * time.Second)
		}
	}
}

// Refresh live status
func (w *Watcher) refresh() {
	err := w.LiveAPI.RefreshLiveInfo()
	if err != nil {
		// fmt.Fprintf(os.Stderr, "refersh err - %s : %s\n", w.LiveAPI.GetLiveURL(), err)
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
		} else {
			fmt.Fprintf(os.Stdout, "%s: %s - 下了，神回，没录，不烤\n",
				time.Now().Format("2006-01-02 15:04:05"),
				w.LiveAPI.GetAuthor(),
			)
		}
	}
}
