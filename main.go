package main

import (
	"github.com/highly/fish/watcher"
)

func main() {
	watcher.NewWatcher().PreBuildRun().Run()
}
