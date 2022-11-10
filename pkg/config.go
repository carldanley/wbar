package pkg

import (
	"os"
	"strconv"
)

const DefaultIntervalCheckSeconds = 300
const DefaultMaxCameraLagSeconds = 10

func GetWyzeBridgeHost() string {
	return os.Getenv("WYZE_BRIDGE_HOST")
}

func GetIntervalCheckSeconds() int {
	intervalCheckSeconds, exists := os.LookupEnv("INTERVAL_CHECK_SECONDS")

	if exists {
		interval, _ := strconv.Atoi(intervalCheckSeconds)

		if interval > 0 {
			return interval
		}
	}

	return DefaultIntervalCheckSeconds
}

func GetMaxCameraLagSeconds() float64 {
	maxCameraLagSeconds, exists := os.LookupEnv("MAX_CAMERA_LAG_SECONDS")

	if exists {
		interval, _ := strconv.Atoi(maxCameraLagSeconds)

		if interval > 0 {
			return float64(interval)
		}
	}

	return DefaultMaxCameraLagSeconds
}

func GetTimeZone() string {
	tz, exists := os.LookupEnv("TZ")
	if !exists {
		tz = "America/New_York"
	}

	return tz
}
