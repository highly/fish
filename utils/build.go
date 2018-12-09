package utils

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
)

const BinFIle = "localTmpMain"

var ctx context.Context
var cancel context.CancelFunc

func GoBuildAndRun(location string) {
	if GoBuild(location) {
		GoRun(location)
	}
}

func GoBuild(location string) bool {
	TwoColorPrintLn("start building program => ", GetProgramName(location))

	cmd := exec.Command("go", "build", "-v", "-o", BinFIle, location)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		Red.Println(stderr.String())
		return false
	}
	Common.Println("building finished.")
	return true
}

func GoRun(location string) {
	// stop program if already exists
	if cancel != nil {
		Common.Println("closing previous program...")
		cancel()
	}
	TwoColorPrintLn("running new program => ", GetProgramName(location))

	ctx, cancel = context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, filepath.Join(location, BinFIle))
	cmd.SysProcAttr = &syscall.SysProcAttr{Foreground: false}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
}

func Clean() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-signalCh:
		if cancel != nil {
			TwoColorPrintLn("\nclosing program => ", GetProgramName(CurrentPath))
			cancel()
		}
		DelFile(filepath.Join(CurrentPath, BinFIle))
		Red.Println("exit from fish.")
		os.Exit(1)
	}
}
