package pkg

import (
	"os"
	"strconv"
)

const DefaultIntervalCheckSeconds = 300

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
