package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

const BinFIle = "localTmpMain"

var hiddenLog = []string{
	"grpc/pickfirst.go",
	"resolver_conn_wrapper",
	"grpc/clientconn.go",
	"gateway_server.go",
	"tracing/config.go",
	"HandleSubConnStateChange",
	"parsed scheme",
	"pick_first",
	"not registered,",
}

var ctx context.Context
var cancel context.CancelFunc
var outputChan chan logMsg

type logMsg struct {
	Type, Message, Message2 string
}

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
		close(outputChan)
		cancel()
	}
	TwoColorPrintLn("running new program => ", GetProgramName(location))

	ctx, cancel = context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, filepath.Join(location, BinFIle))
	cmd.SysProcAttr = &syscall.SysProcAttr{Foreground: false}
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	cmd.Start()

	outputChan = make(chan logMsg, 1024)
	go OutputStdOut(stdout)
	go OutputStdErr(stderr)
	go Output()
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

type outputLog struct {
	Level   string      `json:"level"`
	Message string      `json:"message"`
	Context interface{} `json:"context"`
}

func OutputStdErr(stderr io.ReadCloser) {
	for {
		buf := make([]byte, 10240)
		logLen, err := stderr.Read(buf)
		if err != nil {
			return
		}
		if logLen <= 0 {
			break
		}
		outContent := string(buf[:logLen])
		show := true
		for _, banned := range hiddenLog {
			if strings.Contains(outContent, banned) {
				show = false
				break
			}
		}
		if !show {
			continue
		}
		var logContent outputLog
		err = json.Unmarshal(buf[:logLen], &logContent)
		if err != nil {
			outputChan <- logMsg{
				Type:    "stdErr_ori",
				Message: outContent,
			}
		} else {
			context := fmt.Sprintf("%v", logContent.Context)
			if context != "map[]" {
				logContent.Message = logContent.Message + "\n" + context
			}
			outputChan <- logMsg{
				Type:     "stdErr_parsed",
				Message:  logContent.Level + ": ",
				Message2: logContent.Message,
			}
		}
	}
}

func OutputStdOut(stdout io.ReadCloser) {
	for {
		buf := make([]byte, 10240)
		logLen, err := stdout.Read(buf)
		if err != nil {
			return
		}
		if logLen <= 0 {
			break
		}
		outputChan <- logMsg{
			Type:    "stdOut",
			Message: string(buf[:logLen]),
		}
	}
}

func Output() {
	for msg := range outputChan {
		if msg.Type == "stdErr_ori" {
			FgMagenta.Printf("%s\n", msg.Message)
		} else if msg.Type == "stdErr_parsed" {
			WhiteAndFgMagenta(msg.Message, msg.Message2)
		} else {
			Yellow.Printf("\n%s\n", msg.Message)
		}
	}
}
