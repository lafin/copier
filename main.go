package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	api "github.com/lafin/vk"
)

func newfileUploadRequest(from, to string) (*http.Request, error) {
	src, err := http.Get(from)
	if err != nil {
		return nil, err
	}
	defer src.Body.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(from))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, src.Body)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", to, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func main() {
	clientID := os.Getenv("CLIENT_ID")
	email := os.Getenv("CLIENT_EMAIL")
	password := os.Getenv("CLIENT_PASSWORD")

	log.Println("start")
	_, err := api.GetAccessToken(clientID, email, password)
	if err != nil {
		log.Fatalf("[main:api.GetAccessToken] error: %s", err)
		return
	}

	server, err := api.GetUploadServer(117456732)
	if err != nil {
		log.Fatalf("[main:api.GetUploadServer] error: %s", err)
		return
	}

	uploadURL := server.Response.UploadURL
	request, err := newfileUploadRequest("https://pp.userapi.com/c639631/v639631968/14282/yS0K3aa6zEM.jpg", uploadURL)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer resp.Body.Close()

		var fileUploadRequest api.ResponseFileUploadRequest
		if err := json.Unmarshal(body, &fileUploadRequest); err != nil {
			log.Fatal(err)
			return
		}

		result, err := api.SavePhoto(117456732, fileUploadRequest.Server, fileUploadRequest.Photo, fileUploadRequest.Hash)
		if err != nil {
			log.Fatalf("[main:api.SavePhoto] error: %s", err)
			return
		}
		log.Println(result)
	}
}
