# SQL dump to XLSX
[![CI](https://github.com/ledoyen/sqldump-to-xlsx/actions/workflows/ci.yml/badge.svg)](https://github.com/ledoyen/sqldump-to-xlsx/actions)
[![codecov](https://codecov.io/gh/ledoyen/sqldump-to-xlsx/branch/main/graph/badge.svg?token=FLEFJVYXM1)](https://codecov.io/gh/ledoyen/sqldump-to-xlsx)

Transform a Database dump into an XLSX file

Tested on rather little dumps (~6Mb), it takes a big second to generate the XLSX file.

## Build from source

0. Verify that you have Go 1.17+ installed

   ```sh
   $ go version
   ```

   If `go` is not installed, follow instructions on [the Go website](https://golang.org/doc/install).

1. Clone this repository

   ```sh
   $ git clone https://github.com/ledoyen/sqldump-to-xlsx.git
   $ cd sqldump-to-xlsx
   ```

2. Build

   ```sh
   $ go build
   ```
