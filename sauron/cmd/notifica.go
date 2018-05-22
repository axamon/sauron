// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/axamon/sms"

	"github.com/axamon/reperibili"
	"github.com/spf13/cobra"
)

// notificaCmd represents the notifica command
var notificaCmd = &cobra.Command{
	Use:   "notifica",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		f := cmd.Flag("f").Value.String()
		//fmt.Println(f)
		if len(args[0]) > 3 {
			fmt.Println("Argomento superiore a 3 lettere a che piattaforma ti riferisci?")
			os.Exit(1)
		}
		contatto, err := reperibili.Reperibiliperpiattaforma2(args[0], f)

		if err != nil {
			fmt.Println("errore", err.Error())
			//os.Exit(1)
		}
		fmt.Println(contatto.Cellulare)

		testo := fmt.Sprintf("%s sarai reperibile il giorno %s per %s ", contatto.Nome, contatto.Assegnazione.Giorno, contatto.Assegnazione.Piattaforma)
		sms.Inviasms(contatto.Cellulare, "+17372041296", testo)

		//fmt.Println("notifica called")
	},
}

func init() {
	reperibileCmd.AddCommand(notificaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// notificaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// notificaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
