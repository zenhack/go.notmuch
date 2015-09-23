package main

// Copyright © 2015 Ian Denhardt <ian@zenhack.net>
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

import (
	"flag"
	"fmt"

	"github.com/zenhack/go.notmuch"
)

var (
	dir         = flag.String("dir", "", "Notmuch database directory")
	queryString = flag.String("query", "", "Query string")
)

func main() {
	flag.Parse()
	if *dir == "" {
		fmt.Println("Please provide a database directory.")
		flag.Usage()
		return
	}
	db, err := notmuch.Open(*dir, notmuch.DBReadOnly)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	threads, err := db.NewQuery(*queryString).Threads()
	if err != nil {
		fmt.Println(err)
	}
	for {
		thread, err := threads.Get()
		if err != nil {
			break
		}
		fmt.Println(thread.GetSubject())
		threads.MoveToNext()
	}
}