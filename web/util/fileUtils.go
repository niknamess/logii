package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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
	indexMap = make(map[string]bool)
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
	UlidC := bleveSI.ProcBleveSearch(fileN, lookFor)

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
	UlidC := bleveSI.ProcBleveSearch(fileN, lookFor)

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

/*
func AddJsonInfo(conn *websocket.Conn) []byte {
	dirpath := "/home/nik/projects/Course/tmcs-log-agent-storage/"
	var idents []FileStruct
	var fileList = make(map[string][]string)
	files, _ := ioutil.ReadDir(dirpath)
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
		//res2B, _ := json.Marshal(idents)
		idents = append(idents, group)
		fmt.Println(idents)
	}
	res2B, _ := json.Marshal(idents)
	conn.WriteMessage(websocket.TextMessage, res2B)
	//result, _ := json.Marshal(idents)
	return res2B
}
*/
