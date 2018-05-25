package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/axamon/reperibili"

	sms "github.com/axamon/sms"

	"github.com/hpcloud/tail"
)

const (
	version = "1.0"
)

//Version mostra la versione attuale del software
func Version() {
	fmt.Println(version)
}

//variabile che punta al file log di Nagios
var nagioslog = flag.String("nagioslog", "/var/log/nagios/nagios.log", "Nagios file di log")

//variabile per recuperare lo storage della reperibilità
var reperibilita = flag.String("reperibilita", "$GOPATH/src/github.com/axamon/sauron/sauron/reperibilita.csv", "Nagios file di log")

func main() {
	flag.Parse()

	//Inizia il tail dalla fine del file senza leggerlo dall'inizio
	var fine tail.SeekInfo
	fine.Offset = 0
	fine.Whence = 2

	//MustExist il file deve esistere Follow fa tail -f e ReOpen gestisce il logrotate
	t, err := tail.TailFile(*nagioslog,
		tail.Config{
			Location:  &fine,
			MustExist: true,
			Follow:    true,
			ReOpen:    true,
		})

	if err != nil {
		fmt.Println("errore: ", err.Error())
	}
	//Per ogni nuova linea nel file
	for line := range t.Lines {

		//fmt.Println(line.Text) //per debug

		//Se la linea è di notifica la analizza se no passa oltre
		//notificabool := strings.Contains(line.Text, "NOTIFICATION")
		switch {

		case strings.Contains(line.Text, "NOTIFICATION") && strings.Contains(line.Text, "CRITICAL"):
			fmt.Println(line.Text)
			//TODO cambiare CDN con qualcosa di variabile
			reperibile, _ := reperibili.Reperibiliperpiattaforma2("CDN", *reperibilita)

			TO := reperibile.Cellulare
			NOME := reperibile.Nome
			COGNOME := reperibile.Cognome
			//debug
			fmt.Println(TO, NOME, COGNOME)

			sid, err := reperibili.Chiamareperibile(TO, NOME, COGNOME)
			if err != nil {
				fmt.Println("Errore", err.Error())
			}
			fmt.Println(sid)

			//esce dallo switch
			break

		case strings.Contains(line.Text, "NOTIFICATION") && strings.Contains(line.Text, "OK"):
			//Se ok allora manda solo sms senza chiamata
			fmt.Println("ricevuto OK") //per debug

			reperibile, err := reperibili.Reperibiliperpiattaforma2("CDN", *reperibilita)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}

			TO := reperibile.Cellulare
			pezzi := strings.Split(line.Text, ";")
			messaggio := "Su " + pezzi[1] + " servizio " + pezzi[2] + " " + pezzi[3]

			go sms.Inviasms(TO, messaggio)
			//esce dallo switch
			break

		default:
			//fmt.Println("debug")
			//esce dallo switch
			break
		}
	}
}
