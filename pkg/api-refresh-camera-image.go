package pkg

import (
	"fmt"
	"log"
	"time"
)

const DeadlineGetOrRefreshCameraImage = 10

func GetOrRefreshCameraImage(camera string, lastImageTime time.Time) ([]byte, error) {
	// check the time on the snapshot
	diff := time.Since(lastImageTime)

	// if snapshot hasn't been refreshed in awhile, force a refresh
	if diff.Seconds() >= float64(GetIntervalCheckSeconds()) {
		log.Printf("Refreshing camera image: %s...\n", camera)

		url := fmt.Sprintf("/snapshot/%s.jpg", camera)
		body, err := PerformGetRequest(NewAPIURL(url), DeadlineGetOrRefreshCameraImage)

		if err != nil {
			return nil, fmt.Errorf("could not refresh camera image: %w", err)
		}

		return body, nil
	}

	log.Printf("Fetching existing camera image: %s...\n", camera)

	// otherwise, use the existing snapshot
	url := fmt.Sprintf("/img/%s.jpg", camera)
	body, err := PerformGetRequest(NewAPIURL(url), DeadlineGetOrRefreshCameraImage)

	if err != nil {
		return nil, fmt.Errorf("could not refresh camera image: %w", err)
	}

	return body, nil
}
