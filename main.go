package copier

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

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

// UploadFiles - upload files
func UploadFiles(files []string, groupID int) ([]string, error) {
	clientID := os.Getenv("CLIENT_ID")
	email := os.Getenv("CLIENT_EMAIL")
	password := os.Getenv("CLIENT_PASSWORD")

	_, err := api.GetAccessToken(clientID, email, password)
	if err != nil {
		return nil, err
	}

	server, err := api.GetUploadServer(groupID)
	if err != nil {
		return nil, err
	}

	var results []string
	uploadURL := server.Response.UploadURL
	for _, file := range files {
		request, err := newfileUploadRequest(file, uploadURL)
		if err != nil {
			return nil, err
		}
		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			return nil, err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var fileUploadRequest api.ResponseFileUploadRequest
		if err := json.Unmarshal(body, &fileUploadRequest); err != nil {
			return nil, err
		}

		result, err := api.SavePhoto(groupID, fileUploadRequest.Server, fileUploadRequest.Photo, fileUploadRequest.Hash)
		if err != nil {
			return nil, err
		}

		for _, item := range result.Response {
			results = append(results, strconv.Itoa(item.OwnerID)+"_"+strconv.Itoa(item.ID)+"_"+item.AccessKey)
		}
	}

	return results, nil
}
