package controllers

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"time"

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
	search        string
	datestartend  string
	savefiles     []string
	stringF       bool
	SearchMap     map[string]logenc.LogList
	date_layout   = "01/02/2006"
	startUnixTime int64
	endUnixTime   int64
)

type MyStruct struct {
	DirN string
	File string
}

// RootHandler - http handler for handling / path
func RootHandler(w http.ResponseWriter, _ *http.Request) {

	files := []string{
		"web/templates/index.tmpl",
		"web/templates/footer.tmpl",
		//"./ui/html/footer.partial.tmpl",
		"web/templates/header.tmpl",
		"web/templates/wscontent.tmpl",
		"web/templates/card.tmpl",
	}
	t := template.New("index").Delims("<<", ">>")

	t, err := t.Parse("footer")
	if err != nil {
		log.Fatal("Problem with template \"footer\"")
	}
	t, err = t.Parse("header")
	if err != nil {
		log.Fatal("Problem with template \"header\"")
	}
	t, err = t.Parse("loading")
	if err != nil {
		log.Fatal("Problem with template \"header\"")
	}
	t, err = t.Parse("wscontent")
	if err != nil {
		log.Fatal("Problem with template \"wscontent\"")
	}
	t, err = t.Parse("card")
	if err != nil {
		log.Fatal("Problem with template \"card\"")
	}
	t, err = t.ParseFiles(files...)
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
	if stringF {
		Indexing(conn, filename)
		savefiles = append(savefiles, filename)

	}

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
		search = ""
	}
	//w.WriteHeader(http.StatusUnauthorized)
}

func SearchHandler(_ http.ResponseWriter, r *http.Request) {
	search = r.URL.Query().Get("search_string")
	fmt.Println("SEARCHHANDLER:", search)
}
func DataHandler(_ http.ResponseWriter, r *http.Request) {
	datestartend = r.URL.Query().Get("daterange")
	fmt.Println("DATAHANDLER:", datestartend)
	//SEARCHHANDLER: 01/01/2021 - 01/15/2021
	datastart := string(datestartend[0:10])
	/* daystart := string(datestartend[0:2])
	monthstart := string(datestartend[3:5])
	yearstart := string(datestartend[6:10])

	dayend := string(datestartend[13:15])
	monthend := string(datestartend[16:18])
	yearend := string(datestartend[19:23]) */
	dataend := string(datestartend[13:23])

	timeendUnix, _ := time.Parse(date_layout, "01/15/2021")
	timestartUnix, _ := time.Parse(date_layout, datastart)

	fmt.Println("Common", dataend, "Unix", timeendUnix.Unix())
	fmt.Println("Common", datastart, "Unix", timestartUnix.Unix())

	//fmt.Println("Parse d:m:y", daystart, ":", monthstart, ":", yearstart)
	//fmt.Println("Parse d:m:y", dayend, ":", monthend, ":", yearend)
	startUnixTime = timestartUnix.Unix()
	endUnixTime = timeendUnix.Unix()

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
			bleveSI.ProcBleve(fileN, fileaddr)
			conn.WriteMessage(websocket.TextMessage, []byte("Indexing complated"))
		}()
		SearchMap = logenc.ProcMapFile(fileaddr)
	}
}

//View List of Dir
//NOT fileUtils !!!
func ViewDir(conn *websocket.Conn, search string) {

	//:TODO
	//Рассортировать файлы по Ulid ??
	//в порядке возрастания по времения от самого старого до нового
	//также они отображаются и в Tollbar
	//

	//Добавить отслеживаеи последнего выведенного ulid
	//Разбить передеваеммую структуру на 3*600 записей
	//Всегда отоюражать 2 из 3 структур (кроме начала)
	//Переходим на 3 структуру удаляем первые 600 записей и подгружаем 600 следующих
	//Dinamic
	//
	var lastUlid string
	var commoncsv logenc.LogList
	var fileList = make(map[string][]string)
	files, _ := ioutil.ReadDir("./repdata")
	//"/home/nik/projects/Course/tmcs-log-agent-storage/"
	//"./view"
	countFiles := (len(files))
	conn.WriteMessage(websocket.TextMessage, []byte("Indexing file, please wait"))

	fileList["FileList"] = util.Conf.Dir
	//String[] values = fileList.get("FileList");
	fmt.Println("start")

	for i := 0; i < countFiles; i++ {
		fileaddr := fileList["FileList"][i]
		fileN := filepath.Base(fileaddr)
		go logenc.Replication(fileaddr)
		bleveSI.ProcBleve(fileN, fileaddr)
		commoncsv = util.TailDir(conn, fileaddr, search, SearchMap, startUnixTime, endUnixTime, commoncsv, lastUlid)
		//conn.WriteMessage(websocket.TextMessage, []byte(filepath.Base(fileList["FileList"][i])))
	}
	conn.WriteMessage(websocket.TextMessage, []byte(logenc.EncodeXML(commoncsv)))
	conn.WriteMessage(websocket.TextMessage, []byte("Indexing complated"))
	search = ""
	startUnixTime = 0
	endUnixTime = 0
}
