package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"weather-on-mars-interface/cache"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

type RootResolver struct {
	redisClient *cache.Redis
}

type NasaData struct {
	redisClient *cache.Redis
	Sols        []Sol
}

type Sol struct {
	SolId string
	Data  struct {
		At AT `json:"AT"`
	}
}

type AT struct {
	Av float64 `json:"av"`
	Ct float64 `json:"ct"`
	Mn float64 `json:"mn"`
	Mx float64 `json:"mx"`
}

func (r *RootResolver) Info() (string, error) {
	return string(r.redisClient.Get("584")), nil
}

// Reads and parses the schema from file.
// Associates root resolver. Panics if can't read.
func parseSchema(path string, resolver interface{}) *graphql.Schema {
	bstr, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	schemaString := string(bstr)
	parsedSchema, err := graphql.ParseSchema(
		schemaString,
		resolver,
	)
	if err != nil {
		panic(err)
	}
	return parsedSchema
}

func newNasaData(rds *cache.Redis) NasaData {
	return NasaData{
		redisClient: rds,
	}
}

func (nd *NasaData) runUpdater() {
	for {
		nd.update()
		time.Sleep(8 * time.Second) // Define frequency for request triggering.
	}
}

func (nd *NasaData) update() {
	response, err := http.Get("https://api.nasa.gov/insight_weather/?api_key=6IoPDJhsNThUqrwvKhNKBMe1z5gNnKYj8hxosLCh&feedtype=json&ver=1.0")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return
	}
	defer response.Body.Close()

	serializedData := make(map[string]interface{})
	json.NewDecoder(response.Body).Decode(&serializedData)
	solKeys := serializedData["sol_keys"].([]interface{})
	newSols := make([]Sol, 0)
	for _, solId := range solKeys {
		sol := Sol{
			SolId: solId.(string),
			Data: struct {
				At AT `json:"AT"`
			}{},
		}
		solData := serializedData[solId.(string)]
		atmosfericTemperature := solData.(map[string]interface{})["AT"].(map[string]interface{})
		sol.Data.At.Av = atmosfericTemperature["av"].(float64)
		sol.Data.At.Ct = atmosfericTemperature["ct"].(float64)
		sol.Data.At.Mn = atmosfericTemperature["mn"].(float64)
		sol.Data.At.Mx = atmosfericTemperature["mx"].(float64)
		newSols = append(newSols, sol)
	}
	nd.Sols = newSols
	nd.saveToRedis()

	// Save to mongoDB
}

func (nd *NasaData) saveToMongo() {
	// Save data to mongo
}

func (nd *NasaData) saveToRedis() {
	for _, data := range nd.Sols {
		nd.redisClient.Set(data.SolId, "smth 	"+data.SolId)
		log.Println(nd.redisClient.Get(data.SolId))
	}
}

// You simply want to fetch the data. We dont care whether it exists
// on disk or memory.
func (nd *NasaData) get(solId string) {
	// If solId in current sols (last 7) -> get from redis
	// If it does not exisst in redis -> get from DB
	// If sol is older than 7 sols (is not returned by NASA) -> fetch from DB
}

func main() {
	rds := cache.NewRedis("redis://:redispassword@127.0.0.1:48000")
	nasaData := newNasaData(rds)
	// We re going to start our periodic data updater here.
	go nasaData.runUpdater()

	http.Handle("/graphql", &relay.Handler{
		Schema: parseSchema("./schema.graphql", &RootResolver{redisClient: rds}),
	})

	fmt.Println("serving on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
