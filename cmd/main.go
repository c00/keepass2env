package main

import (
	"fmt"
	"os"

	"github.com/keepasssecrethelper/config"
	"github.com/keepasssecrethelper/helper"
	"github.com/spf13/cobra"
)

func main() {
	var keyfilePath string
	var outputPath string

	var rootCmd = &cobra.Command{
		Use:   "keepasssecrethelper <database path> <config path>",
		Short: "Extract passwords from keepass to a file",
		Long: `
Extract password from keepass to a file. Useful for seeding secrets 
on a new machine. It will open the database and read out the given 
entries and save them in a .env file. 
It will attempt to update existing entries in the .env file if they
exist.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dbPath := args[0]
			entriesPath := "config.yaml"
			if len(args) > 1 {
				entriesPath = args[1]
			}

			cfg, err := config.FromFile(entriesPath)
			if err != nil {
				return fmt.Errorf("cannot get config file: %w", err)
			}

			// TODO get the goddamn pasword from stdin

			runner := helper.Helper{
				Params: helper.HelperParams{
					KeyfilePath:  keyfilePath,
					DatabasePath: dbPath,
					OutputPath:   outputPath,
					Config:       cfg,
				},
			}

			return runner.Run()
		},
	}

	rootCmd.Flags().StringVarP(&keyfilePath, "keyfile", "k", "", "Path to the keyfile")
	rootCmd.Flags().StringVarP(&outputPath, "out", "o", ".secrets.env", "Output file")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
