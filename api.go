package main

import (
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
)

func getLikes(username string) ([]string, error) {
	_, body, err := fasthttp.Get(nil, cfg.LikesServiceURL+"user/"+username)
	if err != nil {
		return nil, err
	}

	var res likesResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	if res.Error != "" {
		return nil, errors.New(res.Error)
	}

	for i, j := 0, len(res.Likes)-1; i < j; i, j = i+1, j-1 {
		res.Likes[i], res.Likes[j] = res.Likes[j], res.Likes[i]
	}

	return res.Likes, nil
}

func getDownloadURL(id string) (string, error) {
	_, body, err := fasthttp.Get(nil, cfg.DownloadServiceURL+id)
	if err != nil {
		return "", err
	}

	var res downloadResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}

	return res.DownloadURL, nil
}

type likesResponse struct {
	Likes []string `json:"likes"`
	Error string   `json:"error"`
}

type downloadResponse struct {
	DownloadURL              string `json:"nowatermark"`
	WithWatermarkDownloadURL string `json:"watermark"`
}
