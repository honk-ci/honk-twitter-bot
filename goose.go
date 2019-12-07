package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	gooseURL = "https://api.unsplash.com/photos/random?query=goose&client_id=%s"
)

func getGoose() []byte {
	for i := 0; i < 5; i++ {
		image, err := readGoose()
		if err != nil {
			continue
		}
		return image
	}

	return nil
}

type gooseResult struct {
	GooseURLS gooseImages `json:"urls"`
}

type gooseImages struct {
	Small string `json:"small"`
}

func readGoose() (image []byte, err error) {
	client := http.Client{}
	uri := fmt.Sprintf(gooseURL, Config.UnsplashClientID)
	randomGoose, err := client.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer randomGoose.Body.Close()

	if randomGoose.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("no goose found")
	}
	var a gooseResult
	if err = json.NewDecoder(randomGoose.Body).Decode(&a); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	smallGoose, err := client.Get(a.GooseURLS.Small)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer smallGoose.Body.Close()
	if smallGoose.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("no goose found")
	}

	f, err := ioutil.ReadAll(smallGoose.Body)

	return f, nil
}
