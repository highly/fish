package watcher

import (
	"github.com/fsnotify/fsnotify"
	"github.com/highly/fish/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var allowExtension = []string{".go", ".yaml"}

type Watch struct {
	watch *fsnotify.Watcher
}

func NewWatcher() *Watch {
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	return &Watch{watch: watch}
}

func (w *Watch) PreBuildRun() *Watch {
	utils.GoBuildAndRun(utils.CurrentPath)
	return w
}

func (w *Watch) Run() {
	defer func() {
		if r := recover(); r != nil {
			utils.Red.Println("found recover: ", r)
		}
		w.watch.Close()
	}()

	if err := w.watchDir(utils.CurrentPath); err != nil {
		log.Fatalln(err.Error())
	}

	go w.Action()
	go w.Clean()

	select {}
}

func (w *Watch) Action() {
	for {
		select {
		case ev := <-w.watch.Events:
			{
				if ev.Op&fsnotify.Create == fsnotify.Create {
					if utils.IsDirectory(ev.Name) {
						w.watch.Add(ev.Name)
					} else {
						if w.IsFileAllowed(ev.Name) {
							utils.TwoColorPrintLn(ev.Name, " han been modified")
							utils.GoBuildAndRun(utils.CurrentPath)
						}
					}
				}

				if (ev.Op&fsnotify.Remove == fsnotify.Remove) ||
					(ev.Op&fsnotify.Rename == fsnotify.Rename) {
					if utils.IsDirectory(ev.Name) {
						w.watch.Remove(ev.Name)
					}
				}
			}
		case err := <-w.watch.Errors:
			{
				utils.Red.Println("watcher chan error : ", err)
				return
			}
		}
	}
}

func (w *Watch) watchDir(dir string) error {
	utils.TwoColorPrintLn("start watching location => ", dir)
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			if err = w.watch.Add(path); err != nil {
				return err
			}
		}
		return nil
	})
}

func (w *Watch) IsFileAllowed(fileName string) bool {
	for _, ext := range allowExtension {
		if strings.HasSuffix(fileName, ext) {
			return true
		}
	}
	return false
}

func (w *Watch) Clean() {
	utils.Clean()
}
