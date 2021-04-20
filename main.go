package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Response struct {
	MRData MRData `json:"MRData"`
}

type MRData struct {
	Series    string    `json:"series,omitempty"`
	Total     int       `json:"total,omitempty"`
	RaceTable RaceTable `json:"RaceTable"`
}

type RaceTable struct {
	Season string `json:"season,omitempty"`
	Races  []Race `json:"Races"`
}

type Race struct {
	Season      int       `json:"season,omitempty"`
	Round       int       `json:"round,omitempty"`
	Url         string    `json:"url,omitempty"`
	RaceName    string    `json:"raceName,omitempty"`
	CircuitName string    `json:"circuitName,omitempty"`
	Date        string    `json:"date,omitempty"`
	Time        string    `json:"time,omitempty"`
	DateTime    time.Time `json:"omitempty"`
}

func main() {
	for _, race := range getSchedule().RaceTable.Races {
		t, _ := time.Parse(time.RFC3339, race.Date+"T"+race.Time)

		if t.After(time.Now()) {
			race.DateTime = t
			SendNotifications(race)
			break
		}
	}
}

func getSchedule() MRData {
	resp, err := http.Get("http://ergast.com/api/f1/current.json")
	if err != nil {
		log.Fatalln(err)
	}

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var response Response
	json.Unmarshal(body, &response)

	return response.MRData
}

func SendNotifications(race Race) {
	for _, number := range strings.Split(os.Getenv("PHONE_NUMBERS"), ",") {
		fmt.Println(race.Time)
		sendSMS(number, race)
	}
}

func sendSMS(number string, race Race) {
	data := url.Values{}
	data.Set("To", number)
	data.Set("MessagingServiceSid", os.Getenv("MessagingServiceSid"))

	loc, _ := time.LoadLocation("Europe/Dublin")
	data.Set("Body", "Next Race: "+race.RaceName+" at: "+race.DateTime.In(loc).Format(time.RFC1123))

	url, _ := url.ParseRequestURI("https://api.twilio.com/2010-04-01/Accounts/" + os.Getenv("TWILIO_ACCOUNT_SID") + "/Messages.json")

	req, err := http.NewRequest("POST", url.String(), strings.NewReader(data.Encode()))

	req.Header.Set("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(os.Getenv("TWILIO_ACCOUNT_SID"), os.Getenv("TWILIO_TOKEN"))

	cli := &http.Client{}
	resp, err := cli.Do(req)

	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Printf(sb)
}
