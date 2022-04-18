package util

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hpcloud/tail"
	"github.com/oklog/ulid/v2"
	"gitlab.topaz-atcs.com/tmcs/logi2/bleveSI"
	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
)

var paginationUlids map[int]string

var (
	FileName []string
	FileList []string
	visited  map[string]bool

	// Global Map that stores all the files, used to skip duplicates while
	// subsequent indexing attempts in cron trigger
	indexMap           = make(map[string]bool)
	signature   bool   = false
	Fname       string = ""
	currentfile string
	page        int = 0
	hashSumFile string
)

type FileStruct struct {
	ID      int    `json:"id"`
	NAME    string `json:"filename"`
	HASHSUM string `json:"hashsum"`
}

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

type Map map[string]string

// TailFile - Accepts a websocket connection and a filename and tails the
// file and writes the changes into the connection. Recommended to run on
// a thread as this is blocking in nature

func TailFile(conn *websocket.Conn, fileName string, lookFor string, SearchMap map[string]logenc.LogList) {
	fmt.Println("Fname", Fname)
	fileN := filepath.Base(fileName)

	if Fname != fileName {
		if Fname != "" {
			logenc.DeleteOldsFiles("./web/util/replace/"+filepath.Base(Fname), "")
		}
		Fname = fileName
		lookFor = ""

		for k := range paginationUlids {
			delete(paginationUlids, k)
		}
	}
	fmt.Println("Command", page)

	UlidC := bleveSI.ProcBleveSearchv2(fileN, lookFor)
	currentfile = fileN
	page = 0
	if lookFor == "" || lookFor == " " || lookFor == "Search" {

		hashSumFile = logenc.FileMD5(fileName)
		go func() {
			for {
				if hashSumFile != logenc.FileMD5(fileName) {
					hashSumFile = logenc.FileMD5(fileName)
					fmt.Println("hashSumFile", hashSumFile)
				} else if Fname != fileName {
					break
				}
			}
		}()
		go followCodeStatus(conn)
		UlidPaginationFile(conn, fileName)

		logenc.DeleteOldsFiles("./web/util/replace/"+fileN, "")
		conn.WriteMessage(websocket.TextMessage, []byte("<start></start>"))
		tailingLogsInFileAll(fileName, conn, 0, page)
		logenc.DeleteOldsFiles("./web/util/replace/"+fileN, "")

		var countline int = 0
		var currentpage int = 0
		for {
			if (logenc.FileMD5(fileName) != hashSumFile) && currentfile == fileN {
				UlidPaginationFile(conn, fileName)
				logenc.DeleteOldsFiles("./web/util/replace/"+fileN, "")
				conn.WriteMessage(websocket.TextMessage, []byte("<start></start>"))
				countline = tailingLogsInFileAll(fileName, conn, 0, page)
				logenc.DeleteOldsFiles("./web/util/replace/"+fileN, "")
				hashSumFile = logenc.FileMD5(fileName)
			} else if currentfile != fileN {

				break
			} else if countline >= 99 {
				//TransmitUlidPagination(conn, fileName)
				countline = 0
			} else if currentpage != page && currentfile == fileN {
				logenc.DeleteOldsFiles("./web/util/replace/"+fileN, "")
				conn.WriteMessage(websocket.TextMessage, []byte("<start></start>"))
				countline = tailingLogsInFileAll(fileName, conn, 0, page)
				logenc.DeleteOldsFiles("./web/util/replace/"+fileN, "")
				currentpage = page

			}

		}

		return

	} else if len(UlidC) == 0 {
		println("Break")
		return
	} else {

		var countCheck int

		fmt.Println("countSearch", 0)
		for i := 0; i < len(UlidC); i++ {
			_, found := SearchMap[UlidC[i]]
			if found {
				countCheck++
			}
		}
		fmt.Println("...............countCheck", countCheck)
		CountPage := "<countpage>" + strconv.Itoa(countCheck) + "</countpage>"
		conn.WriteMessage(websocket.TextMessage, []byte(CountPage))
		countCheck = 0
		currentpage := 0
		tailLogsInFind(SearchMap, UlidC, page, conn)

		for {

			if Fname != fileName {
				break
			} else if currentpage != page && currentfile == fileN {
				tailLogsInFind(SearchMap, UlidC, page, conn)
				currentpage = page

			}
		}

		return

	}

}
func tailLogsInFind(SearchMap map[string]logenc.LogList, UlidC []string, page int, conn *websocket.Conn) {
	//fmt.Println("//////////////////////////////page", page)
	if page == 0 {
		page = 1
	}

	for i := page*100 - 100; i < page*100; i++ {
		v, found := SearchMap[UlidC[i]]
		if found {
			fmt.Println(logenc.EncodeXML(v))
			conn.WriteMessage(websocket.TextMessage, []byte(logenc.EncodeXML(v)))
		}
	}

}

