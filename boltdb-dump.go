package main

import (
	"fmt"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

var (
	template   = "INSERT INTO media VALUES('%s', '%s');\n"
	bucketName = "files"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <filename.db>\n", os.Args[0])
		os.Exit(2)
	}

	filename := os.Args[1]

	db, err := bolt.Open(filename, 0600, &bolt.Options{
		ReadOnly: true,
		Timeout:  1 * time.Second,
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(2)
	}

	defer db.Close()

	dump(db)
}

func dump(db *bolt.DB) {
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf(template, k, v)
		}

		return nil
	})
}
