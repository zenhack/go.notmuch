package main
// Copyright © 2015 Ian Denhardt <ian@zenhack.net>
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

import (
	".."
	"flag"
	"fmt"
)

var (
	dir = flag.String("dir", "", "Notmuch database directory")
	queryString = flag.String("query", "", "Query string")
)

func main() {
	flag.Parse()
	if *dir == "" {
		fmt.Println("Please provide a database directory.")
		flag.Usage()
		return
	}
	db, err := notmuch.Open(*dir, notmuch.DB_RDONLY)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	threads, err := db.QueryCreate(*queryString).SearchThreads()
	if err != nil {
		fmt.Println(err)
	}
	for ; threads.Valid(); threads.MoveToNext() {
		thread := threads.Get()
		fmt.Println(thread.GetSubject())
	}
}
