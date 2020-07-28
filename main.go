package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strings"
)

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]

		data, err := query(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})
	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Bonjour !"))
}

func query(city string) (weatherData, error) {
	var key string = "640d10ada2a7466c11349355946ad0bd"
	var url string = fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, key)

	resp, err := http.Get(url)

	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()
	var data weatherData

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return weatherData{}, err
	}

	return data, nil

}