# md5sum

This program displays md5 sums for files specified on the command line.

## Example (default output)

    $ md5sum DODO go.mod main.go
    DODO:    1f3e040bfa1a31380075c5c7331a7e03
    go.mod:  492ea9ddfb3a923dd696969338cd21a6
    main.go: f5ae2cc3718cce351da9f257fd11054b

## Options

  * -t (terse)

    Only output 'ok' if all the sums match or 'mismatch' if there are
    mismatches.

  * -q (quiet)

    No output is written to stdout. The program ends with exit(0) if the sum
    match or exit(1) if any of the sums fail to match.
