package controllers

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	ctx "github.com/gorilla/context"
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
	search := "32 "
	//http.HandleFunc("/searchproject", func(w http.ResponseWriter, r *http.Request) {
	//	search = r.URL.Query().Get("search_string")
	//fmt.Fprintf(w, "ПОИСК: %s", search)
	//})
	util.TailFile(conn, filename, search)

}

// GetContext wraps each request in a function which fills in the context for a given request.
// This includes setting the User and Session keys and values as necessary for use in later functions.
func GetContext(handler http.Handler) http.HandlerFunc {
	// Set the context here
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request form
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing request", http.StatusInternalServerError)
		}
		w.Header().Set("X-CSRF-Token", csrf.Token(r))
		handler.ServeHTTP(w, r)
		// Remove context contents
		ctx.Clear(r)
	}
}

// CSRFExemptPrefixes - list of endpoints that does not require csrf protction
var CSRFExemptPrefixes = []string{
	// "/user",
}

// CSRFExceptions - exempts ajax calls from csrf tokens
func CSRFExceptions(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, prefix := range CSRFExemptPrefixes {
			if strings.HasPrefix(r.URL.Path, prefix) {
				r = csrf.UnsafeSkipCheck(r)
				break
			}
		}
		handler.ServeHTTP(w, r)
	}
}
