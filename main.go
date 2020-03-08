package main

import (
    "crypto/md5"
    "fmt"
    "io"
    "os"
)

type md5Args struct {
    maxFnLen int
    terse bool
    quiet bool
    testing bool
}

type report struct {
    filepath string
    sum      string
}

var args md5Args

var output io.Writer = os.Stdout

// ----------------------------------------------------------------------------
func main() {
    var sumList []string
    var result report
    var resList []report

    handleArgs()
    for _, filepath := range os.Args[1:] {
        if filepath[0] == '-' {
            continue
        }
        sum := md5sum(filepath)
        sumList = append(sumList, sum)
        filepath += ":"
        result.filepath = fmt.Sprintf("%*s",
            -1 * (args.maxFnLen + 1), filepath)
        result.sum = sum
        resList = append(resList, result)
    }
    xval := 0
    for i, res := range resList {
        if !args.terse && !args.quiet {
            fmt.Fprintf(output, "%s %s\n", res.filepath, res.sum)
        }
        if (0 < i) && (res.sum != resList[i-1].sum) {
            xval = 1
        }
    }
    if args.terse {
        if xval == 0 {
            fmt.Fprintf(output, "ok\n")
        } else {
            fmt.Fprintf(output, "mismatch\n")
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
        case arg == "-t":
            if args.quiet {
                fmt.Fprintf(output, "-t and -q are not compatible\n")
            }
            args.terse = true
        case arg == "-q":
            if args.terse {
                fmt.Fprintf(output, "-q and -t are not compatible\n")
            }
            args.quiet = true
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
        fatalf("%v", err)
    }
    defer f.Close()

    h := md5.New()
    if _, err := io.Copy(h, f); err != nil {
        fatalf("%v", err)
    }
    return fmt.Sprintf("%x", h.Sum(nil))
}

// ----------------------------------------------------------------------------
// Put out a message and exit
func fatalf(sfmt string, items ...interface{}) {
	sfmt += "\n"
	fmt.Fprintf(output, sfmt, items...)
	if output == os.Stdout {
		os.Exit(0)
	}
}
