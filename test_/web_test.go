package test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gitlab.topaz-atcs.com/tmcs/logi2/web"
)

func TestProcWeb(t *testing.T) {

	port := 15001

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Yeah!"))
		//		fmt.Println("/")
		//		http.ServeFile(w, r, "index.tmpl")
	})

	csrfRouter := web.Use((router).ServeHTTP)

	server := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", port), Handler: handlers.CombinedLoggingHandler(os.Stdout, csrfRouter)}

	go func() {
		<-time.NewTimer(time.Second * 3).C

		server.Close()
		t.Log("server stop")

		<-time.NewTimer(time.Second * 1).C
	}()

	go func() {
		<-time.NewTimer(500 * time.Millisecond).C

		resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d", port))
		if err != nil {
			t.Error(err.Error())
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error(err.Error())
			return
		}
		t.Log(string(body))
	}()

	t.Log("server start")

	if err := server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			t.Error(err.Error())
		}
	}

} // else {
/*func TestWebPost(t *testing.T) {

	port := 15000

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Yeah!"))
		//		fmt.Println("/")
		//		http.ServeFile(w, r, "index.tmpl")
	})

	router.HandleFunc("/second", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("second yhu!"))
		//		fmt.Println("/")
		//		http.ServeFile(w, r, "index.tmpl")
	})
	router.HandleFunc("/testGet", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, nil)

		//		fmt.Println("/")
		//		http.ServeFile(w, r, "index.tmpl")
	})

	csrfRouter := Use((router).ServeHTTP)
	server := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", port), Handler: handlers.CombinedLoggingHandler(os.Stdout, csrfRouter)}

	go func() {
		<-time.NewTimer(time.Second * 3).C

		server.Close()
		t.Log("server stop")

		<-time.NewTimer(time.Second * 1).C
	}()

	go func() {
		<-time.NewTimer(500 * time.Millisecond).C

		for _, s := range []string{"/", "/second", "/testGet"} {
			resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d%s", port, s))
			if err != nil {
				t.Error(err.Error())
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err.Error())
				return
			}
			t.Log(string(body))
			fmt.Println(string(body))

		}

	}()

	t.Log("server start")

	if err := server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			t.Error(err.Error())
		}
	}

} // else {
*/