func tailingLogsInFileAll(fileName string, conn *websocket.Conn, current int64, page int) int {

	var statusPagination bool = false

	fileN := filepath.Base(fileName)
	//logenc.DeleteOldsFiles("./web/util/replace/"+fileN, "")
	original, err := os.Open(fileName)
	if err != nil {
		log.Println(err)
	}
	exec.Command("/bin/bash", "-c", "echo > "+"./web/util/replace/"+fileN).Run()
	logenc.CopyFile("./web/util/replace/", fileN, original)
	var countline int = 0
	taillog, err := tail.TailFile("./web/util/replace/"+fileN,
		tail.Config{
			ReOpen: true,
			Follow: true,
			Location: &tail.SeekInfo{
				//Offset: current,
				Whence: io.SeekStart, //!!!

			},
		})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error occurred in opening the file: ", err)
		return countline
	}
	go taillog.StopAtEOF()
	//conn.WriteMessage(websocket.TextMessage, []byte("<start></start>"))

	for line := range taillog.Lines {
		if page != 0 && line.Text != "" && line.Text != " " {

			pagUlid := paginationUlids[page]
			csvsimpl := logenc.ProcLineDecodeXML(line.Text)
			currentUlid := csvsimpl.XML_RECORD_ROOT[0].XML_ULID
			if pagUlid == currentUlid {
				statusPagination = true
			}
			if statusPagination {
				countline++
				conn.WriteMessage(websocket.TextMessage, []byte(logenc.EncodeXML(csvsimpl)))
			}
		} else {
			csvsimpl := logenc.ProcLineDecodeXML(line.Text)
			countline++
			conn.WriteMessage(websocket.TextMessage, []byte(logenc.EncodeXML(csvsimpl)))
			//logenc.DeleteOldsFiles("./web/util/replace/"+fileN, "")
		}

		if countline == 100 {
			//taillog.StopAtEOF()
			//logenc.DeleteOldsFiles("./web/util/replace/"+fileN, "")

			return countline
		}

	}
	//logenc.DeleteOldsFiles("./web/util/replace/"+fileN, "")

	return countline

}

func followCodeStatus(conn *websocket.Conn) {
	//Reset(conn)
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err, "followCodeStatus")
			return
		}
		fmt.Println("msgType", msgType)
		page, err = strconv.Atoi(string(msg[:]))
		if err != nil {
			currentfile = string(msg)
		}
		fmt.Println("Page", page)
	}

}

func UlidPaginationFile(conn *websocket.Conn, fileName string) {
	var CountPage string
	paginationUlids = make(map[int]string)
	var (
		strSlice []string

		countline int
		page      int    = 0
		firstUlid string = " "
	)
	taillog, err := tail.TailFile(fileName,
		tail.Config{
			Follow: true,
			Location: &tail.SeekInfo{
				Whence: io.SeekStart, //!!!

			},
		})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error occurred in opening the file: ", err)
		return
	}
	for line := range taillog.Lines {
		strSlice = append(strSlice, logenc.ProcLineDecodeXMLUlid(line.Text))
		countline++
		if countline == 100 {
			page++
			countline = 0
			firstUlid = strSlice[0]
			paginationUlids[page] = firstUlid
			strSlice = nil

		}
		go taillog.StopAtEOF() //end tail and stop service
	}
	page++
	CountPage = "<countpage>" + strconv.Itoa(page) + "</countpage>"
	conn.WriteMessage(websocket.TextMessage, []byte(CountPage))
	firstUlid = strSlice[1]
	paginationUlids[page] = firstUlid

	for key, value := range paginationUlids {
		fmt.Println("Key:", key, "Value:", value)
	}

	fmt.Println("map", (paginationUlids))
	countline = 0
	strSlice = nil

}

