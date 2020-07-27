package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"weather-on-mars-interface/cache"
	"weather-on-mars-interface/database"
	"weather-on-mars-interface/resolvers"
	"weather-on-mars-interface/types"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NasaData struct {
	cache  *cache.Cache
	client *mongo.Client
}

// Reads and parses the schema from file.
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

// TODO: Add Subscription instead of sleep delay.
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
		}
		solData := serializedData[solId.(string)]
		atmosphericTemperature := solData.(map[string]interface{})["AT"].(map[string]interface{})
		sol.At.Av = atmosphericTemperature["av"].(float64)
		sol.At.Ct = atmosphericTemperature["ct"].(float64)
		sol.At.Mn = atmosphericTemperature["mn"].(float64)
		sol.At.Mx = atmosphericTemperature["mx"].(float64)

		horizontalWindSpeed := solData.(map[string]interface{})["HWS"].(map[string]interface{})
		sol.Hws.Av = horizontalWindSpeed["av"].(float64)
		sol.Hws.Ct = horizontalWindSpeed["ct"].(float64)
		sol.Hws.Mn = horizontalWindSpeed["mn"].(float64)
		sol.Hws.Mx = horizontalWindSpeed["mx"].(float64)

		atmosphericPressure := solData.(map[string]interface{})["PRE"].(map[string]interface{})
		sol.Pre.Av = atmosphericPressure["av"].(float64)
		sol.Pre.Ct = atmosphericPressure["ct"].(float64)
		sol.Pre.Mn = atmosphericPressure["mn"].(float64)
		sol.Pre.Mx = atmosphericPressure["mx"].(float64)

		newSols = append(newSols, sol)
	}

	nd.saveToCache(newSols)
	nd.saveToMongo(newSols)

}

func (nd *NasaData) saveToMongo(newSols []types.Sol) {
	solsCollections := nd.client.Database("weather-on-mars").Collection("sols")

	for _, solData := range newSols {
		_, err := solsCollections.UpdateOne(context.TODO(), bson.M{"solID": bson.M{"$eq": solData.SolId}}, bson.M{"$set": solData}, options.Update().SetUpsert(true))
		if err != nil {
			fmt.Printf("error %s\n", err)
			return
		}
	}
}

func (nd *NasaData) saveToCache(newSols []types.Sol) {
	nd.cache.SolMutex.Lock()
	defer nd.cache.SolMutex.Unlock()
	for _, solData := range newSols {
		nd.cache.SolCache[solData.SolId] = solData
	}
}

func main() {
	dbHost := "localhost"
	dbPort := "48000"

	client := database.NewMongoDB(dbHost, dbPort)

	// We re going to start our periodic data updater here.
	cache := cache.Cache{SolCache: make(map[string]types.Sol)}
	nasaData := NasaData{&cache, client}

	go nasaData.runUpdater()

	http.Handle("/graphql", &relay.Handler{
		Schema: parseSchema("./schema.graphql", resolvers.NewRootResolver(&cache, client)),
	})

	fmt.Println("serving on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
