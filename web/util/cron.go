package util

import (
	"fmt"
<<<<<<< HEAD
	"log"
	"net"
=======
>>>>>>> 7cc21d8bc26936e7ef731a8b2d1dc24da8cf5e15
	"os"
	"time"
)

// MakeAndStartCron - Creates a ticket with an interval of 'repeat' and pushes
// in a channel being read by the for loop in the function. Every time a value is
// pushed, the Cron executes the function passed
func MakeAndStartCron(repeat time.Duration, run func(...interface{}) error, v ...interface{}) {
	ticker := time.Tick(repeat)
	for _ = range ticker {
		fmt.Fprintf(os.Stderr, "Running cron job @%v\n", time.Now())
		//fmt.Println("length of arg :", len(v))
		err := run(v...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cron job failed: %s\n", err)
		}
	}
}
<<<<<<< HEAD

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	str := localAddr.IP
	str1 := str.String()
	//fmt.Println(str)
	//fmt.Println(str1)
	///return localAddr.IP
	return str1
}
=======
>>>>>>> 7cc21d8bc26936e7ef731a8b2d1dc24da8cf5e15
