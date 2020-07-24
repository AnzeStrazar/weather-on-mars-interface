package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"weather-on-mars-interface/cache"
	"weather-on-mars-interface/resolvers"
	"weather-on-mars-interface/types"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

type NasaData struct {
	cache *cache.Cache
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
	newSols := make([]types.Sol, 0)
	for _, solId := range solKeys {
		sol := types.Sol{
			SolId: solId.(string),
			Data: struct {
				At  types.AT  `json:"AT"`
				Hws types.HWS `json:"HWS"`
				Pre types.PRE `json:"PRE"`
			}{},
		}
		solData := serializedData[solId.(string)]
		atmosphericTemperature := solData.(map[string]interface{})["AT"].(map[string]interface{})
		sol.Data.At.Av = atmosphericTemperature["av"].(float64)
		sol.Data.At.Ct = atmosphericTemperature["ct"].(float64)
		sol.Data.At.Mn = atmosphericTemperature["mn"].(float64)
		sol.Data.At.Mx = atmosphericTemperature["mx"].(float64)

		horizontalWindSpeed := solData.(map[string]interface{})["HWS"].(map[string]interface{})
		sol.Data.Hws.Av = horizontalWindSpeed["av"].(float64)
		sol.Data.Hws.Ct = horizontalWindSpeed["ct"].(float64)
		sol.Data.Hws.Mn = horizontalWindSpeed["mn"].(float64)
		sol.Data.Hws.Mx = horizontalWindSpeed["mx"].(float64)

		atmosphericPressure := solData.(map[string]interface{})["PRE"].(map[string]interface{})
		sol.Data.Pre.Av = atmosphericPressure["av"].(float64)
		sol.Data.Pre.Ct = atmosphericPressure["ct"].(float64)
		sol.Data.Pre.Mn = atmosphericPressure["mn"].(float64)
		sol.Data.Pre.Mx = atmosphericPressure["mx"].(float64)

		newSols = append(newSols, sol)
	}
	fmt.Printf("%v", newSols)
	nd.saveToCache(newSols)

	// Save to mongoDB
	// nd.saveToMongo()
}

func (nd *NasaData) saveToMongo() {
	// Save data to mongo
	/*
		for solId, data := range nd.Sols {
			//save (solId, data )
			// Data to JSON and SAVE the data to mongoDB
			nd.mongoDBClient.Set(data.SolId, "smth 	"+data.SolId)
			log.Println(nd.mongoDBClient.Get(data.SolId))

		}
	*/
}

func (nd *NasaData) saveToCache(newSols []types.Sol) {
	nd.cache.SolMutex.Lock()
	defer nd.cache.SolMutex.Unlock()
	for _, solData := range newSols {
		nd.cache.SolCache[solData.SolId] = solData
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
	// We re going to start our periodic data updater here.
	cache := cache.Cache{SolCache: make(map[string]types.Sol)}
	nasaData := NasaData{&cache}
	go nasaData.runUpdater()

	http.Handle("/graphql", &relay.Handler{
		Schema: parseSchema("./schema.graphql", resolvers.NewRootResolver(&cache)),
	})

	fmt.Println("serving on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
