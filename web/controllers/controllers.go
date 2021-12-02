package controllers

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"

	//"encoding/json"

	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"gitlab.topaz-atcs.com/tmcs/logi2/bleveSI"
	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
	"gitlab.topaz-atcs.com/tmcs/logi2/web/util"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	search    string
	savefiles []string
	stringF   bool
	SearchMap map[string]string
)

type MyStruct struct {
	DirN string
	File string
}

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

	t.Execute(w, fileList)
}

// WSHandler - Websocket handler
func WSHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	filenameB, _ := base64.StdEncoding.DecodeString(mux.Vars(r)["b64file"])

	filename := string(filenameB)
	if filenameB == nil {
		return
	}

	if filename == "undefined" {
		ViewDir(conn, search)
	}

	if savefiles == nil {
		Indexing(conn, filename)
		savefiles = append(savefiles, filename)
	} else {
		for i := 0; i < len(savefiles); i++ {
			if filename != savefiles[i] {
				stringF = true
			} else {
				stringF = false
			}
		}

	}
	if stringF == true {
		Indexing(conn, filename)
		savefiles = append(savefiles, filename)

	}
	///logenc.ProcFileBreveSLOWLY(fileN, filename)
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

		//util.TailFile(conn, filename, search, SearchMap, false)
		util.TailFile(conn, filename, search, SearchMap)
	}
	w.WriteHeader(http.StatusUnauthorized)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	search = r.URL.Query().Get("search_string")

}

//NOT fileUtils !!!
func Indexing(conn *websocket.Conn, fileaddr string) {
	//var SearchMap map[string]string
	if fileaddr == "undefined" {
		return
	} else {
		fileN := filepath.Base(fileaddr)
		fmt.Println(fileaddr)
		go logenc.Replication(fileaddr)
		go func() {
			conn.WriteMessage(websocket.TextMessage, []byte("Indexing file, please wait"))
			bleveSI.ProcBlev(fileN, fileaddr)
			conn.WriteMessage(websocket.TextMessage, []byte("Indexing complated"))
		}()
		SearchMap = logenc.ProcMapFile(fileaddr)
	}
}

//View List of Dir
//NOT fileUtils !!!
func ViewDir(conn *websocket.Conn, search string) {
	var fileList = make(map[string][]string)
	files, _ := ioutil.ReadDir("./repdata")
	//"/home/nik/projects/Course/tmcs-log-agent-storage/"
	//"./view"
	countFiles := (len(files))
	conn.WriteMessage(websocket.TextMessage, []byte("Indexing file, please wait"))
	if len(search) == 0 {

		fileList["FileList"] = util.Conf.Dir
		//String[] values = fileList.get("FileList");
		fmt.Println("start")
		for i := 0; i < countFiles; i++ {
			fileaddr := fileList["FileList"][i]
			fileN := filepath.Base(fileaddr)
			go logenc.Replication(fileaddr)
			bleveSI.ProcBlev(fileN, fileaddr)
			conn.WriteMessage(websocket.TextMessage, []byte(fileList["FileList"][i]))

		}

	} else {
		fileList["FileList"] = util.Conf.Dir
		//String[] values = fileList.get("FileList");
		fmt.Println("start")
		for i := 0; i < countFiles; i++ {
			fileaddr := fileList["FileList"][i]
			fileN := filepath.Base(fileaddr)
			go logenc.Replication(fileaddr)
			bleveSI.ProcBlev(fileN, fileaddr)
			if util.TailDir(conn, fileN, search, SearchMap) {
				conn.WriteMessage(websocket.TextMessage, []byte(fileList["FileList"][i]))
			}
			//fmt.Println(fileaddr)
		}
	}
	conn.WriteMessage(websocket.TextMessage, []byte("Indexing complated"))

}
