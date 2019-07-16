[![Build Status][ci-img]][ci]

Go binding for [notmuch mail][notmuch].

Licensed under the GPLv3 or later (like notmuch itself).

# Development

## Running tests
The project uses `make` to setup and download additional assets for the tests.

Run `make test` to run the tests.

## Pre PR checks
Next to the tests, you should also run gofmt on the sourcecode.
Running `make fmtcheck` checks for formatting issues.

To run both tests and format checks, use `make ci`.

[notmuch]: http://notmuchmail.org/
[ci-img]: https://gitlab.com/isd/go-notmuch/badges/master/build.svg
[ci]: https://gitlab.com/isd/go-notmuch/pipelines
