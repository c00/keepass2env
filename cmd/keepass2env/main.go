package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/keepass2env/config"
	"github.com/keepass2env/runner"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const version = "0.0.1"

func main() {
	var configPath string
	var keyfilePath string
	var outputPath string
	var databasePath string

	var rootCmd = &cobra.Command{
		Use:     "keepass2env",
		Example: "keepass2env -c config.yaml",
		Short:   "Extract passwords from keepass to a file",
		Version: version,
		Long: `
Extract password from keepass to a file. Useful for seeding secrets 
on a new machine. It will open the database and read out the given 
entries and put them in a .env file.

It will add or update entries in the output .env file.`,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var password string

			cfg, err := config.FromFile(configPath)
			if err != nil {
				return fmt.Errorf("cannot get config file: %w", err)
			}

			if keyfilePath == "" && cfg.KeyfilePath != "" {
				keyfilePath = cfg.KeyfilePath
			}

			if outputPath == ".secrets.env" && cfg.OutputPath != "" {
				outputPath = cfg.OutputPath
			}

			if databasePath == "" && cfg.DatabasePath != "" {
				databasePath = cfg.DatabasePath
			}

			if cfg.PasswordEnv != "" {
				password = os.Getenv(cfg.PasswordEnv)
			}

			if password == "" {
				// get the pasword from stdin
				fmt.Print("Enter Keepass Database Password: ")
				bytePassword, err := term.ReadPassword(int(syscall.Stdin))
				if err != nil {
					return fmt.Errorf("cannot read password: %w", err)
				}
				password = string(bytePassword)
				fmt.Println("")
			}

			runner := runner.Helper{
				Params: runner.HelperParams{
					KeyfilePath:      keyfilePath,
					DatabasePath:     databasePath,
					DatabasePassword: password,
					OutputPath:       outputPath,
					Entries:          cfg.Entries,
				},
			}

			return runner.Run()
		},
	}

	rootCmd.Flags().StringVarP(&configPath, "config", "c", "keepass2env.yaml", "Configuration file")
	rootCmd.Flags().StringVarP(&databasePath, "database", "d", "", "Database file")
	rootCmd.Flags().StringVarP(&keyfilePath, "keyfile", "k", "", "Path to the keyfile")
	rootCmd.Flags().StringVarP(&outputPath, "out", "o", ".secrets.env", "Output file")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
