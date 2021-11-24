package controllers

import (
	"encoding/base64"
	"io"
	"log"
	"os"
	"strings"

	//"encoding/json"

	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
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

//type FileList struct {
//JNAME json.Name    `json:"filelist"`
//files []FileStruct `json:"file"`
//}

type FileStruct struct {
	ID      int    `json:"id"`
	NAME    string `json:"filename"`
	HASHSUM string `json:"hashsum"`
}
type UlpoadFileStruct struct {
	ID      int    `json:"id"`
	NAME    string `json:"filename"`
	Content []byte
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

//Bleve indexing and create map
func Indexing(conn *websocket.Conn, filename string) {

	if filename == "undefined" {
		return
	} else {
		fileN := filepath.Base(filename)
		fmt.Println(filename)
		go logenc.Replication(filename)
		go func() {
			conn.WriteMessage(websocket.TextMessage, []byte("Indexing file, please wait"))
			bleveSI.ProcFileBreveSLOWLY(fileN, filename)
			conn.WriteMessage(websocket.TextMessage, []byte("Indexing complated"))
		}()
		SearchMap = logenc.ProcMapFile(filename)
	}
}

//View List of Dir
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
			bleveSI.ProcFileBreveSLOWLY(fileN, fileaddr)
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
			bleveSI.ProcFileBreveSLOWLY(fileN, fileaddr)
			if util.TailDir(conn, fileN, search, SearchMap) == true {
				conn.WriteMessage(websocket.TextMessage, []byte(fileList["FileList"][i]))
			}
			//fmt.Println(fileaddr)
		}
	}
	conn.WriteMessage(websocket.TextMessage, []byte("Indexing complated"))

}

/*
func JsonEncode() []FileStruct {

	var idents []FileStruct
	var fileList = make(map[string][]string)
	files, _ := ioutil.ReadDir("/home/nik/projects/Course/tmcs-log-agent-storage/")
	countFiles := (len(files))

	fileList["FileList"] = util.Conf.Dir
	for i := 0; i < countFiles; i++ {
		fileaddr := fileList["FileList"][i]
		fileN := filepath.Base(fileaddr)
		IDfile, _ := strconv.Atoi(logenc.Remove(fileN, '-'))
		hashsumfile := logenc.FileMD5(fileaddr)
		group := FileStruct{
			ID:      IDfile,
			NAME:    fileN,
			HASHSUM: hashsumfile,
		}
		//res2B, _ := json.Marshal(group)
		idents = append(idents, group)
		fmt.Println(idents)
	}
	//result, _ := json.Marshal(idents)
	return idents
}
*/
/*
func Sendfile(FileN string) UlpoadFileStruct {
	var group UlpoadFileStruct
	dirpath := "/home/nik/projects/Course/tmcs-log-agent-storage/"
	IDfile, _ := strconv.Atoi(logenc.Remove(FileN, '-'))
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {

		log.Fatal(err)
	}

	for _, f := range files {
		//fmt.Println(f.Name())
		if f.Name() == FileN {
			original, err := os.ReadFile(dirpath + FileN)
			if err != nil {
				log.Fatal(err)
			}

			group = UlpoadFileStruct{
				ID:      IDfile,
				NAME:    FileN,
				Content: original,
			}
			return group
		}
	}
	return group
}
*/
/*
func BodyHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	filenameB, _ := base64.StdEncoding.DecodeString(mux.Vars(r)["b64file"])

	FilePath := string(filenameB)
	if filenameB == nil {
		return
	}

	ok := false
	for _, wFile := range util.Conf.Dir {
		if FilePath == wFile {
			ok = true
			break
		}
	}
	if ok {

		//util.TailFile(conn, filename, search, SearchMap, false)
		jsonmes := util.AddJsonInfo(conn)
		fmt.Println(jsonmes)
	}

	w.WriteHeader(http.StatusUnauthorized)
}
*/
//upload file on other port
/*
func Upload(filename string, hostName string) string {
	file, err := os.Open("../../send-files/" + filename)
	if err != nil {
		return "file not found"
	}
	defer file.Close()

	res, err := http.Post("http://"+hostName+":8080/upload?filename="+filename, "binary/octet-stream", file)
	if err != nil {
		return "file not send"
	}
	defer res.Body.Close()
	message, _ := ioutil.ReadAll(res.Body)
	fmt.Printf(string(message))

	return string(message)
}

func Client() {

}
*/
func GetFiles(port string) {
	resp, err := http.Get("http://localhost:" + port + "/vfs/data/")
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range logenc.GetLinks(resp.Body) {

		fullURLFile = "http://localhost:" + port + "/vfs/data/" + v

		fileURL, err := url.Parse(fullURLFile)
		if err != nil {
			log.Fatal(err)
		}
		path := fileURL.Path
		segments := strings.Split(path, "/")
		fileName = segments[len(segments)-1]

		file, err := os.OpenFile("./testsave/"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		client := http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				r.URL.Opaque = r.URL.Path
				return nil
			},
		}
		// Put content on file
		resp, err := client.Get(fullURLFile)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		size, err := io.Copy(file, resp.Body)

		defer file.Close()

		fmt.Printf("Downloaded a file %s with size %d", fileName, size)
	}

}
