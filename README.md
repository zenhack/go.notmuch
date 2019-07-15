[![Build
Status](https://travis-ci.org/zenhack/go.notmuch.svg?branch=master)](https://travis-ci.org/zenhack/go.notmuch)

Go binding for [notmuch mail][1].

Licensed under the GPLv3 or later (like notmuch itself).

[1]: http://notmuchmail.org/

# Development

## Running tests
The project uses `make` to setup and download additional assets for the tests.

Run `make test` to run the tests.

## Pre PR checks
Next to the tests, you should also run gofmt on the sourcecode.
Running `make fmtcheck` checks for formatting issues.

To run both tests and format checks, use `make ci`.
