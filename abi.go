package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"gopkg.in/yaml.v2"
)

func main() {
	jsonData, err := ioutil.ReadFile("api.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}
	abiJSON := []byte(result["result"].(string))

	abiObj, err := abi.JSON(abiJSON)
	if err != nil {
		log.Fatalf("Error parsing ABI: %v", err)
	}

	yamlData, err := yaml.Marshal(abiObj)
	if err != nil {
		log.Fatalf("Error marshaling to YAML: %v", err)
	}

	fmt.Println(string(yamlData))
}
