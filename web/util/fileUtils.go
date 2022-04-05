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
	//"gitlab.topaz-atcs.com/tmcs/logi2/web/controllers"
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
	countSearch int
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
	//paginationUlids = make(map[string]int64)
	fmt.Println("Fname", Fname)
	var (
		strSlice  []string
		firstUlid string = " "
	)
	fileN := filepath.Base(fileName)

	if Fname != fileName {
		if Fname != "" {
			logenc.DeleteOldsFiles("./web/util/replace/"+filepath.Base(Fname), "")
		}
		Fname = fileName
		lookFor = ""
		//current = 638
		countSearch = 0
		firstUlid = " "
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
		TransmitUlidPagination(conn, fileName)
		tailingLogsInFileAll(fileName, conn, 0, page)
		//var countline int = 0
		for {
			if (logenc.FileMD5(fileName) != hashSumFile) && currentfile == fileN {
				TransmitUlidPagination(conn, fileName)
				tailingLogsInFileAll(fileName, conn, 0, page)
				//TransmitUlidPagination(conn, fileName)
				hashSumFile = logenc.FileMD5(fileName)
				//currentpage = page
			} else if currentfile != fileN {
				//countline = 0
				break
			}

		}

		return

	} else if len(UlidC) == 0 {
		println("Break")
		return
	} else {
		var commoncsv logenc.LogList
		var countCheck int
		var count int = 0
		for i := 0; i < len(UlidC); i++ {

			_, found := SearchMap[UlidC[i]]

			if found {
				count++
			}

		}
		fmt.Println(count)
		fmt.Println("countSearch", countSearch)
		for i := countSearch; i < len(UlidC); i++ {
			v, found := SearchMap[UlidC[i]]
			if found {

				commoncsv.XML_RECORD_ROOT = append(commoncsv.XML_RECORD_ROOT, v.XML_RECORD_ROOT...)
				countCheck++

				strSlice = append(strSlice, v.XML_RECORD_ROOT[0].XML_ULID)
				if countCheck == 1000 {
					conn.WriteMessage(websocket.TextMessage, []byte(logenc.EncodeXML(commoncsv)))
					countCheck = 0
					commoncsv = logenc.LogList{}
					countSearch = i
					firstUlid = strSlice[0]
					fmt.Println("firstUlid", firstUlid)
					strSlice = nil
					return
				} else if len(UlidC) == i-1 && countCheck < 1000 {
					conn.WriteMessage(websocket.TextMessage, []byte(logenc.EncodeXML(commoncsv)))
					countCheck = 0
					commoncsv = logenc.LogList{}
					countSearch = 0
					firstUlid = strSlice[0]
					fmt.Println("firstUlid", firstUlid)
					strSlice = nil
					return
				}

			}

		}
		conn.WriteMessage(websocket.TextMessage, []byte(logenc.EncodeXML(commoncsv)))
		commoncsv = logenc.LogList{}

		return

	}

}

func tailingLogsInFileAll(fileName string, conn *websocket.Conn, current int64, page int) int {

	fileN := filepath.Base(fileName)
	logenc.DeleteOldsFiles("./web/util/replace/"+fileN, "")
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
	conn.WriteMessage(websocket.TextMessage, []byte("<start></start>"))

	for line := range taillog.Lines {
		//taillog.StopAtEOF()
		/* testinfo, err := taillog.Tell()
		fmt.Println("Taill............", testinfo)

		if err != nil {
			taillog.Stop()
		}
		if currentfile != fileN {
			taillog.Stop()
			return countline
		}

		log.Println("File change ", logenc.FileMD5(fileName))
		if logenc.FileMD5(fileName) != hashSumFile {

			inf, _ := taillog.Tell()
			log.Println("File change ", "OLD:", hashSumFile, "NEW:", logenc.FileMD5(fileName), "current tail:", inf)
			//hashSumFile = logenc.FileMD5(fileName)
			conn.WriteMessage(websocket.TextMessage, []byte("<start></start>"))
			countline = 0
			taillog.Stop()
			logenc.DeleteOldsFiles("./web/util/replace/"+fileN, "")
			taillog.Cleanup()
			return countline */

		/*}  else if page != 0 {
		//current, _ = taillog.Tell()
		//fmt.Println("---------", current)
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
		if countline == 510 {
			taillog.Stop()
			return countline
		} */
		//	}
		//	conn.WriteMessage(websocket.TextMessage, []byte("<start></start>"))
		csvsimpl := logenc.ProcLineDecodeXML(line.Text)
		countline++
		conn.WriteMessage(websocket.TextMessage, []byte(logenc.EncodeXML(csvsimpl)))
		logenc.DeleteOldsFiles("./web/util/replace/"+fileN, "")

		if countline == 510 {
			taillog.Stop()
			logenc.DeleteOldsFiles("./web/util/replace/"+fileN, "")

			return countline
		}
		/*  */

		//taillog.StopAtEOF()

	}
	logenc.DeleteOldsFiles("./web/util/replace/"+fileN, "")

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
		//fmt.Println("msg", string(msg[:]))
		//fmt.Println(msg)
		page, err = strconv.Atoi(string(msg[:]))
		if err != nil {
			currentfile = string(msg)
		}
		fmt.Println("Page", page)
	}
	//code, _ = strconv.Atoi(string(msg[:]))

}

