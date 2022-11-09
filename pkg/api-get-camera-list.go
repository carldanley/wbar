package pkg

import (
	"encoding/json"
	"fmt"
	"time"
)

const LastImageTimeDivider = 1000
const DeadlineCameraList = 1

type APICameraResponse struct {
	Cameras map[string]APICameraResponseCamera `json:"cameras"`
}

type APICameraResponseCamera struct {
	Name          string `json:"name_uri"`
	Connected     bool   `json:"connected"`
	Enabled       bool   `json:"enabled"`
	LastImageTime int64  `json:"img_time"`
}

func GetCameraList() (APICameraResponse, error) {
	body, err := PerformGetRequest(NewAPIURL("/api"), DeadlineCameraList)

	if err != nil {
		return APICameraResponse{}, err
	}

	var jsonData APICameraResponse
	err = json.Unmarshal(body, &jsonData)

	if err != nil {
		return APICameraResponse{}, fmt.Errorf("could not unmarshal body: %w", err)
	}

	return jsonData, nil
}

func (c *APICameraResponseCamera) GetLastImageTime() time.Time {
	return time.Unix(c.LastImageTime/LastImageTimeDivider, 0)
}