func UlidPaginationDir(conn *websocket.Conn, countFiles int, fileList map[string][]string) {
	//var CountPage string
	var CountPage string
	paginationUlids = make(map[int]string)
	var (
		strSlice []string

		countline int
		page      int    = 1
		firstUlid string = " "
	)

	//fileList["FileList"] = util.Conf.Dir

	for i := 0; i < countFiles; i++ {
		fileName := fileList["FileList"][i]
		taillog, err := tail.TailFile(fileName,
			tail.Config{
				Follow: true,
				Location: &tail.SeekInfo{
					Whence: io.SeekStart, //!!!

				},
			})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error occurred in opening the file: ", err)
			return
		}
		for line := range taillog.Lines {
			strSlice = append(strSlice, logenc.ProcLineDecodeXMLUlid(line.Text))
			countline++
			if countline == 100 {
				page++
				countline = 0
				firstUlid = strSlice[1]
				paginationUlids[page] = firstUlid
				strSlice = nil

			}
			go taillog.StopAtEOF() //end tail and stop service

		}
		if countline != 0 && countline < 100 && page == 0 {
			page++
			firstUlid = strSlice[1]
			paginationUlids[page] = firstUlid
		}
	}
	CountPage = "<countpage>" + strconv.Itoa(page) + "</countpage>"
	conn.WriteMessage(websocket.TextMessage, []byte(CountPage))
	firstUlid = strSlice[1]
	paginationUlids[page] = firstUlid
	for key, value := range paginationUlids {
		fmt.Println("Key:", key, "Value:", value)
	}

	fmt.Println("map", (paginationUlids))
	countline = 0
	strSlice = nil
}
func (m Map) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, v := range m {
		e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: v})
	}

	return e.EncodeToken(start.End())
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
		if !v {
			delete(indexMap, k)
			continue
		}
		//fmt.Fprintln(os.Stderr, k)
		FileList = append(FileList, k)
	}
	//filepath.Base
	//filename
	for _, f := range FileList {
		fileN := filepath.Base(f)
		FileName = append(FileName, fileN)
	}
	//fmt.Println("FileName   ", FileName)
	//fmt.Println("FileList   ", FileList)
	Conf.Dir = FileList
	Conf.Dir1 = FileName
	FileName = nil
	fmt.Fprintln(os.Stderr, "Indexing complete !, file index length: ", len(Conf.Dir))
	//fmt.Println(Conf.Dir)
	return nil
}

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
			//dfs(basepath + string(os.PathSeparator) + f.Name())
		}
	} else if strings.ContainsAny(s.Mode().String(), "alTLDpSugct") {
		// skip these files
		// try including names PIPES
	} else {
		// only remaining file are ascii files that can be then differentiated
		// by the user as golang has only these many categorization
		// Note : this appends the absolute paths
		// Insert the absPath into the Map, avoids duplicates in successive cron runs
		indexMap[absPath] = true
	}
}

