package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	imageURL = "https://api.unsplash.com/photos/random?query=%s&client_id=%s"
)

func getImage(animalType string) []byte {
	for i := 0; i < 5; i++ {
		image, err := readImage(animalType)
		if err != nil {
			continue
		}
		return image
	}

	return nil
}

type imageResult struct {
	ImageURLS imagesSize `json:"urls"`
}

type imagesSize struct {
	Small string `json:"small"`
}

func readImage(animalType string) (image []byte, err error) {
	client := http.Client{}
	uri := fmt.Sprintf(imageURL, animalType, Config.UnsplashClientID)
	randomImage, err := client.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer randomImage.Body.Close()

	if randomImage.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("no goose found")
	}
	var a imageResult
	if err = json.NewDecoder(randomImage.Body).Decode(&a); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	smallImage, err := client.Get(a.ImageURLS.Small)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer smallImage.Body.Close()
	if smallImage.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("no %s found", animalType)
	}

	f, err := ioutil.ReadAll(smallImage.Body)

	return f, nil
}
