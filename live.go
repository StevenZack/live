package main

import (
	"github.com/StevenZack/live/views"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"

	"github.com/StevenZack/tools/netToolkit"

	"github.com/fsnotify/fsnotify"
)

var (
	upgrader   = websocket.Upgrader{}
	notifyChan = make(chan bool)
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("args not enough")
		return
	}
	watcher, e := fsnotify.NewWatcher()
	if e != nil {
		fmt.Println(`new watcher error :`, e)
		return
	}
	defer watcher.Close()
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("modified file:", event.Name)
					notifyChan <- true
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println(err)
			}
		}
	}()

	e = watcher.Add(os.Args[1])
	if e != nil {
		fmt.Println(`add error :`, e)
		return
	}
	http.HandleFunc("/", home)
	http.HandleFunc("/live/ws", preview)
	http.HandleFunc("/live/live.js", handleJs)
	fmt.Println(strings.Join(netToolkit.GetIPs(), "\n"))
	e = http.ListenAndServe(":8080", nil)
	if e != nil {
		fmt.Println(`listen error :`, e)
		return
	}
}
func home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, os.Args[1])
}
func preview(w http.ResponseWriter, r *http.Request) {
	c, e := upgrader.Upgrade(w, r, nil)
	if e != nil {
		fmt.Println(`upgrade error :`, e)
		return
	}
	fmt.Println("websocket connected")
	defer c.Close()
	defer fmt.Println("websocket disconnected")
	for v := range notifyChan {
		if !v {
			break
		}
		c.WriteMessage(websocket.TextMessage, []byte("changed"))
	}
}
func handleJs(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/javascript")
	w.Header().Add("Content-Length", fmt.Sprintf("%v", views.Str_live))
	fmt.Fprint(w, views.Str_live)
}