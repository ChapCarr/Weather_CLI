package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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

func CleanInput(city string) string {
	// Replace whitespace with _ and remove the newline char
	cityClean := strings.Replace(city, " ", "_", -1)
	return strings.ReplaceAll(cityClean, "\n", "")
}

func errHandler(err error) {
	if err != nil {
		fmt.Printf("ERROR: %s", err)
	}
}

func GetRequest(url string) WeatherResponse {
	// Make the request
	resp, err := http.Get(url)
	errHandler(err)
	defer resp.Body.Close()

	//Read in the response
	body, err := io.ReadAll(resp.Body)
	errHandler(err)

	// Process the JSON
	var weather WeatherResponse
	err = json.Unmarshal(body, &weather)
	errHandler(err)

	return weather

}

func printEntry(city string){
	fmt.Println(city)
}

func main() {
	

 	scanner := bufio.NewReader(os.Stdin)
	app := app.New()

	//Build window
	window := app.NewWindow("Weather App")
	window.Resize(fyne.NewSize(500,500))

	// Entry widget
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter a city")

	content := container.NewVBox(input, widget.NewButton("Enter",func(){
		 CleanInput(input.Text)
	}))
	
	window.SetContent(content)

	window.ShowAndRun()

	// Scan and Parse user input
	fmt.Print("Enter a city: ")
	cityRaw, err := scanner.ReadString('\n')
	errHandler(err)
	cityURL := CleanInput(cityRaw)

	url := "http://api.weatherapi.com/v1/current.json?key=" + getKey("key.txt") + "&q=" + cityURL + "&aqi=no"
	
	weather := GetRequest(url) // GET REQUEST
	
	
	window.ShowAndRun()
	fmt.Printf("City: %s\n", weather.Location.Name)
	fmt.Printf("Temperature %.1fF / %.1fC \n", weather.Current.TempF, weather.Current.TempC)

}
