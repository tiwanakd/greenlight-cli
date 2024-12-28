package cmd

import (
	"fmt"

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
		fmt.Printf("New Movie\nTitle: %s\nYear: %d\nRuntime: %s\nGenres: %v\n", movieTitle, movieYear, movieRuntime, movieGenres)
		return nil
	},
}

func init() {
	movieCreateCmd.Flags().StringVarP(&movieTitle, "title", "t", "", "provie the movie title")
	movieCreateCmd.Flags().Int32VarP(&movieYear, "year", "y", 0, "provie the movie release year")
	movieCreateCmd.Flags().StringVarP(&movieRuntime, "runtime", "r", "", "provie the movie runtime (e.g: 120 mins)")
	movieCreateCmd.Flags().StringArrayVarP(&movieGenres, "genres", "g", nil, "provie the movie genres")
	movieCreateCmd.MarkFlagRequired("title")
	movieCreateCmd.MarkFlagRequired("year")
	movieCreateCmd.MarkFlagRequired("runtime")
	movieCreateCmd.MarkFlagRequired("genres")

	movieCmd.AddCommand(movieCreateCmd)
}
