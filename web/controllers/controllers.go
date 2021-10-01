package controllers

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"gitlab.topaz-atcs.com/tmcs/logi2/web/util"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// RootHandler - http handler for handling / path
func RootHandler(w http.ResponseWriter, r *http.Request) {
	t := template.New("index").Delims("<<", ">>")
	t, err := t.ParseFiles("web/templates/index.tmpl")
	t = template.Must(t, err)
	if err != nil {
		panic(err)
	}
	fmt.Println("root_h")
	var fileList = make(map[string]interface{})

	fileList["FileList"] = util.Conf.Dir
	fileList[csrf.TemplateTag] = csrf.Token(r)
	fileList["token"] = csrf.Token(r)
	t.Execute(w, fileList)
}

// WSHandler - Websocket handler
func WSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	filenameB, _ := base64.StdEncoding.DecodeString(mux.Vars(r)["b64file"])
	filename := string(filenameB)
	// sanitize the file if it is present in the index or not.
	filename = filepath.Clean(filename)
	ok := false
	for _, wFile := range util.Conf.Dir {
		if filename == wFile {
			ok = true
			break
		}
	}

	// If the file is found, only then start tailing the file.
	// This is to prevent arbitrary file access. Otherwise send a 403 status
	// This should take care of stacking of filenames as it would first
	// be searched as a string in the index, if not found then rejected.
	if ok {
		//go util.TailFile(conn, filename, search)
		go Runw(conn, filename)
	}
	w.WriteHeader(http.StatusUnauthorized)
}

func Runw(conn *websocket.Conn, filename string) {
	search := "32"
	//http.HandleFunc("/searchproject",  func(w http.ResponseWriter, r *http.Request))
	//search = r.URL.Query().Get("search_string")
	//fs := http.FileServer(http.Dir(self.staticFs))
	//mux := http.NewServeMux()
	//mux.HandleFunc("/searchproject/", handleGet)
	//http.HandleFunc("/searchproject", handleGet)

	//http.HandleFunc("/searchproject", func(w http.ResponseWriter, r *http.Request) {
	//	search = r.URL.Query().Get("search_string")
	//	fmt.Fprintf(w, "ПОИСК: %s", search)
	//})

	util.TailFile(conn, filename, search)

}
