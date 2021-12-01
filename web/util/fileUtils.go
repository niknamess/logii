package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hpcloud/tail"
	"gitlab.topaz-atcs.com/tmcs/logi2/bleveSI"
	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
	//"gitlab.topaz-atcs.com/tmcs/logi2/web/controllers"
)

var (
	// FileList - list of files that were parsed from the provided config
	FileList []string
	visited  map[string]bool

	// Global Map that stores all the files, used to skip duplicates while
	// subsequent indexing attempts in cron trigger
	indexMap  = make(map[string]bool)
	SearchMap map[string]string
)

type FileStruct struct {
	ID      int    `json:"id"`
	NAME    string `json:"filename"`
	HASHSUM string `json:"hashsum"`
}

// TailFile - Accepts a websocket connection and a filename and tails the
// file and writes the changes into the connection. Recommended to run on
// a thread as this is blocking in nature
func TailFile(conn *websocket.Conn, fileName string, lookFor string, SearchMap map[string]string) {

	fileN := filepath.Base(fileName)
	UlidC := bleveSI.ProcBleveSearchv2(fileN, lookFor)

	taillog, err := tail.TailFile(fileName,
		tail.Config{
			Follow: true,
			Location: &tail.SeekInfo{
				Whence: os.SEEK_CUR, //!!!

			},
		})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error occurred in opening the file: ", err)
		return
	}

	if lookFor == "" || lookFor == " " || lookFor == "Search" {
		for line := range taillog.Lines {
			//if dir == false {
			conn.WriteMessage(websocket.TextMessage, []byte(logenc.ProcLine(line.Text)))
			//} else {
			//	return true
			//}

		}
	} else if len(UlidC) == 0 {
		println("Break")
		return
	} else {

		for i := 0; i < len(UlidC); i++ {

			v, found := SearchMap[UlidC[i]]
			if found == true {

				conn.WriteMessage(websocket.TextMessage, []byte(v))

			}
		}

	}
	//return true
}

// IndexFiles - takes argument as a list of files and directories and returns
// a list of absolute file strings to be tailed
func IndexFiles(fileList []string) error {
	// Re-initialize the visited array
	visited = make(map[string]bool)

	// marking all entries possible stale
	// will be removed if not updated
	for k := range indexMap {
		indexMap[k] = false
	}

	for _, file := range fileList {
		dfs(file)
	}
	// Re-initialize the file list array
	FileList = make([]string, 0)

	// Iterate through the map that contains the filenames
	for k, v := range indexMap {
		if v == false {
			delete(indexMap, k)
			continue
		}
		//fmt.Fprintln(os.Stderr, k)
		FileList = append(FileList, k)
	}
	Conf.Dir = FileList
	fmt.Fprintln(os.Stderr, "Indexing complete !, file index length: ", len(Conf.Dir))
	return nil
}

/* skip all files that are :
   a: append-only
   l: exclusive use
   T: temporary file; Plan 9 only
   L: symbolic link
   D: device file
   p: named pipe (FIFO)
   S: Unix domain socket
   u: setuid
   g: setgid
   c: Unix character device, when ModeDevice is set
   t: sticky
*/
func dfs(file string) {
	// Mostly useful for first entry, as the paths may be like ../dir or ~/path/../dir
	// or some wierd *nixy style, Once the file is cleaned and made into an absolute
	// path, it should be safe as the next call is basepath(abspath) + "/" + name of
	// the file which should be accurate in all terms and absolute without any
	// funky conversions used in OS
	file = filepath.Clean(file)
	absPath, err := filepath.Abs(file)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get absolute path of the file %s; err: %s\n", file, err)
	}
	if _, ok := visited[file]; ok {
		// if the absolute path has been visited, return without processing
		return
	}
	visited[file] = true
	s, err := os.Stat(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to stat file %s; err: %s\n", file, err)
		return
	}
	// check if the file is a directory
	if s.IsDir() {
		basepath := filepath.Clean(file)
		filelist, _ := ioutil.ReadDir(absPath)
		for _, f := range filelist {
			dfs(basepath + string(os.PathSeparator) + f.Name())
		}
	} else if strings.ContainsAny(s.Mode().String(), "alTLDpSugct") {
		// skip these files
		// @TODO try including names PIPES
	} else {
		// only remaining file are ascii files that can be then differentiated
		// by the user as golang has only these many categorization
		// Note : this appends the absolute paths
		// Insert the absPath into the Map, avoids duplicates in successive cron runs
		indexMap[absPath] = true
	}
}

func TailDir(conn *websocket.Conn, fileName string, lookFor string, SearchMap map[string]string) bool {

	fileN := filepath.Base(fileName)
	UlidC := bleveSI.ProcBleveSearchv2(fileN, lookFor)

	if len(UlidC) == 0 {
		println("Break")
		return false
	} else {

		for i := 0; i < len(UlidC); i++ {

			_, found := SearchMap[UlidC[i]]
			if found == true {
				return true

			}
		}

	}
	return false

}

func Indexing(conn *websocket.Conn, filename string) {

	if filename == "undefined" {
		return
	} else {
		fileN := filepath.Base(filename)
		fmt.Println(filename)
		go logenc.Replication(filename)
		go func() {
			conn.WriteMessage(websocket.TextMessage, []byte("Indexing file, please wait"))
			bleveSI.ProcBlev(fileN, filename)
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

		fileList["FileList"] = Conf.Dir
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
		fileList["FileList"] = Conf.Dir
		//String[] values = fileList.get("FileList");
		fmt.Println("start")
		for i := 0; i < countFiles; i++ {
			fileaddr := fileList["FileList"][i]
			fileN := filepath.Base(fileaddr)
			go logenc.Replication(fileaddr)
			bleveSI.ProcBlev(fileN, fileaddr)
			if TailDir(conn, fileN, search, SearchMap) == true {
				conn.WriteMessage(websocket.TextMessage, []byte(fileList["FileList"][i]))
			}
			//fmt.Println(fileaddr)
		}
	}
	conn.WriteMessage(websocket.TextMessage, []byte("Indexing complated"))

}

func GetFiles(address string, port string) {
	resp, err := http.Get("http://" + address + ":" + port + "/vfs/data/")
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range logenc.GetLinks(resp.Body) {

		fullURLFile := "http://" + address + ":" + port + "/vfs/data/" + v

		fileURL, err := url.Parse(fullURLFile)
		if err != nil {
			log.Fatal(err)
		}
		path := fileURL.Path
		segments := strings.Split(path, "/")
		fileName := segments[len(segments)-1]

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

		_, err = io.Copy(file, resp.Body)
		logenc.Replication("./testsave/" + fileName)
		fmt.Println("Merge", fileName)
		defer file.Close()
		logenc.DeleteOldsFiles("./testsave/", fileName, "")
		//fmt.Printf("Downloaded a file %s with size %d", fileName, size)
	}

}

func DeleteFile90(dir string) {
	for {
		var cutoff = 24 * time.Hour * 90
		fileInfo, err := ioutil.ReadDir("./")
		if err != nil {
			log.Fatal(err.Error())
		}
		now := time.Now()
		fmt.Println(now)
		for _, info := range fileInfo {
			if diff := now.Sub(info.ModTime()); diff > cutoff {
				fmt.Printf("Deleting %s which is %s old\n", info.Name(), diff)
				logenc.DeleteOldsFiles(dir, info.Name(), "")

			}
		}
	}
}
