package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
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

		jsonMap := newJSONMap()
		jsonMap.add("email", email)

		jsReader, err := jsonMap.createJSONReader()
		if err != nil {
			return err
		}

		err, code, _, body := apiClient.NewRequest(http.MethodPost, "/v1/tokens/activation", jsReader, nil)
		if err != nil {
			return err
		}

		if code != http.StatusAccepted {
			return customError(cmd, body)
		}

		fmt.Println(body)
		return nil
	},
}

var authenticationTokenCmd = &cobra.Command{
	Use:   "authentication",
	Short: "generate an authtentication token",
	Long: `this command generates a token for user authtication;
this token allow usage of futher transactions with the api`,
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonMap := newJSONMap()
		jsonMap.add("email", userEmail)
		jsonMap.add("password", userPassword)

		js, err := jsonMap.createJSONReader()
		if err != nil {
			return err
		}

		err, code, _, body := apiClient.NewRequest(http.MethodPost, "/v1/tokens/authentication", js, nil)
		if err != nil {
			return err
		}

		if code != http.StatusCreated {
			return customError(cmd, body)
		}

		fmt.Println(body)
		return nil
	},
}

func init() {
	activationTokenCmd.Flags().StringVarP(&userEmail, "email", "e", "", "email of registered user")

	authenticationTokenCmd.Flags().StringVarP(&userEmail, "email", "e", "", "email of the user to authenticate")
	authenticationTokenCmd.Flags().StringVarP(&userPassword, "password", "p", "", "password of the user to authenticate")
	authenticationTokenCmd.MarkFlagRequired("email")
	authenticationTokenCmd.MarkFlagRequired("password")

	tokenCmd.AddCommand(activationTokenCmd)
	tokenCmd.AddCommand(authenticationTokenCmd)
}
