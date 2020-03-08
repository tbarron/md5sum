package main

import (
    "os"
    "os/exec"
    "strings"
    "testing"
)

// ----------------------------------------------------------------------------
func TestExtAlignment(t *testing.T) {
    frog, err := externalTest(t, "md5sum main.go main_test.go")
    x := strings.Split(frog, "\n")
    exp := "main.go:      "
    if x[0][0:len(exp)] != exp {
        t.Errorf("%q != %q", x[0][0:len(exp)], exp)
    }
    if err.Error() != "exit status 1" {
        t.Errorf("Expected error 'exit status 1', got '%s'", err.Error())
    }
}

// ----------------------------------------------------------------------------
func TestHandleArgs(t *testing.T) {
    cases := []struct {
        inp string
        exp md5Args
    }{
        { "foo", md5Args{ 0, false, false, true } },
        { "foo -t", md5Args{0, true, false, true } },
        { "foo -q", md5Args{0, false, true, true } },
        { "foo three four fifteenish -q -t", md5Args{10, true, true, true } },
        { "foo one two three four fifteen", md5Args{7, false, false, true } },
    }

    for _, c := range cases {
        os.Args = strings.Split(c.inp, " ")
        handleArgs()
        if args.maxFnLen != c.exp.maxFnLen {
            t.Errorf("Expected %d in maxFnLen, got %d",
                c.exp.maxFnLen, args.maxFnLen)
        }
        if args.terse != c.exp.terse {
            t.Errorf("Expected %v in terse, got %v", c.exp.terse, args.terse)
        }
        if args.quiet != c.exp.quiet {
            t.Errorf("Expected %v in quiet, got %v", c.exp.quiet, args.quiet)
        }
        args.maxFnLen = 0
        args.terse = false
        args.quiet = false
    }
}

// ----------------------------------------------------------------------------
func TestMain(t *testing.T) {
    os.Args = []string{"foo", "-q", "-T", "main.go", "main_test.go"}
    main()
}

// ----------------------------------------------------------------------------
func TestMd5sum(t *testing.T) {
    exp := "d41d8cd98f00b204e9800998ecf8427e"
    result := md5sum("/dev/null")
    if result != exp {
        t.Errorf("Expected %q, got %q", exp, result)
    }
}

// ----------------------------------------------------------------------------
func TestMd5sumNosuch(t *testing.T) {
    exp := "open ./nosuchfile: no such file or directory\n"
    result, err := externalTest(t, "md5sum ./nosuchfile")
    if err != nil {
        t.Errorf("Failure: %q", err.Error())
    }
    if result != exp {
        t.Errorf("Expected %q, got %q", exp, result)
    }
}

// ----------------------------------------------------------------------------
func externalTest(t *testing.T, cmd string) (string, error) {
    words := strings.Fields(cmd)
    xbl := exec.Command(words[0], words[1:]...)
    op, err := xbl.Output()
    return string(op), err
}
