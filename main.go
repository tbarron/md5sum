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

var output io.Writer = os.Stdout

// ----------------------------------------------------------------------------
func main() {
    var sumList []string
    handleArgs()
    for _, filepath := range os.Args[1:] {
        if filepath[0] == '-' {
            continue
        }
        sum := md5sum(filepath)
        sumList = append(sumList, sum)
        filepath += ":"
        fmt.Fprintf(output, "%*s %s\n",
            -1 * (args.maxFnLen + 1), filepath, sum)
    }
    xval := 0
    prev := sumList[0]
    for _, sum := range sumList {
        if prev != sum {
            xval = 1
        }
    }
    if ! args.testing {
        os.Exit(xval)
    }
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
