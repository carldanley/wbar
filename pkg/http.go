package pkg

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func NewAPIURL(postfix string) string {
	return fmt.Sprintf("%s/%s", GetWyzeBridgeHost(), postfix)
}

func PerformGetRequest(url string, deadlineSeconds int) ([]byte, error) {
	if deadlineSeconds == 0 {
		deadlineSeconds = 1
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(deadlineSeconds))

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)

	if err != nil {
		return []byte{}, fmt.Errorf("could not create new request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return []byte{}, fmt.Errorf("could not fetch new request: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("could not read response body: %w", err)
	}

	return body, nil
}
