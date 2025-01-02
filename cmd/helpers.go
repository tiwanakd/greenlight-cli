package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/tiwanakd/greenlight-cli/client"
)

type JsonMap struct {
	kv map[string]any
}

func newJSONMap() *JsonMap {
	return &JsonMap{make(map[string]any)}
}

func (j *JsonMap) add(key string, value any) {
	if _, ok := j.kv[key]; !ok {
		j.kv[key] = value
	}
}

func (j *JsonMap) createJSON() ([]byte, error) {
	js, err := json.Marshal(j.kv)
	if err != nil {
		return nil, err
	}

	return js, err
}

func (j *JsonMap) createJSONReader() (*bytes.Reader, error) {
	js, err := j.createJSON()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(js), nil
}

// create a method that will return an error if the API does not return the desired response code.
// Silence usage and errors when a command returns an error as we do not want the Usage for the command
// to be printed. Dislabing errors will avoid the error body from our API to duplicated on the terminal
// This will allow to return the body as an error the way RunE propery of Cobra.Command expects.
func customError(cmd *cobra.Command, body string) error {
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	return errors.New("ERROR: " + body)
}

func listBody(cmd *cobra.Command, resp client.Response) error {
	if resp.Err != nil {
		return resp.Err
	}

	if resp.Code != http.StatusOK {
		return customError(cmd, resp.Body)
	}

	fmt.Fprintln(os.Stdout, resp.Body)
	return nil
}
