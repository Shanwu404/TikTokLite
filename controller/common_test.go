package controller

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func SendRequest(method string, url string, Body io.Reader) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://www.apifox.cn)")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println(string(body))
}
