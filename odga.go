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

// type ParetoFront struct {
// 	Fields []float64
// }

// type JsonType struct {
// 	Array []ParetoFront
// }

/**
REST API expose: http://gaserver-svc.default.svc.cluster.local:9090
*/
func main() {
	// testRedis()

	router := mux.NewRouter()

	rootRouter := router.PathPrefix("/").Subrouter()
	rootRouter.HandleFunc("/", rootHandler).Methods("GET")

	pushRouter := router.PathPrefix("/push").Subrouter()
	pushRouter.HandleFunc("/{name}", pushHandler).Methods("POST")

	pollRouter := router.PathPrefix("/poll").Subrouter()
	pollRouter.HandleFunc("/{name}", pollHandler).Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(":9090", nil) //standard http
}

func rootHandler(httpResp http.ResponseWriter, httpReq *http.Request) {

	httpResp.Header().Add("Content-Type", "application/json")
	httpResp.WriteHeader(200)
	json.NewEncoder(httpResp).Encode("Welcome to the root directory ...")
}

func pushHandler(httpResp http.ResponseWriter, httpReq *http.Request) {
	vars := mux.Vars(httpReq)
	islandID := vars["name"]
	fmt.Println(islandID, httpReq.Body)
}

func pollHandler(httpResp http.ResponseWriter, httpReq *http.Request) {
	vars := mux.Vars(httpReq)
	islandID := vars["name"]
	fmt.Println(islandID)
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
