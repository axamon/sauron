package sms

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

//Verificacellulare si assicura che il cellulare inserito sia nel formato corretto
func Verificacellulare(CELLULARE string) (ok bool) {

	re := regexp.MustCompile(`^\+3[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]$`)
	return re.MatchString(CELLULARE)

}

//Quanto sono compliant con almeno alcune delle 12 best practices di GO!
//https://talks.golang.org/2013/bestpractices.slide#1
func recuperavariabile(variabile string) (result string, err error) {
	if result, ok := os.LookupEnv(variabile); ok && len(result) != 0 {
		return result, nil
	}
	return "", fmt.Errorf("la variabile %s non esiste o è vuota", variabile)
}

//Inviasms invia sms via Twilio
func Inviasms(to, from, body string) {

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

	//TODO vedere se riesce a prendere anche le variabili da ambiente windows...
	//...ma anche no! :)

	//Crea la URL necessaria per richiamare la funzionalità degli SMS di Twilio
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	//Valorizza i campi per l'invio del SMS
	v := url.Values{}
	v.Set("To", to)     //Esempio: "+393357291532"
	v.Set("From", from) //Esempio "+17372041296"
	v.Set("Body", body)

	//impacchettiamo tutte le variabile insieme
	rb := *strings.NewReader(v.Encode())

	//Creiamo un client http
	client := &http.Client{}

	//Creiamo la http request da inviare dopo
	req, err := http.NewRequest("POST", urlStr, &rb)
	if err != nil {
		fmt.Fprintln(os.Stdout, "OH noooo! Qualcosa è andata storta nel creare la richiesta", err)
		os.Exit(103)
	}

	//Utiliziamo l'autenticazione basic
	req.SetBasicAuth(accountSid, authToken)
	//Inseriamo un po' di headers come piacciono a Twilio
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	//Finalmente inviamo la request e salviamo la http response
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stdout, "OH noooo! Qualcosa è andata storta nell'inviare la richiesta", err)
		os.Exit(104)
	}

	//controlliamo che ha da dire la response
	//Restituisce codice e significato, se ricevi 201 CREATED allora è ok.
	fmt.Println(resp.Status)

	//Usciamo con zero che significa tutto ok!
	os.Exit(0)
}
