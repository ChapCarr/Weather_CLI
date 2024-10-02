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
	"fyne.io/fyne/v2/layout"
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

	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		key := scanner.Text()

		return key
	} else {
		fmt.Println("Failed to read key from file, Scan returned false.")
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

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

func printEntry(city string) {
	fmt.Println(city)
}

func main() {

	app := app.New()

	//Build window
	window := app.NewWindow("Weather App")
	window.Resize(fyne.NewSize(500, 400))

	// Entry widget
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter a city")
	input.Resize(fyne.NewSize(100, 30))

	title := widget.NewLabelWithStyle("Weather Forecast", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	resultLabel := widget.NewLabel("")

	button := widget.NewButton("Enter", func() {

		city := CleanInput(string(input.Text))
		if city == "" {
			resultLabel.SetText("Please enter correct value eg 'austin' or 'san antonio'")
			return
		}
		key := getKey("key.txt")
		url := "http://api.weatherapi.com/v1/current.json?key=" + key + "&q=" + city + "&aqi=no"

		weather := GetRequest(url)

		resultText := fmt.Sprintf("City: %s\nTemperature: %.1fF / %.1fC\nCondition: %s",
			weather.Location.Name, weather.Current.TempF, weather.Current.TempC, weather.Current.Condition.Text)
		resultLabel.SetText(resultText)

	})

	content := container.NewVBox(
		title,
		container.NewPadded(input),
		container.NewPadded(button),
		layout.NewSpacer(),
		container.NewCenter(resultLabel),
		layout.NewSpacer(),
	)

	//content := container.NewVBox(input, button, resultLabel)
	window.SetContent(content)
	window.ShowAndRun()

}