func TailDir(conn *websocket.Conn, fileName string, lookFor string, SearchMap map[string]logenc.LogList, startUnixTime int64, endUnixTime int64) {
	fileN := filepath.Base(fileName)
	UlidC := bleveSI.ProcBleveSearchv2(fileN, lookFor)
	//println("Find", lookFor)
	//println(startUnixTime)
	//println(endUnixTime)
	//println(fileName)
	if (lookFor == "" || lookFor == " " || lookFor == "Search") && (startUnixTime == 0 || endUnixTime == 0) {

		go followCodeStatus(conn)
		var countline int = 0
		var currentpage int = 0

		logenc.DeleteOldsFiles("./web/util/replace/"+fileN, "")
		conn.WriteMessage(websocket.TextMessage, []byte("<start></start>"))
		countline = tailingLogsInFileAll(fileName, conn, 0, page)
		logenc.DeleteOldsFiles("./web/util/replace/"+fileN, "")
		for {

			/* if currentfile != fileN {

				break
			} else  */

			if countline >= 99 {
				//TransmitUlidPagination(conn, fileName)
				countline = 0
			} else if currentpage != page {
				logenc.DeleteOldsFiles("./web/util/replace/"+fileN, "")
				conn.WriteMessage(websocket.TextMessage, []byte("<start></start>"))
				countline = tailingLogsInFileAll(fileName, conn, 0, page)
				fmt.Println("countline", countline)
				logenc.DeleteOldsFiles("./web/util/replace/"+fileN, "")
				currentpage = page
				if countline < 95 {
					break
				}

			}

		}

		//conn.WriteMessage(websocket.TextMessage, []byte(logenc.EncodeXML(commoncsv)))
		//commoncsv = logenc.LogList{}

	} else if len(UlidC) == 0 {
		println("Break")
		return
	} else if (startUnixTime != 0 || endUnixTime != 0) && (lookFor == "" || lookFor == " " || lookFor == "Search") {
		for i := 0; i < len(UlidC); i++ {
			ulidS := UlidC[i]
			ulidU, _ := ulid.Parse(ulidS)
			Unixtime := ulidU.Time()
			log.Println("Uint64", Unixtime)
			log.Println("int64", int64(Unixtime))
			if startUnixTime <= int64(Unixtime) && endUnixTime >= int64(Unixtime) {
				conn.WriteMessage(websocket.TextMessage, []byte(ulidS))
			}

		}
	} else {
		var commoncsv logenc.LogList
		for i := 0; i < len(UlidC); i++ {
			//	ulidS := UlidC[i]
			//	ulidU, _ := ulid.Parse(ulidS)
			//	Unixtime := ulidU.Time()
			v, found := SearchMap[UlidC[i]]
			log.Println(v)
			fmt.Println(v)
			//	if found && startUnixTime <= int64(Unixtime) && endUnixTime >= int64(Unixtime) {
			if found {
				//:TODO create common structure
				//PS: Merge xml structure
				//:TODO map with xml structure
				//structure <loglist> append <log></log>......<log></log></loglist>
				commoncsv.XML_RECORD_ROOT = append(commoncsv.XML_RECORD_ROOT, v.XML_RECORD_ROOT...)
			}
		}
		return
		//:TODO transmit to websoket

	}

}

