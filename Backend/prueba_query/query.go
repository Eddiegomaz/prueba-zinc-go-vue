package main

import (
	"Project/zincShare"
	"fmt"
)

func main() {
	p, err := zincShare.Query("", 0, 10)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(p)
}
