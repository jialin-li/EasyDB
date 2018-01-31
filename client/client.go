package main

import "fmt"

func main() {
	client(5)
}
func client(id int) int {
	for i := 0; i < id; i++ {
		fmt.Println("in client")
	}
	return 0

}
