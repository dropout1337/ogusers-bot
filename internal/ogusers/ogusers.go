package ogusers

import (
	"encoding/json"
	"io"
	"net/http"
)

func GetUserByUID(uid string) (*Response, error) {
	return request("https://ogusers.gg/0_link.php?uid=" + uid)
}

func GetUserByUsername(username string) (*Response, error) {
	return request("https://ogusers.gg/0_link.php?username=" + username)
}

func request(url string) (*Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var response Response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