func GetFiles(address string, port string) error {
	//var signature bool = false
	resp, err := http.Get("http://" + address + ":" + port + "/vfs/data/")
	if err != nil {

		return err
		//log.Fatal(err)

	}
	for _, v := range logenc.GetLinks(resp.Body) {
		fmt.Println(address)

		fullURLFile := "http://" + address + ":" + port + "/vfs/data/" + v

		fileURL, err := url.Parse(fullURLFile)
		if err != nil {

			log.Fatal("Parse", err)
		}
		path := fileURL.Path
		segments := strings.Split(path, "/")
		fileName := segments[len(segments)-1]

		func() { // lambda for defer file.Close()
			file, err := os.OpenFile("./testsave/"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {

				log.Fatal("Getfiles", err)
				//file.Close()
				//return
			}

			defer file.Close()

			client := http.Client{
				CheckRedirect: func(r *http.Request, _ []*http.Request) error {
					r.URL.Opaque = r.URL.Path
					return nil
				},
			}
			// Put content on file
			resp, err := client.Get(fullURLFile)
			if err != nil {

				logenc.DeleteOldsFiles("./testsave/"+fileName, "")
				return
				//log.Fatal(err)
			}
			defer resp.Body.Close()
			contain := strings.Contains(fileName, "md5")
			if contain && logenc.CheckFileSum("./testsave/"+fileName, "rep", "") {
				signature = true

				fileS, err := os.OpenFile("./"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
				if err != nil {

					log.Println("Open for copy", err)
				}
				defer fileS.Close()
				_, err = io.Copy(fileS, resp.Body)
				if err != nil {

					log.Println("Copy", err)
				}
				logenc.WriteFileSum("./testsave/"+fileName, "rep", "")
				log.Println("*1")
				logenc.DeleteOldsFiles("./testsave/"+fileName, "")

			} else if !contain {
				_, err = io.Copy(file, resp.Body)
				if err != nil {

					log.Println("Copy", err)
				}
			}

			if signature && !contain {
				last3 := fileName[len(fileName)-3:]
				if logenc.CheckFileSum("./testsave/"+fileName, last3, "") {
					log.Println("*2")
					logenc.DeleteOldsFiles("./repdata/"+fileName, "")
					logenc.Replication("./testsave/" + fileName)
					logenc.WriteFileSum("./testsave/"+fileName, "rep", "")
					fmt.Println("Merge", fileName)

					log.Println("*3")
					logenc.DeleteOldsFiles("./testsave/"+fileName, "")

				} else {
					logenc.Replication("./testsave/" + fileName)
					logenc.WriteFileSum("./testsave/"+fileName, "rep", "")
					fmt.Println("Merge", fileName)
					log.Println("*4")
					logenc.DeleteOldsFiles("./testsave/"+fileName, "")
				}

			} else if !signature && !contain {
				//time.Sleep(15 * time.Second)
				logenc.Replication("./testsave/" + fileName)
				logenc.WriteFileSum("./testsave/"+fileName, "rep", "")
				fmt.Println("Merge", fileName)
				log.Println("*5")
				logenc.DeleteOldsFiles("./testsave/"+fileName, "")
			}

		}()
	}
	return nil
}

//Disk Check
type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

func DiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
	//limit big.Float = 0.8
)

func DiskInfo(dir string) {
	time.Sleep(time.Second * 55)
	for {
		disk := DiskUsage(dir)
		x, y := big.NewFloat(float64(disk.All)/float64(GB)), big.NewFloat(float64(disk.Used)/float64(GB))
		z := new(big.Float).Quo(y, x)

		if z.Cmp(big.NewFloat(0.8)) == 1 || z.Cmp(big.NewFloat(0.8)) == 0 {
			//fmt.Println("HAHA")ste
			FindOldestfile(dir)

		} else {
			//fmt.Println("HA")
			DeleteFile90(dir)
		}
		//fmt.Printf("All: %.2f GB\n", float64(disk.All)/float64(GB))
		//fmt.Printf("Used: %.2f GB\n", float64(disk.Used)/float64(GB))
		//fmt.Printf("Free: %.2f GB\n", float64(disk.Free)/float64(GB))
	}

}

func FindOldestfile(dir string) {
	var name string
	var cutoff = time.Hour
	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal("FindOldestfile", err.Error())
	}
	now := time.Now()
	//fmt.Println(now)
	for _, info := range fileInfo {
		if diff := now.Sub(info.ModTime()); diff > cutoff {
			cutoff = now.Sub(info.ModTime())
			name = info.Name()

		}
	}
	logenc.DeleteOldsFiles(dir+name, "")
}

func DeleteFile90(dir string) {

	var cutoff = 24 * time.Hour * 90
	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal("DeleteFile90", err.Error())
	}
	now := time.Now()
	//fmt.Println(now)
	for _, info := range fileInfo {
		//fmt.Println(info.Name())
		if diff := now.Sub(info.ModTime()); diff > cutoff {
			fmt.Printf("Deleting %s which is %s old\n", info.Name(), diff)
			logenc.DeleteOldsFiles(dir+info.Name(), "")

		}
	}

}

func CheckIPAddress(ip string) bool {
	/* if ip == "localhost" {
		fmt.Printf("IP Address: %s - Valid\n", ip)
		return true
	} else  */if net.ParseIP(ip) == nil {
		fmt.Printf("IP Address: %s - Invalid\n", ip)
		return false
	} else {
		fmt.Printf("IP Address: %s - Valid\n", ip)
		return true
	}

}
