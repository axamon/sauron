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

	"github.com/axamon/reperibili"

	"github.com/spf13/cobra"
)

// reperibileCmd represents the reperibile command
var reperibileCmd = &cobra.Command{
	Use:   "reperibile",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		contatto, err := reperibili.Reperibiliperpiattaforma2(args[0], args[1])
		if args[1] == "" {
			fmt.Println("devi passare il file della reperibilità")
			os.Exit(1)
		}
		if err != nil {
			fmt.Println("errore", err.Error())
			//os.Exit(1)
		}
		fmt.Println(contatto)

		//fmt.Println("reperibile called")
	},
}

func init() {
	rootCmd.AddCommand(reperibileCmd)

	//Assegnazione è la variabile con i dati relativi alla ruota di reperibilità

	//limite delle 7 fino alle 7 del mattino seguente il reperibile che viene visualizzato è quello del giorno prima

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reperibileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reperibileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
