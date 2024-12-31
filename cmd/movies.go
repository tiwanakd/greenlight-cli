package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var movieCmd = &cobra.Command{
	Use:   "movies",
	Short: "make requests related to the movies endpoint",
}

var movieCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new movie",
	Long: `create a new movie with the provided details like movie title,
year, runtime and Genres`,
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonMap := newJSONMap()
		jsonMap.add("title", movieTitle)
		jsonMap.add("year", movieYear)
		jsonMap.add("runtime", movieRuntime)
		jsonMap.add("genres", movieGenres)

		js, err := jsonMap.createJSONReader()
		if err != nil {
			return err
		}

		authToken, err := extractAuthToken()
		if authToken.Token == "" || authToken.Expiry.Before(time.Now()) {
			return customError(cmd, "autorization token expired or not found; Please login again using [login] command")
		}

		//Since the "POST /v1/movies" expects a "Authorization: Bearer [token]" header
		//create an empty header map and add the authorizion token
		authorizationHeader := http.Header{}
		authorizationHeader.Add("Authorization", "Bearer "+authToken.Token)

		err, code, _, body := apiClient.NewRequest(http.MethodPost, "/v1/movies", js, authorizationHeader)
		if err != nil {
			return err
		}

		if code != http.StatusCreated {
			return customError(cmd, body)
		}

		fmt.Printf("New Movie\n%s", body)

		return nil
	},
}

var movieListCmd = &cobra.Command{
	Use:   "list",
	Short: "list movies",
	Long: `list the movies as per the flags provided,
if not flags are provided all the movies will be listed`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("list movies")
		return nil
	},
}

func init() {
	movieCreateCmd.Flags().StringVar(&movieTitle, "title", "", "title of the new movie")
	movieCreateCmd.Flags().Int32VarP(&movieYear, "year", "y", 0, "movie release year")
	movieCreateCmd.Flags().StringVarP(&movieRuntime, "runtime", "r", "", "movie runtime (e.g: 120 mins)")
	movieCreateCmd.Flags().StringSliceVarP(&movieGenres, "genres", "g", nil, "movie genres (values should be comma separated with no spaces following commas)")
	movieCreateCmd.MarkFlagRequired("title")
	movieCreateCmd.MarkFlagRequired("year")
	movieCreateCmd.MarkFlagRequired("runtime")
	movieCreateCmd.MarkFlagRequired("genres")

	movieCmd.AddCommand(movieCreateCmd)
	movieCmd.AddCommand(movieListCmd)
}
