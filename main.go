package main

import (
	"Gcache/cache"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"tom":  "100",
	"jack": "10",
	"sam":  "60",
}

func startCacheServer(addr string, addrs []string, gee *cache.Group) {
	peers := cache.NewHTTPPool(addr)
	peers.Set(addrs...)
	gee.RegisterPeers(peers)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, gee *cache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := gee.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}
	c := cache.NewGroup("score", 2<<10, cache.GetterFunc(func(key string) ([]byte, error) {
		log.Println("slow DB search key:", key)
		fmt.Println(db, key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		} else {
			return nil, errors.New("not found in db ")
		}
	}))
	if api {
		go startAPIServer(apiAddr, c)
	}
	startCacheServer(addrMap[port], addrs, c)

}
