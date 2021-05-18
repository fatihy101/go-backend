package flags

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"enstrurent.com/server/db"
	"go.mongodb.org/mongo-driver/bson"
)

type Configuration struct {
	SERVER_PORT string
	CON_STR     string
	DBNAME      string
	JWT_KEY     string
}

var config Configuration

func InitConfig() Configuration {
	err := os.Mkdir("assets/images", 0755)
	if err != nil {
		fmt.Println("WARNING: " + err.Error())
	}
	// Open our jsonFile
	jsonFile, err := os.Open("assets/json/config.json")

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &config)
	return config
}

func GetConfig() Configuration {
	return config
}

type mytype []map[string]interface{}

func InitCities(mdb *db.DBHandle) {
	ctx := context.TODO()
	collection := mdb.CitiesCollection()
	cursor, _ := collection.Find(ctx, bson.D{})
	if !cursor.Next(ctx) { // If collection empty
		var data mytype
		byteValue, _ := ioutil.ReadFile("assets/json/cities.json")
		err := json.Unmarshal(byteValue, &data)
		if err != nil {
			fmt.Print(err)
		}
		var cities []interface{}
		for _, t := range data {
			cities = append(cities, t)
		}
		if _, err = collection.InsertMany(context.Background(), cities); err != nil {
			fmt.Println(err)
		}
		return
	}
}
