package main

import (
    "crypto/md5"
    "fmt"
    "io"
    "os"

    U "../swak"
)

type md5Args struct {
    maxFnLen int
    terse bool
    quiet bool
    testing bool
}

var args md5Args

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

// ----------------------------------------------------------------------------
func handleArgs() {
    for _, arg := range os.Args[1:] {
        switch {
        case arg == "-T": args.testing = true
        case arg == "-t": args.terse = true
        case arg == "-q": args.quiet = true
        case arg[0] != '-':
            if args.maxFnLen < len(arg) {
                args.maxFnLen = len(arg)
            }
        }
    }
}

// ----------------------------------------------------------------------------
func md5sum(filepath string) string {
    f, err := os.Open(filepath)
    if err != nil {
        U.Fatalf("%v", err)
    }
    defer f.Close()

    h := md5.New()
    if _, err := io.Copy(h, f); err != nil {
        U.Fatalf("%v", err)
    }
    return fmt.Sprintf("%x", h.Sum(nil))
}
