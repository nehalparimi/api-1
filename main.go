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

type weatherProvider interface {
	temperature(city string) (float64, error)
}

type openWeatherMap struct {}

func (w openWeatherMap) temperature(city string) (float64, error) {
	var key string = "640d10ada2a7466c11349355946ad0bd"
	var url string = fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, key)

	resp, err := http.Get(url)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	var d struct {
		Main struct {
			Kelvin float64 `json:"temp"`
		} `json:"main"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return 0, err
	}

	log.Printf("openWeatherMap: %s: %.2f", city, d.Main.Kelvin)
	return d.Main.Kelvin, nil
}

type darkSky struct {
	apiKey string
}

func (w darkSky) temperature(city string) {
	var url string = "https://api.darksky.net/forecast/2d12851525d65bbbefac0728f7b7a9b4/37.8267,-122.4233"

	resp, err := http.Get(url)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	var d struct {
		Observation struct {
			Celsius float64 `json:"temp_c"`
		} `json:"current_observation"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return 0, err
	}

	kelvin := d.Observation.Celsius + 273.15
	log.Printf("Dark Sky: %s, %.2f", city, kelvin)

	return kelvin, nil
}

func temperature(city string, providers ...weatherProvider) (float64, error) {
	sum := 0.0

	for _, provider := range providers {
		k, err := provider.temperature(city)

		if err != nil {
			return 0, err
		}

		sum += k
	}

	return sum / float64(len(providers)), nil
}
