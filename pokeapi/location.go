package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type PokemonApiLocationResult struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// Api calls
func PokemonApiLocationCall(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed HTTP request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, errors.New("Error with ReadAll")
	}

	return body, nil
}

func UnmarshalToLocationResult(body []byte) (PokemonApiLocationResult, error) {
	var locations PokemonApiLocationResult
	if err := json.Unmarshal(body, &locations); err != nil {
		return locations, errors.New("Error unmarshaling json")
	}

	return locations, nil
}