func TransmitUlidPagination(conn *websocket.Conn, fileName string) {
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

		//current, _ = taillog.Tell()

		strSlice = append(strSlice, logenc.ProcLineDecodeXMLUlid(line.Text))
		countline++
		if countline == 50 {
			page++
			countline = 0
			firstUlid = strSlice[1]
			//strconv.Itoa(page)
			//	paginationUlids[strconv.Itoa(page)] = ir_table{ulid: firstUlid, point: current}
			paginationUlids[page] = firstUlid
			strSlice = nil

		}
		go taillog.StopAtEOF() //end tail and stop service
	}
	page++
	CountPage = "<countpage>" + strconv.Itoa(page) + "</countpage>"
	conn.WriteMessage(websocket.TextMessage, []byte(CountPage))
	firstUlid = strSlice[1]
	//paginationUlids[logenc.Convert1to1000(page)] = firstUlid
	paginationUlids[page] = firstUlid

	//x, _ := xml.MarshalIndent(Map(paginationUlids), " ", "  ")
	//fmt.Println(string(x))
	//conn.WriteMessage(websocket.TextMessage, []byte(string(x)))
	for key, value := range paginationUlids {
		fmt.Println("Key:", key, "Value:", value)
	}

	fmt.Println("map", (paginationUlids))
	//fmt.Println("func", createKeyValuePairs(paginationUlids))
	countline = 0
	strSlice = nil
	//taillog.Stop()

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

/* func createKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%v=\"%v\"\n", key, value)
	}
	return b.String()
} */

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
	println("Find", lookFor)
	println(startUnixTime)
	println(endUnixTime)
	println(fileName)
	if (lookFor == "" || lookFor == " " || lookFor == "Search") && (startUnixTime == 0 || endUnixTime == 0) {
		var (
			countline int = 0
			commoncsv logenc.LogList
		)
		for line := range taillog.Lines {
			countline++
			//Найти lastUlid и только потом продолжать

			csvsimpl := logenc.ProcLineDecodeXML(line.Text)
			commoncsv.XML_RECORD_ROOT = append(commoncsv.XML_RECORD_ROOT, csvsimpl.XML_RECORD_ROOT...)
			if countline == 1000 {
				conn.WriteMessage(websocket.TextMessage, []byte(logenc.EncodeXML(commoncsv)))
				countline = 0
				commoncsv = logenc.LogList{}
			}

			go taillog.StopAtEOF() //end tail and stop service

		}
		conn.WriteMessage(websocket.TextMessage, []byte(logenc.EncodeXML(commoncsv)))
		commoncsv = logenc.LogList{}
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
	//return from Ulid timestamp
	//if in period return
	/* startUnixTime
	endUnixTime
	*/
	/* fileN := filepath.Base(fileName)
	UlidC := bleveSI.ProcBleveSearchv2(fileN, lookFor)

	if len(UlidC) == 0 {
		println("Break")
		return false
	} else {

		for i := 0; i < len(UlidC); i++ {

			_, found := SearchMap[UlidC[i]]
			if found {
				conn.WriteMessage(websocket.TextMessage, []byte("LOL"))
				return true

			}
		}

	}
	return false */

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

//Check 30 seconds file chancge
//:TODO
/* func FindNewChangeFile(dir string) bool {

	var cutoff = 30 * time.Second
	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal("FindNewChangeFile", err.Error())
	}
	now := time.Now()
	//fmt.Println(now)
	for _, info := range fileInfo {
		//fmt.Println(info.Name())
		if diff := now.Sub(info.ModTime()); diff > cutoff {
			fmt.Printf("Deleting %s which is %s old\n", info.Name(), diff)
			//logenc.DeleteOldsFiles(dir+info.Name(), "")
			return true
		}
	}

	return false

} */

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

/*
func lineCounter(path string) (int, error) {

	r, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v", err)
		os.Exit(1)
	}
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
*/
