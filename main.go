package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type localitiy struct {
	Name       string `json:"name"`
	Code       string `json:"code"`
	ZoneNumber string `json:"zone_number"`
}

var localities = map[string]string{}

var localitiesCount = 0

func main() {
	initData()
	startServer()
}

func initData() {
	log.Println("Initializing data..")
	content, err := ioutil.ReadFile("koatuu.json")
	if err != nil {
		return
	}
	var data map[string]interface{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		return
	}
	getLocalities(data["level1"].([]interface{}), 1)
	log.Printf("%v localitites initialized.\n", localitiesCount)
}

func getLocalities(data []interface{}, currentLevel int) {
	for _, v := range data {
		value := v.(map[string]interface{})
		code, name := value["code"].(string), value["name"].(string)
		localities[code] = name
		localitiesCount++
		nextLevelIndex := fmt.Sprintf("level%v", currentLevel+1)
		nextLevel, ok := value[nextLevelIndex]
		if ok {
			getLocalities(nextLevel.([]interface{}), currentLevel+1)
		}
	}
}

func startServer() {
	log.Println("Starting server..")
	http.HandleFunc("/locality", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request URI: %v\n", r.URL.RequestURI())
		w.Header().Set("Content-Type", "application/json")
		code, okCode := r.URL.Query()["code"]
		zoneNumber, okZn := r.URL.Query()["zone_number"]
		if okCode && okZn {
			name := getLocalityName(code[0], zoneNumber[0])
			if name != "" {
				res := localitiy{
					Name:       name,
					Code:       code[0],
					ZoneNumber: zoneNumber[0],
				}
				js, err := json.Marshal(res)
				if err == nil {
					w.Write(js)
					return
				}
			}
		}
		w.Write([]byte("{}"))
	})
	log.Fatal(http.ListenAndServe(":7979", nil))
}

func getLocalityName(koatuu, zoneNumber string) string {
	if len(koatuu) < len(zoneNumber) {
		return ""
	}
	koatuuWithZone := koatuu[:len(koatuu)-len(zoneNumber)] + zoneNumber
	name, ok := localities[koatuuWithZone]
	if ok {
		return name
	}
	name, ok = localities[koatuu]
	if ok {
		return name
	}
	return ""
}
