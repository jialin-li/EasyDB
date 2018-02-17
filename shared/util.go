package shared

import (
	"log"
	"net"
	"strconv"
)

func Dial(port int) (net.Conn, error) {
	conn, err := net.Dial("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		log.Println(err)
	}
	return conn, err
}
