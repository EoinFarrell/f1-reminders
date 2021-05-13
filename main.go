package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type ErgastResponse struct {
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

type DbResponse struct {
	Users []User `json:"users"`
}

type User struct {
	UserId   string `json:"userId,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Email    string `json:"email,omitempty"`
	Timezone string `json:"timezone,omitempty"`
}

func main() {
	for _, race := range GetSchedule().RaceTable.Races {
		t, _ := time.Parse(time.RFC3339, race.Date+"T"+race.Time)

		if t.After(time.Now()) {
			race.DateTime = t
			SendNotifications(race)
			break
		}
	}
}

func GetSchedule() MRData {
	resp, err := http.Get("http://ergast.com/api/f1/current.json")
	if err != nil {
		log.Fatalln(err)
	}

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var response ErgastResponse
	json.Unmarshal(body, &response)

	return response.MRData
}

func GetUsers() DbResponse {
	url, _ := url.ParseRequestURI("https://phrhyp7dx2.execute-api.eu-west-1.amazonaws.com/Production")

	req, err := http.NewRequest("GET", url.String(), nil)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	req.Header.Set("Accept", "application/json")

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

	var response DbResponse
	json.Unmarshal(body, &response)

	return response
}

func SendNotifications(race Race) {
	for _, user := range GetUsers().Users {
		SendSMS(user, race)
	}
}

func SendSMS(user User, race Race) {
	data := url.Values{}
	data.Set("To", user.Phone)
	data.Set("MessagingServiceSid", os.Getenv("MessagingServiceSid"))

	loc, _ := time.LoadLocation(user.Timezone)
	data.Set("Body", "Next Race: "+race.RaceName+" at: "+race.DateTime.In(loc).Format(time.RFC1123)+". https://www.formula1.com/")

	url, _ := url.ParseRequestURI("https://api.twilio.com/2010-04-01/Accounts/" + os.Getenv("TWILIO_ACCOUNT_SID") + "/Messages.json")

	req, err := http.NewRequest("POST", url.String(), strings.NewReader(data.Encode()))
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

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
