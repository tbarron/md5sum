package main

import (
    "crypto/md5"
    "fmt"
    "io"
    "log"
    "os"
)

func main() {
	for _, filepath := range os.Args[1:] {
		f, err := os.Open(filepath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		h := md5.New()
		if _, err := io.Copy(h, f); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s: %x\n", filepath, h.Sum(nil))
	}
}
