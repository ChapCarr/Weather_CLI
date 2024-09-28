package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type WeatherResponse struct {
	Location struct {
		Name      string `json:"name"`
		Region    string `json:"region"`
		Localtime string `json:"localtime"`
	} `json:"location"`

	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`

		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
}

func getKey(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	defer file.Close()

	return ""
}

func main() {
	var cityRaw string
	scanner := bufio.NewReader(os.Stdin)

	// Scan and Parse user input
	fmt.Print("Enter a city: ")
	cityRaw, err := scanner.ReadString('\n')
	cityClean := strings.Replace(cityRaw, " ", "_", -1)
	cityClean = strings.ReplaceAll(cityClean, "\n", "")

	//GET request
	url := "http://api.weatherapi.com/v1/current.json?key="+getKey("key.txt")+"&q=" + cityClean + "&aqi=no"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// Read in
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var weather WeatherResponse
	err = json.Unmarshal(body, &weather)
	if err != nil {
		log.Fatalln("ERROR: unmarshalling: ", err)
	}

	fmt.Printf("City: %s\n", weather.Location.Name)
	fmt.Printf("Temperature %.1fF / %.1fC \n", weather.Current.TempF, weather.Current.TempC)

}
