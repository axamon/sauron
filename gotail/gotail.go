package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/axamon/reperibili"

	"github.com/hpcloud/tail"
)

var nagioslog = flag.String("f", "/var/log/nagios/nagios.log", "Nagios file di log")

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
		notificabool := strings.Contains(line.Text, "NOTIFICATION")
		switch notificabool {

		case true:
			fmt.Println(line.Text)
			reperibile, _ := reperibili.Reperibiliperpiattaforma2("CDN", "reperibilita.csv")

			TO := reperibile.Cellulare
			NOME := reperibile.Nome

			reperibili.Chiamareperibile(TO, NOME)

		default:
			//fmt.Println("debug")
			continue
		}
	}
}
