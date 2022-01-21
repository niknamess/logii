package web

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/gorilla/mux"
	"github.com/pelletier/go-toml"
	"github.com/spf13/viper"
	"gitlab.topaz-atcs.com/tmcs/logi2/generate_logs"
	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
	"gitlab.topaz-atcs.com/tmcs/logi2/web/controllers"
	"gitlab.topaz-atcs.com/tmcs/logi2/web/util"
)

var (
	//dir = kingpin.Arg("dir", "Directory path(s) to look for files").Default("./view").ExistingFilesOrDirs()
	//dir = kingpin.Arg("dir", "Directory path(s) to look for files").Default("/home/nik/projects/Course/logi2/repdata/").ExistingFilesOrDirs()
	//dir = kingpin.Arg("dir", "Directory path(s) to look for files").Default("/home/nik/projects/Course/tmcs-log-agent-storage/").ExistingFilesOrDirs()
	dir   = kingpin.Arg("dir", "Directory path(s) to look for files").Default("./repdata").ExistingFilesOrDirs()
	port  = kingpin.Flag("port", "Port number to host the server").Short('p').Default("15000").Int()
	portx = kingpin.Flag("portx", "Port number to host the server").Short('x').Default("15000").Int()
	ports = kingpin.Flag("ports", "Port number to host the server").Short('s').Default("15000").Int()
	//port            *int
	cron            = kingpin.Flag("cron", "configure cron for re-indexing files, Supported durations:[h -> hours, d -> days]").Short('t').Default("0h").String()
	cert            = kingpin.Flag("Test", "Test").Short('c').Default("").String()
	missadr         []string
	limit           string
	ipaddr          []string
	wg              sync.WaitGroup
	status          bool = false
	ctxCF, cancelCF      = context.WithCancel(context.Background())
)

type DatabaseConfig struct {
	Host  []string `mapstructure:"hostname"`
	Hostt []string `toml:"hostname"`
	Port  string
}

type Config struct {
	Db       DatabaseConfig `mapstructure:"database"`
	DataBase DatabaseConfig `toml:"database"`
}

func ProcWeb(dir1 string, slice []string, ctx context.Context) (err error) {
	status = false
	if dir1 == "-x" {
		status = true
	}

	fmt.Println("dir", dir1)

	generate_logs.Remove("./testsave/", "gen_logs_coded")
	generate_logs.Remove("./testsave/", "md5")

	logenc.CreateDir("./repdata/", "")
	logenc.CreateDir("./testsave/", "")

	kingpin.Parse()

	err = util.ParseConfig(*dir, *cron, *cert) //INDEXING FILE

	if err != nil {
		panic(err)

	}
	//go util.DeleteFile90("./repdata")
	go func() {
		time.Sleep(time.Second * 55)
		util.DiskInfo("./repdata")
	}()
	if status {
		EnterIpReady(slice)
	} else {
		EnterIp()
	}

	Ip, CPort := CheckConfig()
	for i := 0; i < len(Ip); i++ {
		//TODO
		go CheckFiles(Ip[i], CPort, ctxCF)
		time.Sleep(time.Second * 2)
	}

	go CheckFiles("localhost", "10015", ctxCF)

	time.Sleep(time.Second * 10)

	router := mux.NewRouter()
	router.HandleFunc("/ws/{b64file}", Use(controllers.WSHandler)).Methods("GET")
	router.HandleFunc("/", Use(controllers.RootHandler)).Methods("GET")
	router.HandleFunc("/searchproject", controllers.SearchHandler)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/static")))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.tmpl")
	})
	fmt.Println(port)
	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", 15000), //port
		Handler: router}
	go func() {
		if err = server.ListenAndServe(); err != nil {
			log.Println("listen:", err)
		}

	}()
	<-ctx.Done()
	go func() {
		cancelCF()
		ctxCF, cancelCF = context.WithCancel(context.Background())
	}()

	log.Printf("server stopped")

	if err = server.Shutdown(context.Background()); err != nil {
		log.Fatalf("server Shutdown Failed:%+s", err)
	}

	log.Printf("server exited properly")
	if err == http.ErrServerClosed {
		err = nil
	}
	return err
}

// Use - Stacking middlewares
func Use(handler http.HandlerFunc, mid ...func(http.Handler) http.HandlerFunc) http.HandlerFunc {
	for _, m := range mid {
		handler = m(handler)
	}
	return handler
}

func CheckFiles(address string, port string, ctx context.Context) {

	for range time.Tick(time.Second * 3) {
		if len(missadr) == 0 {
			missadr = append(missadr, "nope")
		}
		for {
			select {
			case <-ctx.Done():
				fmt.Println("stop CheckFiles")
				return
			default:
				for i := range missadr {
					if missadr[i] != address {
						err := util.GetFiles(address, port)
						if err != nil {
							log.Println(err)
							fmt.Println(address)
							missadr = append(missadr, address)
						}
						err = util.ParseConfig(*dir, *cron, *cert) //INDEXING FILE

						if err != nil {
							log.Println("LOOP", err)
							panic(err)
						}

					}
				}

				fmt.Println(missadr)
				wg.Add(1)
				go reconect(address)
				wg.Wait()
				time.Sleep(time.Second * 5)
				continue

			}
		}

	}

}

//FIXME never used

func reconect(address string) {
	defer wg.Done()

	for i := range missadr {
		if missadr[i] == address {
			copy(missadr[i:], missadr[i+1:]) // Shift a[i+1:] left one index.
			missadr[len(missadr)-1] = ""     // Erase last element (write zero value).
			missadr = missadr[:len(missadr)-1]
		}

	}

}

func EnterIp() {
	var data []byte
	for {
		fmt.Print("Enter IP:  ")
		fmt.Scanln(&limit)

		if limit == "stop" {
			break
		} else if util.CheckIPAddress(limit) {
			ipaddr = append(ipaddr, limit)
			limitSlice, _ := CheckConfig()
			ipaddr = append(ipaddr, limitSlice...)
			ipaddr = removeDuplicateStr(ipaddr)
			config := Config{DataBase: DatabaseConfig{Hostt: ipaddr, Port: "10015"}}
			data, _ = toml.Marshal(&config)
		}
	}
	//TODO
	err3 := ioutil.WriteFile("config.toml", data, 0666)

	if err3 != nil {

		log.Fatal(err3)
	}
	fmt.Println("Written")

}

func CheckConfig() ([]string, string) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		fmt.Println("couldn't load config:", err)
		os.Exit(1)
	}
	var c Config
	if err := v.Unmarshal(&c); err != nil {
		fmt.Printf("couldn't read config: %s", err)
	}
	Ip := c.Db.Host
	Port := c.Db.Port
	Ip = removeDuplicateStr(Ip)

	return Ip, Port
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func EnterIpReady(ipmas []string) {
	var data []byte
	ipaddr = ipmas
	limitSlice, _ := CheckConfig()
	ipaddr = append(ipaddr, limitSlice...)
	ipaddr = removeDuplicateStr(ipaddr)
	config := Config{DataBase: DatabaseConfig{Hostt: ipaddr, Port: "10015"}}
	data, _ = toml.Marshal(&config)

	//TODO
	err3 := ioutil.WriteFile("config.toml", data, 0666)

	if err3 != nil {

		log.Fatal(err3)
	}
	fmt.Println("Written")

}
