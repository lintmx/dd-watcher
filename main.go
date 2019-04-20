package main

import (
	"context"
	"fmt"
	"github.com/lintmx/dd-watcher/api"
	"github.com/lintmx/dd-watcher/watcher"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

// Config struct
type Config struct {
	Interval uint32   `yaml:"interval"`
	Proxy    string   `yaml:"proxy"`
	Rooms    []string `yaml:"rooms"`
}

func main() {
	fmt.Fprintf(os.Stdout, "誰でも大好き！\n")
	var wait sync.WaitGroup
	defer wait.Wait()
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Get Pwd Error.\n")
		os.Exit(1)
	}

	conf := filepath.Join(currentDir, "config.yml")
	// Read Configuration file
	file, err := ioutil.ReadFile(conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read configuration file - %s\n", conf)
		os.Exit(1)
	}

	config := &Config{}
	// Parse yaml
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Configuration file parsing failed - %s\n", conf)
		os.Exit(1)
	}

	if config.Proxy != "" {
		proxyURL, err := url.Parse(config.Proxy)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Read Http Proxy Error.\n")
			os.Exit(1)
		}
		http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}

	ctx, cannel := context.WithCancel(context.Background())

	for _, room := range config.Rooms {
		u, err := url.Parse(room)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Room Url Parse Error - %s\n", room)
			continue
		}

		api := api.Check(u)
		if api == nil {
			fmt.Fprintf(os.Stderr, "Room not support - %s\n", u.Host)
		} else {
			wait.Add(1)
			w := &watcher.Watcher{
				LiveAPI:    api,
				TimeTicker: time.NewTicker(time.Duration(config.Interval) * time.Second),
				Wait:       &wait,
			}
			go w.Run(ctx)
		}
	}

	// Catch the exit signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		<-sigCh
		cannel()
	}()

	wait.Wait()
	fmt.Fprintf(os.Stdout, "\nさようなら～\n")
}
