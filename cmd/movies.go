package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	movieTitle   string
	movieYear    int32
	movieRuntime string
	movieGenres  []string

	movieID   int64
	movieSort string
)

var movieCmd = &cobra.Command{
	Use:   "movies",
	Short: "make requests related to the movies endpoint",
	Long:  "This command allows to peform CRUD operations on the movies endpoint of the API.",
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
			return errors.New("autorization token expired or not found; Please login again using [login] command")
		}

		//Since the "POST /v1/movies" expects a "Authorization: Bearer [token]" header
		//create an empty header map and add the authorizion token
		authorizationHeader := http.Header{}
		authorizationHeader.Add("Authorization", "Bearer "+authToken.Token)

		err, code, body := apiClient.NewRequest(http.MethodPost, "/v1/movies", js, authorizationHeader)
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
if no flags are provided all the movies will be listed`,
	RunE: func(cmd *cobra.Command, args []string) error {
		authToken, err := extractAuthToken()
		if err != nil {
			return err
		}

		authHeader := http.Header{}
		authHeader.Add("Authorization", "Bearer "+authToken.Token)

		id := cmd.Flag("id")
		title := cmd.Flag("title")
		genres := cmd.Flag("genres")
		sort := cmd.Flag("sort")

		switch {
		case id.Changed:
			err, code, body := apiClient.NewRequest(http.MethodGet, fmt.Sprint("/v1/movies/", movieID), http.NoBody, authHeader)
			return listBody(cmd, err, code, body)
		case title.Changed || genres.Changed || sort.Changed:
			url := fmt.Sprintf("/v1/movies?title=%s&genres=%s&sort=%s", movieTitle, strings.Join(movieGenres, ","), movieSort)
			err, code, body := apiClient.NewRequest(http.MethodGet, url, http.NoBody, authHeader)
			return listBody(cmd, err, code, body)
		default:
			err, code, body := apiClient.NewRequest(http.MethodGet, "/v1/movies", http.NoBody, authHeader)
			return listBody(cmd, err, code, body)
		}

	},
}

var movieUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "update a movie",
	Long:  "update the movie with the provided id",
	RunE: func(cmd *cobra.Command, args []string) error {
		titleFlag := cmd.Flag("title")
		yearFlag := cmd.Flag("year")
		runtimeFlag := cmd.Flag("runtime")
		genresFlag := cmd.Flag("genres")

		jsonMap := newJSONMap()

		if titleFlag.Changed {
			jsonMap.add("title", movieTitle)
		}

		if yearFlag.Changed {
			jsonMap.add("year", movieYear)
		}

		if runtimeFlag.Changed {
			jsonMap.add("runtime", movieRuntime)
		}

		if genresFlag.Changed {
			jsonMap.add("genres", movieGenres)
		}

		js, err := jsonMap.createJSONReader()
		if err != nil {
			return err
		}

		authToken, err := extractAuthToken()
		if err != nil {
			return err
		}

		authHeader := http.Header{}
		authHeader.Add("Authorization", "Bearer "+authToken.Token)

		err, code, body := apiClient.NewRequest(http.MethodPatch, fmt.Sprint("/v1/movies/", movieID), js, authHeader)
		if err != nil {
			return err
		}

		if code != http.StatusOK {
			return customError(cmd, body)
		}

		fmt.Println(body)
		return nil
	},
}

var movieDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a movie",
	Long:  "delete the movie with the provided id",
	RunE: func(cmd *cobra.Command, args []string) error {
		idFlag := cmd.Flag("id")

		if len(args) < 0 && !idFlag.Changed {
			return errors.New("movie id needs to provided")
		}

		var id int64
		if len(args) > 0 {
			id64, err := strconv.ParseInt(args[0], 0, 64)
			if err != nil {
				return err
			}
			id = id64
		} else {
			id = movieID
		}

		authToken, err := extractAuthToken()
		if err != nil {
			return err
		}

		authHeader := http.Header{}
		authHeader.Add("Authorization", "Bearer "+authToken.Token)

		err, code, body := apiClient.NewRequest(http.MethodDelete, fmt.Sprint("/v1/movies/", id), http.NoBody, authHeader)
		if err != nil {
			return err
		}

		if code != http.StatusOK {
			return customError(cmd, body)
		}

		fmt.Println(body)
		return nil
	},
}

func init() {
	movieCreateCmd.Flags().StringVarP(&movieTitle, "title", "t", "", "title of the new movie")
	movieCreateCmd.Flags().Int32VarP(&movieYear, "year", "y", 0, "movie release year")
	movieCreateCmd.Flags().StringVarP(&movieRuntime, "runtime", "r", "", "movie runtime (e.g: 120 mins)")
	movieCreateCmd.Flags().StringSliceVarP(&movieGenres, "genres", "g", nil, "movie genres (values should be comma separated with no spaces following commas)")
	movieCreateCmd.MarkFlagRequired("title")
	movieCreateCmd.MarkFlagRequired("year")
	movieCreateCmd.MarkFlagRequired("runtime")
	movieCreateCmd.MarkFlagRequired("genres")

	movieListCmd.Flags().Int64Var(&movieID, "id", 0, "id of the movie")
	movieListCmd.Flags().StringVarP(&movieTitle, "title", "t", "", "filter title of the movie(s)")
	movieListCmd.Flags().StringSliceVarP(&movieGenres, "genres", "g", nil, "filter by genre(s)")
	movieListCmd.Flags().StringVar(&movieSort, "sort", "", "sort the movies by various fields (e.g. use 'year'/'-year' for ascending/desending)")

	movieUpdateCmd.Flags().Int64Var(&movieID, "id", 0, "id of the movie")
	movieUpdateCmd.Flags().StringVarP(&movieTitle, "title", "t", "", "title of the new movie")
	movieUpdateCmd.Flags().Int32VarP(&movieYear, "year", "y", 0, "movie release year")
	movieUpdateCmd.Flags().StringVarP(&movieRuntime, "runtime", "r", "", "movie runtime (e.g: 120 mins)")
	movieUpdateCmd.Flags().StringSliceVarP(&movieGenres, "genres", "g", nil, "movie genres (values should be comma separated with no spaces following commas)")
	movieUpdateCmd.MarkFlagRequired("id")
	movieUpdateCmd.MarkFlagsOneRequired("title", "year", "runtime", "genres")

	movieDeleteCmd.Flags().Int64Var(&movieID, "id", 0, "delete the movie with the provided id")

	movieCmd.AddCommand(movieCreateCmd)
	movieCmd.AddCommand(movieListCmd)
	movieCmd.AddCommand(movieUpdateCmd)
	movieCmd.AddCommand(movieDeleteCmd)
}
