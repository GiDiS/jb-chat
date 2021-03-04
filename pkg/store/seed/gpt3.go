package seed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const SberClourGtp3 = "https://api.sbercloud.ru/v2/aicloud/gpt3"

func GetAnswer(ctx context.Context, question string) (string, error) {
	req, err := makeGtp3Req(question)
	if err != nil {
		return "", err
	}
	req = req.WithContext(ctx)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 200 && resp.StatusCode <= 201 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		data := make(map[string]string, 0)
		err = json.Unmarshal(body, &data)
		if err != nil {
			return "", err
		}
		if text, ok := data["data"]; ok {
			return text, nil
		}
	} else {
		return "", fmt.Errorf("invalid response status code: %d", resp.StatusCode)
	}
	return "", errors.New("gpt3 ask fails")
}

func makeGtp3Req(question string) (*http.Request, error) {
	raw, err := json.Marshal(question)
	if err != nil {
		return nil, err
	}
	ask := `{"question":` + string(raw) + `}`
	req, err := http.NewRequest("POST", SberClourGtp3, strings.NewReader(ask))
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:84.0) Gecko/20100101 Firefox/84.0)")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Referer", "https://sbercloud.ru/ru/warp/gpt-3")
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Origin", "https://sbercloud.ru")

	return req, nil
}
