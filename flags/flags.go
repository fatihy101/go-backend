package flags

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Configuration struct {
	SERVER_PORT string
	CON_STR     string
	DBNAME      string
	JWT_KEY     string
}

var config Configuration

func InitConfig() Configuration {
	// Open our jsonFile
	jsonFile, err := os.Open("config.json")

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
