package main

import (
	"fmt"
	"log"
	"net/http"

	"birdcache"
)

/*
$ curl http://localhost:9999/_birdcache/scores/Tom
630

$ curl http://localhost:9999/_birdcache/scores/kkk
kkk not exist
*/

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func getDataFromDB(key string) ([]byte, error) {
	log.Println("[SlowDB] search key", key)
	if v, ok := db[key]; ok {
		return []byte(v), nil
	}
	return nil, fmt.Errorf("%s not exist", key)
}

func main() {
	birdcache.NewGroup("scores", 2<<10, birdcache.GetterFunc(getDataFromDB))

	addr := "localhost:9999"
	peers := birdcache.NewHTTPPool(addr)

	log.Println("birdcache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
