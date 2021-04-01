package filewatcher

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"path/filepath"
	"strings"
	"sync"
)

/**
 * Copyright (C), 2019-2019
 * FileName: watcher
 * Author:   LinGuanHong
 * Date:     2019-11-12 10:33
 * Description: note
 */

// note:
// 		文件改变时候的监听，方便实现热加载
// 		watcherMap 对应多个 dir
// 		一个 watcher 对应 一个 dir，一个 dir 可监听多个 file
// 		按下 ctrl + s 保持的时候，才会触发。前后内容，如果没有改变，那么不会触发

type FileChangeWatcher struct {
	muxLock         sync.Mutex
	fileDir         string
	targetFileNames *[]string
	handleFuncMap   map[string]func(fileName string, optName fsnotify.Op)
	runOnce         *sync.Once
	watcherMap      map[string]*fsnotify.Watcher
}

func NewChangeWatcher(fileDir string) *FileChangeWatcher {
	if !strings.HasSuffix(fileDir, "/") {
		fileDir = fileDir + "/"
	}
	return &FileChangeWatcher{
		fileDir:         fileDir,
		muxLock:         sync.Mutex{},
		targetFileNames: &[]string{},
		handleFuncMap:   map[string]func(fileName string, optName fsnotify.Op){},
	}
}

func (w *FileChangeWatcher) PrintlnWatchCount() {
	fmt.Println("fileWatcher -->", "Dir:", len(w.watcherMap), "file:", len(*w.targetFileNames))
}

func (w *FileChangeWatcher) RegisterWatcher(fileName string, handlerFunc func(fileName string, optName fsnotify.Op)) error {
	initFunc := func() error {
		w.muxLock.Lock()
		defer w.muxLock.Unlock()
		if w.runOnce == nil {
			w.runOnce = &sync.Once{}
		}
		if w.watcherMap == nil {
			w.watcherMap = make(map[string]*fsnotify.Watcher, 0)
		}
		if w.watcherMap[w.fileDir] == nil {
			watcherChild, err := fsnotify.NewWatcher()
			if err != nil {
				return fmt.Errorf("FileWatcher new fs notify error: %s ", err.Error())
			}
			if watcherChild == nil {
				return fmt.Errorf("FileWatcher watcher == nil")
			}
			if err = watcherChild.Add(w.fileDir); err != nil {
				return fmt.Errorf("FileWatcher watcher Add dir error: %s, %s", err.Error(), w.fileDir)
			}
			w.watcherMap[w.fileDir] = watcherChild
		}
		*w.targetFileNames = append(*w.targetFileNames, fileName)
		w.handleFuncMap[fileName] = handlerFunc
		return nil
	}
	if err := initFunc(); err != nil {
		return err
	}
	w.runOnce.Do(func() {
		go func() {
			for {
				select {
				case event := <-w.watcherMap[w.fileDir].Events:
					eventFileName := filepath.Clean(event.Name)
					for _, name := range *w.targetFileNames {
						inputName := w.fileDir + name
						originName := filepath.Clean(inputName)
						if eventFileName != originName {
							continue
						}
						handler := w.handleFuncMap[name]
						if event.Op == fsnotify.Create || event.Op == fsnotify.Write || event.Op == fsnotify.Chmod {
							w.muxLock.Lock()
							// 复制粘贴会触发三个 Create,Write,Write
							switch event.Op {
							case fsnotify.Create:
								handler(name, fsnotify.Create)
								break
							case fsnotify.Write:
								handler(name, fsnotify.Write)
								break
							case fsnotify.Chmod:
								handler(name, fsnotify.Chmod)
								break
							}
							w.muxLock.Unlock()
						}
					}
					break
				case err := <-w.watcherMap[w.fileDir].Errors:
					if err != nil {
						fmt.Printf("FileWatcher fs notify err info ===> %s", err.Error())
					}
				}
			}
		}()
	})
	return nil
}
