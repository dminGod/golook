package main

import (
	"time"
	"fmt"
	"github.com/timshannon/bolthold"
	"os"
	"go.etcd.io/bbolt"
)

var store *bolthold.Store


func initLocalDB() {

	var err error

	store, err = bolthold.Open("local.db", 0666, &bolthold.Options{ Options : &bbolt.Options{Timeout: 2 * time.Second}})

	if err != nil {
		fmt.Printf( "Error creating the local database file - '%v' -- Will exit application -- Error: %v \n", "local.db", err )
		os.Exit(1)
		return
	}
}

func addApp(a Application) (err error) {

	err = store.Upsert(a.Name, &a)
	if err != nil {
		fmt.Printf("Got error when trying to upsert record for file read complete reference, File Ref: '%+v', Error: '%v' \n",
			a, err.Error())
	}

	return
}

func getApp(aName string) (a Application) {

	err := store.Get(aName, &a)
	if err != nil {
		if err == bolthold.ErrNotFound {
			fmt.Printf("No records found for the inode number : %v \n", aName)
			return
		}
	}

	return
}