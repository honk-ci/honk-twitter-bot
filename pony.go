package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	ponyURL = "https://theponyapi.com/api/v1/pony/random"
)

// Only the properties we actually use.
type ponyResult struct {
	Pony ponyResultPony `json:"pony"`
}

type ponyResultPony struct {
	Representations ponyRepresentations `json:"representations"`
}

type ponyRepresentations struct {
	Small string `json:"small"`
}

func getPony() (ponyImage []byte) {
	for i := 0; i < 5; i++ {
		image, err := readPony()
		if err != nil {
			continue
		}
		return image
	}

	return nil
}

func readPony() (ponyImage []byte, err error) {
	client := http.Client{}
	resp, err := client.Get(ponyURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("no pony found")
	}
	var a ponyResult
	if err = json.NewDecoder(resp.Body).Decode(&a); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	smallMeow, err := client.Get(a.Pony.Representations.Small)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer smallMeow.Body.Close()
	if smallMeow.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("no Meow found")
	}

	f, err := ioutil.ReadAll(smallMeow.Body)
	return f, err
}
