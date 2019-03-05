package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

/**
REST API expose: http://gaserver-svc.default.svc.cluster.local:9090

Using
docker run --name recorder-redis -p 6379:6379 -d redis
to debug locally
*/

var client *redis.Client
var islands = make([]string, 0)

func main() {
	// init db
	client = redis.NewClient(&redis.Options{
		Addr:     "redis-master.default.svc.cluster.local:6379",
		Password: "",
		DB:       1,
	})

	router := mux.NewRouter()

	rootRouter := router.PathPrefix("/").Subrouter()
	rootRouter.HandleFunc("/", rootHandler).Methods("GET")

	pushRouter := router.PathPrefix("/push").Subrouter()
	pushRouter.HandleFunc("/{name}", pushHandler).Methods("POST")

	pollRouter := router.PathPrefix("/poll").Subrouter()
	pollRouter.HandleFunc("/{name}", pollHandler).Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(":9090", nil)

}

func rootHandler(httpResp http.ResponseWriter, httpReq *http.Request) {

	httpResp.Header().Add("Content-Type", "application/json")
	httpResp.WriteHeader(200)
	json.NewEncoder(httpResp).Encode("Welcome to the root directory ...")
}

func pushHandler(httpResp http.ResponseWriter, httpReq *http.Request) {
	vars := mux.Vars(httpReq)
	islandID := vars["name"]

	body, _ := ioutil.ReadAll(httpReq.Body)
	fmt.Println("Received push request, going to save into db\n\tid=", islandID, "body:", body)
	// save to db
	err := client.Set(islandID, string(body), 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func pollHandler(httpResp http.ResponseWriter, httpReq *http.Request) {
	vars := mux.Vars(httpReq)
	islandID := vars["name"]
	fmt.Println("Received poll request, going to save into db\n\tid=", islandID)
	//read from db
	val2, err := client.Get(islandID).Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		fmt.Println(err)
	} else {
		httpResp.Header().Add("Content-Type", "application/json")
		httpResp.WriteHeader(200)
		fmt.Fprintf(httpResp, "%s", val2)
	}
}
