package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin   float64 `json:"temp"`
		Max      float64 `json:"temp_max"`
		Min      float64 `json:"temp_min"`
		Pressure int     `json:"pressure"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Coordinates struct {
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lon"`
	} `json:"coord"`
	Wind struct {
		Speed  float64 `json:"speed"`
		Degree float64 `json:"deg"`
		Gust   float64 `json:"gust"`
	} `json:"wind"`
}

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return apiConfigData{}, err
	}

	var configData apiConfigData
	err = json.Unmarshal(bytes, &configData)
	if err != nil {
		return apiConfigData{}, err
	}
	return configData, nil
}

func hello(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Hello from Server\n"))

}
func query(city string) (weatherData, error) {
	apiConfigData, err := loadApiConfig("config.json")
	if err != nil {
		return weatherData{}, err
	}
	url := "https://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfigData.OpenWeatherMapApiKey + "&q=" + city
	fmt.Println("Getting: ", url)
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
func main() {
	http.HandleFunc("/hello", hello)

	http.HandleFunc("/weather/", func(rw http.ResponseWriter, r *http.Request) {
		city := strings.SplitAfterN(r.URL.Path, "/", 3)[2]
		data, err := query(city)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(rw).Encode(data)
	})

	fmt.Println("Server Started at port", ":8080")
	http.ListenAndServe(":8080", nil)

}
