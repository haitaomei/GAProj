package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strings"

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

type objective struct {
	List [][]float64 `json:"list"`
}

func main() {
	// init db
	client = redis.NewClient(&redis.Options{
		Addr:     "redis-master.default.svc.cluster.local:6379",
		Password: "",
		DB:       1,
	})
	connectRedis()

	router := mux.NewRouter()

	rootRouter := router.PathPrefix("/").Subrouter()
	rootRouter.HandleFunc("/", rootHandler).Methods("GET")

	pushRouter := router.PathPrefix("/push").Subrouter()
	pushRouter.HandleFunc("/{name}", pushHandler).Methods("POST")

	pushObjectiveRouter := router.PathPrefix("/pushobject").Subrouter()
	pushObjectiveRouter.HandleFunc("/{name}", pushObjectiveHandler).Methods("POST")

	pollRouter := router.PathPrefix("/poll").Subrouter()
	pollRouter.HandleFunc("/{name}", pollHandler).Methods("POST")
	/* alia for poll */
	pullRouter := router.PathPrefix("/pull").Subrouter()
	pullRouter.HandleFunc("/{name}", pollHandler).Methods("POST")

	getAllIslandsRouter := router.PathPrefix("/allislands").Subrouter()
	getAllIslandsRouter.HandleFunc("", getAllIslandsHandler).Methods("GET")

	http.Handle("/", router)
	fmt.Println("* Elastic GA Data Tier Server is starting... (listening on http)")
	log.Fatal(http.ListenAndServe(":9090", nil))
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
	fmt.Println("* Push request, going to save into db.\tID =", islandID) //, "body:", body
	// save to db
	err := client.Set(islandID, string(body), 0).Err()
	if err != nil {
		fmt.Println(err)
	}
	updateAllIslandIDs(islandID)
}

func pushObjectiveHandler(httpResp http.ResponseWriter, httpReq *http.Request) {
	vars := mux.Vars(httpReq)
	islandID := vars["name"]

	body, _ := ioutil.ReadAll(httpReq.Body)
	fmt.Println("* Push request, going to save into db.\tID =", islandID) //, "body:", body
	// save to db
	err := client.Set(islandID+"_Objective_Key", string(body), 0).Err()
	if err != nil {
		fmt.Println(err)
	}
	updateAllIslandIDs(islandID)
}

func pollHandler(httpResp http.ResponseWriter, httpReq *http.Request) {
	vars := mux.Vars(httpReq)
	islandID := vars["name"]
	fmt.Println("* Pull request, ID =", islandID)

	// random select one island
	updateAllIslandIDs(islandID)
	count := len(islands)
	selectedIsland := islandID
	for i := 1; i <= 10; i++ {
		selected := rand.Intn(count)
		if islands[selected] == islandID {
			selected = (selected + 1) % count
		}
		if islands[selected] != "" {
			selectedIsland = islands[selected]
			break
		}
	}

	fmt.Println("\t- randomly selected island:", selectedIsland)
	// read from db
	res, err := client.Get(selectedIsland).Result()
	if err == redis.Nil {
		fmt.Println(selectedIsland, " does not exist")
	} else if err != nil {
		fmt.Println(err)
	} else {
		httpResp.Header().Add("Content-Type", "application/json")
		httpResp.WriteHeader(200)
		// do not use json encoder here again, because we saved the data without decode
		fmt.Fprintf(httpResp, "%s", res)
	}
}

// tracking all the islands, dummy implementation
func updateAllIslandIDs(islandID string) {
	refreshCurRecords()
	// appending current island
	needAppend := true
	for _, s := range islands {
		if s == islandID {
			needAppend = false
			break
		}
	}
	if needAppend {
		islands = append(islands, islandID)
	}

	// write into redis
	record := ""
	for _, s := range islands {
		if s != "" {
			record += (s + ";")
		}
	}
	client.Set("AllIsLands", record, 0).Err()

	fmt.Println("\t- All islands recorded:", record)
}

func refreshCurRecords() {
	islandsInDB := make([]string, 0)
	records, err := client.Get("AllIsLands").Result()
	if err == nil {
		islandsInDB = strings.Split(records, ";")
	}

	if len(islandsInDB) > 0 {
		existing := make(map[string]bool)
		for _, s := range islands {
			existing[s] = true
		}
		for _, s := range islandsInDB {
			if s != "" && !existing[s] {
				islands = append(islands, s)
			}
		}
	}
}

func getAllIslandsHandler(httpResp http.ResponseWriter, httpReq *http.Request) {
	record := "==========All the records of islands till now===========\n"
	record += "-------------Islands List--------------\n"
	refreshCurRecords()
	islandsRecords := make([]string, len(islands))
	copy(islandsRecords, islands)
	sort.Strings(islandsRecords)
	for _, s := range islandsRecords {
		if s != "" {
			record += ("Island ID: " + s + "\n")
		}
	}
	record += "\n-------------Contents--------------\n"
	for _, s := range islandsRecords {
		if s != "" {
			objectiveData, _ := client.Get(s + "_Objective_Key").Result()
			od := objective{}
			err := json.Unmarshal([]byte(objectiveData), &od)
			if err == nil {
				fmt.Printf("%+v\n", od)
			} else {
				fmt.Println(err, objectiveData)
			}
			record += objectiveData
			record += "\n"
		}
	}
	fmt.Fprintf(httpResp, "%s", record)
}

func connectRedis() {
	for true {
		err := client.Set("RedisConnection", "Connected", 0).Err()
		if err != nil {
			fmt.Println("Waitting to connect to Redis...", err)
		} else {
			break
		}
	}
}
