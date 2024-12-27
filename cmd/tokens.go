package cmd

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/tiwanakd/greenlight-cli/client"
)

var tokenCmd = &cobra.Command{
	Use:   "tokens",
	Short: "use for generating tokens",
	Long: `This command allows for generation for various tokens for authentication.
Tokens can be generated for user activation, autnetication and password resets`,
}

var activationTokenCmd = &cobra.Command{
	Use:   "activation [email]",
	Short: "generate a activation token for a user",
	Long: `Generate a new activation for a user who is not activated.
A user may be able to request a new token in case the current token in exipired
or they not recevie the intial activation token for some reason`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var email string
		flag := cmd.Flag("email")
		if len(args) < 1 && !flag.Changed {
			return fmt.Errorf("email is required in order to request a activation token")
		}

		if len(args) > 0 {
			email = args[0]
		} else {
			email = userEmail
		}

		jsonMap := client.NewJSONMap()
		jsonMap.Add("email", email)

		jsReader, err := jsonMap.CreateJSONReader()
		if err != nil {
			return err
		}

		err, code, _, body := apiClient.PostRequest("/v1/tokens/activation", jsReader)
		if err != nil {
			return err
		}

		if code != http.StatusAccepted {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true
			return errors.New(body)
		}

		fmt.Println(body)
		return nil
	},
}

func init() {
	activationTokenCmd.Flags().StringVarP(&userEmail, "email", "e", "", "email of registered user")

	tokenCmd.AddCommand(activationTokenCmd)
}
