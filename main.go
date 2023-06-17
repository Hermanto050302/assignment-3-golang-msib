package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Status struct {
	Water  int    `json:"water"`
	Wind   int    `json:"wind"`
	Status string `json:"status"`
}

type Data struct {
	Status `json:"status"`
}

func updateData() {
	for {
		var data = Data{Status: Status{}}

		data.Status.Water = rand.Intn(100) + 1
		data.Status.Wind = rand.Intn(100) + 1
		data.Status.Status = getStatus(data.Status.Water, data.Status.Wind)

		b, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			log.Fatalln("error while marshalling json data:", err)
		}

		err = ioutil.WriteFile("data.json", b, 0644)
		if err != nil {
			log.Fatalln("error while writing value to data.json file:", err)
		}

		fmt.Println("Menunggu selama 5 detik...")
		time.Sleep(time.Second * 5)
	}
}

func getStatus(water, wind int) string {
	if water < 5 || wind < 6 {
		return "Aman"
	} else if water >= 6 && water <= 8 || wind >= 7 && wind <= 15 {
		return "Siaga"
	} else {
		return "Bahaya"
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	go updateData()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl, err := template.ParseFiles("index.html")
		if err != nil {
			log.Println("error while parsing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var data = Data{Status: Status{}}

		b, err := ioutil.ReadFile("data.json")
		if err != nil {
			log.Println("error while reading data.json file:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(b, &data)
		if err != nil {
			log.Println("error while unmarshalling JSON data:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tpl.ExecuteTemplate(w, "index.html", data.Status)
		if err != nil {
			log.Println("error while executing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
