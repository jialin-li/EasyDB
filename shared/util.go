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

func (e *DepError) Error() string {
	return ERR_DEP
}

func (e *KeyError) Error() string {
	return ERR_KEY
}

// Compare the timestamp with t2, updating any times that are lower than the
// corresponding index in t2
func (t1 *Time) Update(t2 *Time) {
	for i := range t1.Clock {
		if t1.Clock[i] < t2.Clock[i] {
			t1.Clock[i] = t2.Clock[i]
		}
	}
}

func (t1 *Time) IsLaterThan(t2 *Time) bool {
	for i := range t1.Clock {
		if t1.Clock[i] > t2.Clock[i] {
			return true
		}
		if t1.Clock[i] < t2.Clock[i] {
			return false
		}
	}
	return false
}
