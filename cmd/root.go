package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tiwanakd/greenlight-cli/client"
)

var (
	userName     string
	userEmail    string
	userPassword string
)

var apiClient = client.New()

var rootCmd = &cobra.Command{
	Use:   "greenlight",
	Short: "greenlight is client to interact with greenlight-api",
	Long: `This is client that is build to interact with greenlight-api project
the api is deployed at greenlight-api.com this cli-clent provides various 
commands that allows to interract with this api`,
}

var healthCheckCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "check the health of the greenlight api",
	RunE: func(cmd *cobra.Command, args []string) error {
		err, code, _, body := apiClient.GetRequest("/v1/healthcheck")
		if err != nil {
			return err
		}

		fmt.Printf("Code: %d\nBody: %s\n", code, body)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(tokenCmd)
	rootCmd.AddCommand(healthCheckCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
