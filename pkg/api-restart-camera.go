package pkg

import (
	"fmt"
	"time"
)

func RestartCamera(name string) error {
	err := stopCamera(name)
	if err != nil {
		return err
	}

	time.Sleep(time.Second)

	err = startCamera(name)
	if err != nil {
		return err
	}

	return nil
}

func stopCamera(name string) error {
	url := fmt.Sprintf("/api/%s/stop", name)
	_, err := PerformGetRequest(NewAPIURL(url), DeadlineRefreshCameraImage)

	if err != nil {
		return fmt.Errorf("could not stop camera: %w", err)
	}

	return nil
}

func startCamera(name string) error {
	url := fmt.Sprintf("/api/%s/start", name)
	_, err := PerformGetRequest(NewAPIURL(url), DeadlineRefreshCameraImage)

	if err != nil {
		return fmt.Errorf("could not start camera: %w", err)
	}

	return nil
}
