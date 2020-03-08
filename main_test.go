package main

import (
    "io"
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
        expout string
    }{
        { "foo", md5Args{ 0, false, false, true }, "" },
        { "foo -t", md5Args{0, true, false, true }, "" },
        { "foo -q", md5Args{0, false, true, true }, "" },
        { "foo three fifteenish -t", md5Args{10, true, false, true }, "" },
        { "foo one three four fifteen", md5Args{7, false, false, true }, "" },
        { "foo one four -t -q", md5Args{4, true, true, true},
            "-q and -t are not compatible\n"},
        { "foo one four -q -t", md5Args{4, true, true, true},
            "-t and -q are not compatible\n" },
    }

    for _, c := range cases {
        os.Args = strings.Split(c.inp, " ")
        obuf := &strings.Builder{}
        setOutput(obuf)
        handleArgs()
        setOutput(os.Stdout)
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
        if obuf.String() != c.expout {
            t.Errorf("Expected output %q, got %q", c.expout, obuf.String())
        }
        args.maxFnLen = 0
        args.terse = false
        args.quiet = false
    }
}

// ----------------------------------------------------------------------------
func TestMain(t *testing.T) {
    obuf := &strings.Builder{}
    setOutput(obuf)
    os.Args = []string{"foo", "-q", "-T", "main.go", "main_test.go"}
    main()
    setOutput(os.Stdout)
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
func TestQuiet(t *testing.T) {
    result, err := externalTest(t, "md5sum -q main.go main_test.go")
    if err.Error() != "exit status 1" {
        t.Errorf("Expected error 'exit status 1', got '%s'", err.Error())
    }
    if result != "" {
        t.Errorf("Expected output \"\", got %q", result)
    }

    result, err = externalTest(t, "md5sum -q main.go main.go")
    if err != nil {
        t.Errorf("Expected error nil, got '%s'", err.Error())
    }
    if result != "" {
        t.Errorf("Expected output \"\", got %q", result)
    }
}

// ----------------------------------------------------------------------------
func TestTerse(t *testing.T) {
    result, err := externalTest(t, "md5sum -t main.go main_test.go")
    if err.Error() != "exit status 1" {
        t.Errorf("Expected error 'exit status 1', got '%s'", err.Error())
    }
    if result != "mismatch\n" {
        t.Errorf("Expected output \"mismatched\\n\", got %q", result)
    }

    result, err = externalTest(t, "md5sum -t main.go main.go")
    if err != nil {
        t.Errorf("Expected error nil, got '%s'", err.Error())
    }
    if result != "ok\n" {
        t.Errorf("Expected output \"ok\\n\", got %q", result)
    }
}

// ----------------------------------------------------------------------------
func externalTest(t *testing.T, cmd string) (string, error) {
    words := strings.Fields(cmd)
    xbl := exec.Command(words[0], words[1:]...)
    op, err := xbl.Output()
    return string(op), err
}

// ----------------------------------------------------------------------------
func setOutput(op io.Writer) {
    output = op
}
