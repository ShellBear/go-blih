package blih

import (
	"encoding/json"
	"fmt"
	"time"
)

type UtilsService struct {
	s *Service
}

type ApiMessageResponse struct {
	Message string `json:"message"`
}

func NewUtilsService(s *Service) *UtilsService {
	return &UtilsService{s: s}
}

func (s *UtilsService) Ping() (time.Duration, error) {
	start := time.Now()
	resp, err := s.s.SendRequest("/", "GET", nil)
	if err != nil {
		return 0, err
	}
	end := time.Now().Sub(start)

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf(
			"failed to to ping blih. Request failed with code %d (%s)", resp.StatusCode, resp.Status)
	}

	return end, nil
}

func (s *UtilsService) WhoAmI() (*ApiMessageResponse, error) {
	resp, err := s.s.SendRequest("/whoami", "GET", nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(
			"failed to to get identity. Request failed with code %d (%s)", resp.StatusCode, resp.Status)
	}

	var message ApiMessageResponse
	err = json.NewDecoder(resp.Body).Decode(&message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}
