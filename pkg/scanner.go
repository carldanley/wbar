package pkg

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/otiai10/gosseract/v2"
)

const TimestampWidth = 500
const TimestampHeight = 100
const TimestampOffsetX = 570
const TimestampOffsetY = 100

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

		// now, iterate through each camera and check it for restart
		for _, camera := range cameras {
			err := processCamera(camera)
			if err != nil {
				log.Printf("Could not process camera \"%s\": %v", camera.Name, err)
			}

			log.Printf("Processed camera: %s\n", camera.Name)
		}

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
	log.Printf("[%s] Refreshing camera image\n", camera.Name)

	cameraImageBytes, err := RefreshCameraImage(camera.Name, camera.GetLastImageTime())
	if err != nil {
		return fmt.Errorf("could not refresh camera image: %w", err)
	}

	blackedOutImage, err := blackoutNonTextOnImage(cameraImageBytes)
	if err != nil {
		return fmt.Errorf("could not black out non-text on image: %w", err)
	}

	// convert the image to a byte array
	imageByteBuffer := new(bytes.Buffer)

	err = png.Encode(imageByteBuffer, blackedOutImage)
	if err != nil {
		return fmt.Errorf("could not convert image to bytes: %w", err)
	}

	// scan the image for text
	log.Printf("[%s] Scanning camera image for text\n", camera.Name)

	text, err := scanCameraImageForText(imageByteBuffer.Bytes())
	if err != nil {
		return fmt.Errorf("could not scan camera image for text: %w", err)
	}

	// filter text out of string
	regex := regexp.MustCompile(`\d{2}:\d{2}:\d{2}`)
	foundText := regex.FindString(text)

	// house-keeping
	log.Printf("[%s] Text found: %s\n", camera.Name, foundText)
	duration := getTimeDifferenceFromText(foundText)
	log.Printf("[%s] Difference of: %v\n", camera.Name, duration)

	// make sure duration isn't too high
	if duration.Seconds() >= GetMaxCameraLagSeconds() {
		log.Printf("[%s] Restarting Camera\n", camera.Name)

		err = RestartCamera(camera.Name)
		if err != nil {
			return fmt.Errorf("could not restart camera: %w", err)
		}
	}

	return nil
}

func scanCameraImageForText(img []byte) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	err := client.SetImageFromBytes(img)
	if err != nil {
		return "", fmt.Errorf("could not set image from bytes: %w", err)
	}

	text, err := client.Text()
	if err != nil {
		return "", fmt.Errorf("could not parse text: %w", err)
	}

	return text, nil
}

func blackoutNonTextOnImage(imageBytes []byte) (image.Image, error) {
	const ColorValueCutoff = 210

	// convert the image bytes to a go image object
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, fmt.Errorf("could not decode image: %w", err)
	}

	size := img.Bounds().Size()
	rect := image.Rect(0, 0, TimestampWidth, TimestampHeight)
	newImage := image.NewRGBA(rect)

	startX := size.X - TimestampOffsetX
	startY := size.Y - TimestampOffsetY

	// loop though all the x
	for x := startX; x < size.X; x++ {
		// and now loop thorough all of this x's y
		for y := startY; y < size.Y; y++ {
			pixel := img.At(x, y)
			originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)
			newColor := originalColor

			if (originalColor.R < ColorValueCutoff) || (originalColor.G < ColorValueCutoff) || (originalColor.B < ColorValueCutoff) {
				newColor.R = 0
				newColor.G = 0
				newColor.B = 0
			}

			newImage.Set(x-startX, y-startY, newColor)
		}
	}

	return newImage, nil
}

func getTimeDifferenceFromText(text string) time.Duration {
	const ExpectedTimePartsLength = 3

	textParts := strings.Split(text, ":")

	if len(textParts) != ExpectedTimePartsLength {
		return 0
	}

	hour, err := strconv.Atoi(textParts[0])
	if err != nil {
		hour = 0
	}

	minute, err := strconv.Atoi(textParts[1])
	if err != nil {
		minute = 0
	}

	second, err := strconv.Atoi(textParts[2])
	if err != nil {
		second = 0
	}

	location, err := time.LoadLocation("EST")
	if err != nil {
		return 0
	}

	now := time.Now()
	then := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, 0, location)

	return time.Since(then)
}
