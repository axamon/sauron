package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func recuperavariabile(variabile string) (result string, err error) {
	if result, ok := os.LookupEnv(variabile); ok && len(result) != 0 {
		return result, nil
	}
	return "", fmt.Errorf("la variabile %s non esiste o è vuota", variabile)
}

func main() {
	TO := os.Args[1]
	NOME := os.Args[2]

	sid, err := Chiamareperibile(TO, NOME)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(sid)

}

func call(w http.ResponseWriter, r *http.Request) {

	twilionumber, err := recuperavariabile("TWILIONUMBER")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(100)
	}
	// Let's set some initial default variables

	//Recupera l'accountsid di Twilio dallla variabile d'ambiente
	accountSid, err := recuperavariabile("TWILIOACCOUNTSID")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(101)
	}

	//Recupera il token supersegreto dalla variabile d'ambiente
	authToken, err := recuperavariabile("TWILIOAUTHTOKEN")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(102)
	}

	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Calls.json"
	vars := mux.Vars(r)
	// Build out the data for our message
	v := url.Values{}
	//v.Set("status_callback", "http://sauron1.westeurope.cloudapp.azure.com:3000/status")
	//v.Set("status_callback_event", "initiated")
	v.Set("status_callback_method", "POST")
	v.Set("To", vars["TO"])
	v.Set("NOME", "Gringo")
	v.Set("From", twilionumber)
	//Sfortunatamente la URL deve essere Pubblica se no twilio non può arrivarci
	v.Set("Url", "https://handler.twilio.com/twiml/EH5cef42aa1454fc2326780c8f08c6d568")
	rb := *strings.NewReader(v.Encode())

	// Create Client
	client := &http.Client{}

	//Prepara la richiesta HTTP
	req, _ := http.NewRequest("POST", urlStr, &rb)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// make request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("errore", err.Error())
		os.Exit(500)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal(bodyBytes, &data)
		if err == nil {
			fmt.Println(data["sid"])

			fmt.Println(data)
		}

	}
}

//Chiamareperibile e comunica il problema
func Chiamareperibile(TO, NOME string) (sid interface{}, err error) {

	twilionumber, err := recuperavariabile("TWILIONUMBER")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(100)
	}
	// Let's set some initial default variables

	//Recupera l'accountsid di Twilio dallla variabile d'ambiente
	accountSid, err := recuperavariabile("TWILIOACCOUNTSID")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(101)
	}

	//Recupera il token supersegreto dalla variabile d'ambiente
	authToken, err := recuperavariabile("TWILIOAUTHTOKEN")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(102)
	}

	body := strings.NewReader("Url=https://handler.twilio.com/twiml/EH5cef42aa1454fc2326780c8f08c6d568?NOME=" + NOME + "&To=" + TO + "&From=" + twilionumber)
	req, err := http.NewRequest("POST", "https://api.twilio.com/2010-04-01/Accounts/"+accountSid+"/Calls", body)
	if err != nil {
		fmt.Println(err)
	}
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	// make request
	var data map[string]interface{}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {

		bodyBytes, errb := ioutil.ReadAll(resp.Body)
		if errb != nil {
			fmt.Println(errb)
		}
		err := json.Unmarshal(bodyBytes, &data)
		if err != nil {
			//fmt.Println("ok")
			return "", err
		}
	}
	fmt.Println(data)
	sid = data["sid"]
	return sid, nil
}
