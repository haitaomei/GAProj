package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/go-redis/redis"
)

// docker run --name recorder-redis -p 6379:6379 -d redis to debug locally

//OnlineServers ...
type OnlineServers struct {
	Servers []string
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis-cluster.default.svc.cluster.local:6379", //redis-cluster.default.svc.cluster.local:6379
		Password: "",
		// DB:       1,
	})

	// pong, err := client.Ping().Result()

	err := client.Set("key1", "value1", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get("key1").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key1", val)

	val2, err := client.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}

	var testsvrs []string
	testsvrs = append(testsvrs, "localhost")
	testsvrs = append(testsvrs, "localhost2")
	testsvrs = append(testsvrs, "localhost3")
	osvrs := &OnlineServers{
		Servers: testsvrs,
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(osvrs)
	client.Set("ServerList", b.String(), 0)
	val3, err := client.Get("ServerList").Result()

	dec := json.NewDecoder(strings.NewReader(val3))
	for {
		var m OnlineServers
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Println(m)
	}

}
