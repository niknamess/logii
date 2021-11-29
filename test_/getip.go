package test

import (
	"fmt"
	"log"
	"net"
)

// Get preferred outbound ip of this machine
func GetOutboundIP() {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	str := localAddr.IP
	str1 := str.String()
	fmt.Println(str)
	fmt.Println(str1)
	///return localAddr.IP
}

/*
func createIndex() (bleve.Index, error) {
	indexName := "history.bleve"
	index, err := bleve.Open(indexName)
	if err == bleve.ErrorIndexPathDoesNotExist {
		mapping := buildMapping()
		kvStore := goleveldb.Name
		kvConfig := map[string]interface{}{
			"create_if_missing": true,
		}
		index, err = bleve.NewUsing(indexName, mapping, "upside_down", kvStore, kvConfig)
	}
	if err != nil {
		return err
	}
}
*/
