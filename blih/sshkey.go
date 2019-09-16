package blih

import (
	"encoding/json"
	"fmt"
)

type SSHKeyService struct {
	s *Service
}

type SSHKey struct {
	SSHKey string `json:"sshkey"`
}

type KeyListResponse map[string]string

func NewSSHKeyService(s *Service) *SSHKeyService {
	return &SSHKeyService{s: s}
}

func (s *SSHKeyService) List() (*KeyListResponse, error) {
	resp, err := s.s.SendRequest("/sshkeys", "GET", nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(
			"failed to list SSH keys. Request failed with code %d (%s)", resp.StatusCode, resp.Status)
	}

	var message KeyListResponse
	err = json.NewDecoder(resp.Body).Decode(&message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (s *SSHKeyService) Create(key *SSHKey) (*ApiMessageResponse, error) {
	resp, err := s.s.SendRequest("/sshkeys", "POST", key)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(
			"failed to create SSH key. Request failed with code %d (%s)", resp.StatusCode, resp.Status)
	}

	var message ApiMessageResponse
	err = json.NewDecoder(resp.Body).Decode(&message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (s *SSHKeyService) Delete(name string) (*ApiMessageResponse, error) {
	resp, err := s.s.SendRequest("/sshkey/"+name, "DELETE", nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(
			"failed to delete SSH key. Request failed with code %d (%s)", resp.StatusCode, resp.Status)
	}

	var message ApiMessageResponse
	err = json.NewDecoder(resp.Body).Decode(&message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}
