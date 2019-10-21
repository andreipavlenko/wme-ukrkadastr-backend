package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type localitiy struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

var localities = map[string]string{}

func main() {
	initData()
	startServer()
}

func initData() {
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
}

func getLocalities(data []interface{}, currentLevel int) {
	for _, v := range data {
		value := v.(map[string]interface{})
		code, name := value["code"].(string), value["name"].(string)
		localities[code] = name
		nextLevelIndex := fmt.Sprintf("level%v", currentLevel+1)
		nextLevel, ok := value[nextLevelIndex]
		if ok {
			getLocalities(nextLevel.([]interface{}), currentLevel+1)
		}
	}
}

func startServer() {
	http.HandleFunc("/locality", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request URI: %v\n", r.URL.RequestURI())
		w.Header().Set("Content-Type", "application/json")
		code, ok := r.URL.Query()["code"]
		if ok {
			l, ok := localities[code[0]]
			if ok {
				res := localitiy{
					Name: l,
					Code: code[0],
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
	log.Fatal(http.ListenAndServe(":8080", nil))
}
