package pkg

import (
	"fmt"
	"time"
)

const DeadlineRefreshCameraImage = 10

func RefreshCameraImage(camera string, lastImageTime time.Time) ([]byte, error) {
	url := fmt.Sprintf("/snapshot/%s.jpg", camera)
	body, err := PerformGetRequest(NewAPIURL(url), DeadlineRefreshCameraImage)

	if err != nil {
		return nil, fmt.Errorf("could not refresh camera image: %w", err)
	}

	return body, nil
}
