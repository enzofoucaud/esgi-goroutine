package msx

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

func GetAccountTransactionCount(address string) (int, error) {
	//Encode the data
	var (
		url = MSX_API + "accounts/" + address + "/transactions/count"
	)

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Err(err).Msg("Error when client create GET request to " + url)
		return 0, err
	}

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Err(err).Msg("Error when client do request to " + url)
		return 0, err
	}
	defer resp.Body.Close()

	// Check error code
	if resp.StatusCode != http.StatusOK {
		var errorElrond ErrorElrond
		body, _ := io.ReadAll(resp.Body)
		_ = json.Unmarshal(body, &errorElrond)
		return 0, errors.New(errorElrond.Message)
	}

	// Convert text/html to int
	var count int
	err = json.NewDecoder(resp.Body).Decode(&count)
	if err != nil {
		log.Err(err).Msg("Error when Decode JSON")
		return 0, err
	}

	return count, nil
}
