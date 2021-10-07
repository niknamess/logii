package controllers

import (
	"encoding/base64"
	"html/template"
	"net/http"
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
	search = ","
)

// RootHandler - http handler for handling / path
func RootHandler(w http.ResponseWriter, r *http.Request) {
	t := template.New("index").Delims("<<", ">>")
	t, err := t.ParseFiles("web/templates/index.tmpl")
	t = template.Must(t, err)
	if err != nil {
		panic(err)
	}
	//fmt.Println("root_h")
	var fileList = make(map[string]interface{})

	fileList["FileList"] = util.Conf.Dir
	fileList[csrf.TemplateTag] = csrf.Token(r)
	fileList["token"] = csrf.Token(r)
	t.Execute(w, fileList)
}

// WSHandler - Websocket handler
func WSHandler(w http.ResponseWriter, r *http.Request) {
	//search = r.URL.Query().Get("search_string")
	conn, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		//fmt.Fprintln(os.Stderr, err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	filenameB, _ := base64.StdEncoding.DecodeString(mux.Vars(r)["b64file"])
	//ProcFileBreve(filename)

	filename := string(filenameB)
	//logenc.ProcFileBreve(filename)
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
	//search = "0001GD3TH0Y1V8PBD2Z6DEY7PP"
	//search = "NTP"
	if ok {
		//go util.TailFile(conn, filename, search)
		util.TailFile(conn, filename, search)
	}
	w.WriteHeader(http.StatusUnauthorized)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	search = r.URL.Query().Get("search_string")

}
