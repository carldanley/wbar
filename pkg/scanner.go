package pkg

import (
	"fmt"
	"log"
	"regexp"
	"sync"
	"time"

	"github.com/otiai10/gosseract/v2"
)

const SecondsToWaitAfterImageRefresh = 5

func StartScanning() {
	for {
		log.Println("Fetching camera list...")

		// fetch the cameras
		response, err := GetCameraList()
		if err != nil {
			log.Printf("Could not fetch camera list: %v\n", err)

			time.Sleep(time.Second * time.Duration(GetIntervalCheckSeconds()))

			continue
		}

		// weed out the cameras that are not connected or not enabled
		log.Printf("Found %d camera(s); filtering for connected and enabled...\n", len(response.Cameras))
		cameras := filterCamerasForConnectedAndEnabled(response)
		log.Printf("Processing %d camera(s)...", len(cameras))

		// create a new wait group for this batch of cameras
		var wg sync.WaitGroup

		// now, iterate through each camera and check it for restart
		for _, camera := range cameras {
			wg.Add(1)

			go func(camera APICameraResponseCamera) {
				defer wg.Done()

				err := processCamera(camera)
				if err != nil {
					log.Printf("Could not process camera \"%s\": %v", camera.Name, err)
				}

				log.Printf("Processed camera: %s\n", camera.Name)
			}(camera)
		}

		wg.Wait()
		log.Println("Cameras scanned; waiting for next loop...")
		time.Sleep(time.Second * time.Duration(GetIntervalCheckSeconds()))
	}
}

func filterCamerasForConnectedAndEnabled(response APICameraResponse) []APICameraResponseCamera {
	tmp := []APICameraResponseCamera{}

	for _, camera := range response.Cameras {
		if camera.Connected && camera.Enabled {
			tmp = append(tmp, camera)
		}
	}

	return tmp
}

func processCamera(camera APICameraResponseCamera) error {
	log.Printf("Processing camera: %s\n", camera.Name)

	cameraImageBytes, err := GetOrRefreshCameraImage(camera.Name, camera.GetLastImageTime())
	if err != nil {
		return fmt.Errorf("could not refresh camera image: %w", err)
	}

	// scan the image for text
	log.Printf("Scanning camera image for text: %s\n", camera.Name)
	log.Printf("bytes found: %d\n", len(cameraImageBytes))

	text, err := scanCameraImageForText(cameraImageBytes)
	if err != nil {
		return fmt.Errorf("could not scan camera image for text: %w", err)
	}

	// filter text out of string
	regex := regexp.MustCompile(`\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}`)
	foundText := regex.FindString(text)

	// make sure we actually have a date
	log.Printf("camera: %s text found: %s\n", camera.Name, foundText)

	return nil
}

func scanCameraImageForText(image []byte) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	err := client.SetImageFromBytes(image)
	if err != nil {
		return "", fmt.Errorf("could not set image from bytes: %w", err)
	}

	text, err := client.Text()
	if err != nil {
		return "", fmt.Errorf("could not parse text: %w", err)
	}

	return text, nil
}
