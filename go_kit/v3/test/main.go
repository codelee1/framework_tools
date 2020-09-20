package main

import (
	"fmt"
	"hash/crc32"
)

func main() {
	fmt.Println(crc32.ChecksumIEEE([]byte("127.0.0.1:8080")))
	fmt.Println(crc32.ChecksumIEEE([]byte("127.0.0.1:8081")))
}
