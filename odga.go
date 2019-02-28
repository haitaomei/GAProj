package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

/**
REST API expose: http://gaserver-svc.default.svc.cluster.local:9090
*/
func main() {
	testRedis()

	router := mux.NewRouter()

	rootRouter := router.PathPrefix("/").Subrouter()
	rootRouter.HandleFunc("/", rootHandler).Methods("GET")

	helloRouter := router.PathPrefix("/helloAPI").Subrouter()
	helloRouter.HandleFunc("/{name}", helloHandler).Methods("GET")

	http.Handle("/", router)
	http.ListenAndServe(":9090", nil) //standard http
}

func rootHandler(httpResp http.ResponseWriter, httpReq *http.Request) {

	httpResp.Header().Add("Content-Type", "application/json")
	httpResp.WriteHeader(200)
	json.NewEncoder(httpResp).Encode("Welcome to the root directory ...")
}

func helloHandler(httpResp http.ResponseWriter, httpReq *http.Request) {
	vars := mux.Vars(httpReq)
	username := vars["name"]
	var responseText = "Hi " + username + ", how are you?"

	httpResp.Header().Add("Content-Type", "application/json")
	httpResp.WriteHeader(200)
	json.NewEncoder(httpResp).Encode(responseText)
}

/*







 */
// docker run --name recorder-redis -p 6379:6379 -d redis to debug locally
//OnlineServers ...
type OnlineServers struct {
	Servers []string
}

func testRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis-master.default.svc.cluster.local:6379",
		Password: "",
		DB:       1,
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
