package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func recuperavariabile(variabile string) (result string, err error) {
	if result, ok := os.LookupEnv(variabile); ok && len(result) != 0 {
		return result, nil
	}
	return "", fmt.Errorf("la variabile %s non esiste o è vuota", variabile)
}

func main() {

	sid := "CA10c7de2b487b59ceb51917ab81aa2367"

	for n := 0; n < 3; n++ {
		status := Retrievestatus(sid)
		if status == "Completed" {
			fmt.Println("Call andata a buon fine")
			break
		}
		fmt.Println(status)
		time.Sleep(10 * time.Second)
	}
	fmt.Println("Call andata male")
	os.Exit(1)
}

//Retrievestatus trova lo status di una call
func Retrievestatus(sid string) (status string) {

	accountSid, err := recuperavariabile("TWILIOACCOUNTSID")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(101)
	}

	url := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Calls/" + sid

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Basic QUM2MTU1NWQ2NDYyODE2NjAxMWM4YzU3NzZhM2JlOTU3ZTo1NDliNGRjOTQ5NmQ3MDg1YTA1M2FkZjQwNzBhOWFkYQ==")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Postman-Token", "decb8b3e-3689-4de0-bba9-d84c74fd0bf7")

	res, errres := http.DefaultClient.Do(req)
	if errres != nil {
		log.Fatal(errres)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	//fmt.Println(res)
	//fmt.Println(string(body))

	//Creo tipo per estrarre singolo valore da file XML
	type TwilioResponse struct {
		Status string `xml:"Call>Status"`
	}

	v := TwilioResponse{}
	errstat := xml.Unmarshal(body, &v)
	if errstat != nil {
		fmt.Printf("error: %v", err)
		return
	}

	//fmt.Printf("Status: %s\n", v.Status)

	return v.Status
}
