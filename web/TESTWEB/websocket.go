package testweb

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// 1
//type longLatStruct struct {
//	Long float64 `json:"longitude"`
//	Lat  float64 `json:"latitude"`
//}

//file test
type longLatStruct struct {
	Id      int64  `json:"id"`      //Long
	Hashsum string `json:"hashsum"` //Lat
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan *longLatStruct)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Testweb() {
	// 2
	router := mux.NewRouter()
	router.HandleFunc("/", rootHandler).Methods("GET")
	router.HandleFunc("/longlat", longLatHandler).Methods("POST")
	router.HandleFunc("/ws", wsHandler)
	curl()
	go echo()
	//curl()
	log.Fatal(http.ListenAndServe(":8844", router))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "home")
}

func writer(coord *longLatStruct) {
	broadcast <- coord
}

func longLatHandler(w http.ResponseWriter, r *http.Request) {
	var coordinates longLatStruct
	if err := json.NewDecoder(r.Body).Decode(&coordinates); err != nil {
		log.Printf("ERROR: %s", err)
		http.Error(w, "Bad request", http.StatusTeapot)
		return
	}
	defer r.Body.Close()
	go writer(&coordinates)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// register client
	clients[ws] = true
}

func curl() {
	//hashsum := dop.FileMD5(file)
	params := url.Values{}
	params.Add("'{\"id\": 220122, \"hashsum\": aslkdplaks23}}'", ``)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "localhost:8844/longlat", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	//defer resp.Body.Close()

}

// 3
func echo() {
	for {
		curl()
		val := <-broadcast
		latlong := fmt.Sprintf("%f %f %s", val.Hashsum, val.Id)
		// send to every client that is currently connected
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(latlong))
			if err != nil {
				log.Printf("Websocket error: %s", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
