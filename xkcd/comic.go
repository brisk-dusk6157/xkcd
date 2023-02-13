package xkcd

import (
	"encoding/json"
)

type Comic struct {
	Num        int
	SafeTitle  string `json:"safe_title"`
	Year       string
	Img        string
	Alt        string
	Transcript string
}

func Parse(b []byte) (*Comic, error) {
	var comic Comic
	err := json.Unmarshal(b, &comic)
	if err != nil {
		return nil, err
	}
	return &comic, nil
}
