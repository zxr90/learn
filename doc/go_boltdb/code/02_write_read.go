package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/renstrom/shortuuid"
)

var (
	dbPath     = shortuuid.New() + ".db"
	bucketName = shortuuid.New()

	numKeys = 10
	keyLen  = 3
	valLen  = 7

	keys = make([][]byte, numKeys)
	vals = make([][]byte, numKeys)
)

func init() {
	fmt.Println("Generating random data...")
	for i := range keys {
		keys[i] = randBytes(keyLen)
		vals[i] = randBytes(valLen)
	}
	fmt.Println("Done with random data...")
}

func main() {
	fmt.Println("dbPath:", dbPath)
	fmt.Println("bucketName:", bucketName)

	defer os.Remove(dbPath)

	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(bucketName))
		if err != nil {
			return err
		}
		for i := range keys {
			if err := b.Put(keys[i], vals[i]); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}

	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		for i := range keys {
			fmt.Printf("%s ---> %s\n", keys[i], b.Get(keys[i]))
		}
		return nil
	}); err != nil {
		panic(err)
	}
	fmt.Println("Done with db.View")
}

/*
Generating random data...
Done with random data...
dbPath: xHUfNFaPy4YPcKDhbVC3qM.db
bucketName: gSjrPgznxEa8q7roSRXpLd
zWN ---> LKckYKJ
yLp ---> lvWgBBc
BAD ---> xasSjyf
ilB ---> wVWExop
sSZ ---> kSwzVtf
Ntv ---> NkcpxBO
EVq ---> dXnWnZR
PJu ---> TTQCqLc
WEU ---> HyCQkFw
dnL ---> WYBvMaH
Done with db.View
*/

func randBytes(n int) []byte {
	const (
		letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return b
}
