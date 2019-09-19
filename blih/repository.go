package blih

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultRepositoryType = "git"
)

type DateTime struct {
	time.Time
}

type BoolString bool

func (bs BoolString) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), "\"")

	bs = str == "True"
	return nil
}

func (d *DateTime) UnmarshalJSON(b []byte) error {
	i, err := strconv.ParseInt(strings.Trim(string(b), "\""), 10, 64)
	if err != nil {
		return err
	}

	*d = DateTime{time.Unix(i, 0)}
	return nil
}

type RepositoryService struct {
	s *Service
}

type RepositoryACL struct {
	User string `json:"user"`
	ACL  string `json:"acl"`
}

// A structure which represents a Blih repository.
type Repository struct {
	// The repository name
	Name string `json:"name"`

	// The repository type (usually "git")
	Type string `json:"type"`

	// An optional description
	Description string `json:"description,omitempty"`
}

type RepositoryListEntry struct {
	UUID string `json:"uuid"`
	URL  string `json:"url"`
}

type RepositoryListResponse struct {
	Message      string                         `json:"message"`
	Repositories map[string]RepositoryListEntry `json:"repositories"`
}

type RepositoryInfosResponse struct {
	Message struct {
		URL          string     `json:"url"`
		UUID         string     `json:"uuid"`
		Public       BoolString `json:"public,string"`
		Description  string     `json:"description"`
		CreationTime DateTime   `json:"creation_time,string"`
	}
}

type RepositoryACLResponse map[string]string

func NewRepositoryService(s *Service) *RepositoryService {
	return &RepositoryService{s: s}
}

func (a *RepositoryACL) Validate() bool {
	str := strings.ToLower(a.ACL)

	for _, l := range str {
		if l != 'r' && l != 'w' && l != 'a' {
			return false
		}
	}

	return true
}

func (s *RepositoryService) Create(repo *Repository) (*ApiMessageResponse, error) {
	if repo.Type == "" {
		repo.Type = DefaultRepositoryType
	}

	resp, err := s.s.SendRequest("/repositories", "POST", repo)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(
			"failed to create repository. Request failed with code %d (%s)", resp.StatusCode, resp.Status)
	}

	var message ApiMessageResponse
	err = json.NewDecoder(resp.Body).Decode(&message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (s *RepositoryService) SetACL(name string, acl *RepositoryACL) (*ApiMessageResponse, error) {
	if !acl.Validate() {
		return nil, fmt.Errorf("ACL should only contains a, r or w characters")
	}

	resp, err := s.s.SendRequest("/repository/"+name+"/acls", "POST", acl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(
			"failed to set ACL. Request failed with code %d (%s)", resp.StatusCode, resp.Status)
	}

	var message ApiMessageResponse
	err = json.NewDecoder(resp.Body).Decode(&message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (s *RepositoryService) GetACL(name string) (*RepositoryACLResponse, error) {
	resp, err := s.s.SendRequest("/repository/"+name+"/acls", "GET", nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(
			"failed to get ACL. Request failed with code %d (%s)", resp.StatusCode, resp.Status)
	}

	var message RepositoryACLResponse
	err = json.NewDecoder(resp.Body).Decode(&message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (s *RepositoryService) Infos(name string) (*RepositoryInfosResponse, error) {
	resp, err := s.s.SendRequest("/repository/"+name, "GET", nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(
			"failed to get repository infos. Request failed with code %d (%s)", resp.StatusCode, resp.Status)
	}

	var message RepositoryInfosResponse
	err = json.NewDecoder(resp.Body).Decode(&message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (s *RepositoryService) List() (*RepositoryListResponse, error) {
	resp, err := s.s.SendRequest("/repositories", "GET", nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(
			"failed to list repositories. Request failed with code %d (%s)", resp.StatusCode, resp.Status)
	}

	var message RepositoryListResponse
	err = json.NewDecoder(resp.Body).Decode(&message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (s *RepositoryService) Delete(name string) (*ApiMessageResponse, error) {
	resp, err := s.s.SendRequest("/repository/"+name, "DELETE", nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(
			"failed to delete repository. Request failed with code %d (%s)", resp.StatusCode, resp.Status)
	}

	var message ApiMessageResponse
	err = json.NewDecoder(resp.Body).Decode(&message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}
