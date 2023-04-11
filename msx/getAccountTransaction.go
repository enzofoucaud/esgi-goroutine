package msx

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

type NftAccount struct {
	GasUsed       int    `json:"gasUsed"`
	Receiver      string `json:"receiver"`
	ReceiverShard int    `json:"receiverShard"`
	Sender        string `json:"sender"`
}

func GetCollectionAccount(account, from, size string) ([]NftAccount, error) {
	//Encode the data
	var (
		url = MSX_API + "accounts/" + account + "/transactions?from=" + from + "&size=" + size
	)

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Err(err).Msg("Error when client create GET request to " + url)
		return []NftAccount{}, err
	}

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Err(err).Msg("Error when client do request to " + url)
		return []NftAccount{}, err
	}
	defer resp.Body.Close()

	// Check error code
	if resp.StatusCode != http.StatusOK {
		var errorElrond ErrorElrond
		body, _ := io.ReadAll(resp.Body)
		_ = json.Unmarshal(body, &errorElrond)
		return []NftAccount{}, errors.New(errorElrond.Message)
	}

	// Convert JSON into struct
	var nftAccount []NftAccount
	err = json.NewDecoder(resp.Body).Decode(&nftAccount)
	if err != nil {
		log.Err(err).Msg("Error when Decode JSON")
		return []NftAccount{}, err
	}

	return nftAccount, nil
}
