package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (r *repository) sendPostRequest(ctx context.Context, data any, url string, out any) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", r.authURL+"/"+url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	request.WithContext(ctx)
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode >= 300 {
		return fmt.Errorf("bad status code: %d", response.StatusCode)
	}

	return json.NewDecoder(response.Body).Decode(&out)
}
