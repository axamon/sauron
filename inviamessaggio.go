package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

//Santa Palomba
var twimlurl = "http://sauron1.westeurope.cloudapp.azure.com:3000/twiml"

//var twimlurl = "https://handler.twilio.com/twiml/EHf9986fbef2c724000473a181c2de9864"

//TwiML parte superiore del file xml che si vuole creare
type TwiML struct {
	XMLName xml.Name `xml:"Response"`
	Say     []Say
}

//Say parte inerna della response
type Say struct {
	Value string `xml:",chardata"`
	Voice string `xml:"voice,attr"`
	Lang  string `xml:"language,attr"`
}

//MESSAGE è il testo da inviare a TWILIO
var MESSAGE = os.Args[1]

//Attrezza tutto per le API
func main() {

	r := mux.NewRouter()
	r.HandleFunc("/twiml", twimlfunc)
	//Il parametro passato dopo call deve essere un cellulare italiano nel formato +39xxxxxxxxxx
	//Se non è ben formattato allora restituisce un 404
	r.HandleFunc("/call/{TO:\\+39\\d{10}}", call)
	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}

//deve far vedere il file XML che voglio io
func twimlfunc(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(MESSAGE)

	twiml := TwiML{Say: []Say{Say{Value: MESSAGE, Lang: "it-IT", Voice: "alice"}}}

	//fmt.Println(twiml)

	x, err := xml.Marshal(twiml)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(x)
}

func recuperavariabile(variabile string) (result string, err error) {
	if result, ok := os.LookupEnv(variabile); ok && len(result) != 0 {
		return result, nil
	}
	return "", fmt.Errorf("la variabile %s non esiste o è vuota", variabile)
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
        v.Set("status_callback", "http://www.myapp.com/events")
        v.Set("status_callback_event", "initiated")
        v.Set("status_callback_method", "POST")
	v.Set("To", vars["TO"])
	v.Set("From", twilionumber)
	//Sfortunatamente la URL deve essere Pubblica se no twilio non può arrivarci
	v.Set("Url", twimlurl)
	rb := *strings.NewReader(v.Encode())

	// Create Client
	client := &http.Client{}

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
		}
	} else {
		fmt.Println(resp.Status)
		w.Write([]byte("Grosso guaio a ChinaTown"))
	}
}
