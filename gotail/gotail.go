package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/axamon/reperibili"

	sms "github.com/axamon/sms"

	"github.com/hpcloud/tail"
)

var nagioslog = flag.String("nagioslog", "/var/log/nagios/nagios.log", "Nagios file di log")
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
		//fmt.Println(line.Text)
		//Se la linea è di notifica la analizza se no passa oltre
		//notificabool := strings.Contains(line.Text, "NOTIFICATION")
		switch {

		case strings.Contains(line.Text, "NOTIFICATION") && strings.Contains(line.Text, "CRITICAL"):
			fmt.Println(line.Text)
			//TODO cambiare CDN con qualcosa di variabile
			reperibile, _ := reperibili.Reperibiliperpiattaforma2("CDN", *reperibilita)

			TO := reperibile.Cellulare
			NOME := reperibile.Nome
			//debug
			fmt.Println(TO, NOME)

			reperibili.Chiamareperibile(TO, NOME)

		case strings.Contains(line.Text, "NOTIFICATION") && strings.Contains(line.Text, "OK"):
			//Se ok allora manda solo sms niente chiamata
			fmt.Println("ricevuto OK")
			reperibile, _ := reperibili.Reperibiliperpiattaforma2("CDN", *reperibilita)

			TO := reperibile.Cellulare
			pezzi := strings.Split(line.Text, ";")
			messaggio := "Su " + pezzi[1] + " servizio " + pezzi[2] + " " + pezzi[3]

			go sms.Inviasms(TO, messaggio)

		default:
			//fmt.Println("debug")
			continue
		}
	}
}
