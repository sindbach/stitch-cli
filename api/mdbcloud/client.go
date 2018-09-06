package mdbcloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sindbach/stitch-cli/utils"

	"github.com/edaniels/digest"
)

var errCommonServerError = fmt.Errorf("an unexpected server error has occurred")

type projectResponse struct {
	Results []Project `json:"results"`
}

// Project represents a mongodb atlas project
type Project struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	OrgID           string `json:"orgId"`
	ReplicaSetCount int    `json:"replicaSetCount"`
	ShardCount      int    `json:"shardCount"`
}

type userResponse struct {
	Results []User `json:"results"`
}

type User struct {
	ID        string     `json:"id"`
	Email     string     `json:"emailAddress"`
	Firstname string     `json:"firstName"`
	Lastname  string     `json:"lastName"`
	Username  string     `json:"username"`
	Roles     []UserRole `json:"roles"`
}

type UserRole struct {
	OrgID     string `json:"orgId"`
	ProjectID string `json:"groupId"`
	Name      string `json:"roleName"`
}

type orgResponse struct {
	Results []Org `json:"results"`
}

// Org represents a mongodb atlas organizations
type Org struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Client provides access to the MongoDB Cloud Manager APIs
type Client interface {
	WithAuth(username string, apiKey string) Client
	Orgs() ([]Org, error)
	UserByName(string) (*User, error)
	Projects() ([]Project, error)
	ProjectByID(string) (*Project, error)
	ProjectByName(string) (*Project, error)
	DeleteDatabaseUser(projectID string, username string) error
}

type simpleClient struct {
	transport       *digest.Transport
	atlasAPIBaseURL string
}

// NewClient constructs and returns a new Client given a username, API key,
// the public Cloud API base URL, and the atlas API base url
func NewClient(atlasAPIBaseURL string) Client {
	return &simpleClient{
		atlasAPIBaseURL: atlasAPIBaseURL,
	}
}

func (client simpleClient) WithAuth(username, apiKey string) Client {
	// digest.NewTransport will use http.DefaultTransport
	client.transport = digest.NewTransport(username, apiKey)
	return &client
}

// Orgs returns all available Orgs for the user
func (client *simpleClient) Orgs() ([]Org, error) {
	resp, err := client.do(
		http.MethodGet,
		fmt.Sprintf("%s/api/atlas/v1.0/orgs", client.atlasAPIBaseURL),
		nil,
		true,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch available Organization IDs: %s", resp.Status)
	}

	dec := json.NewDecoder(resp.Body)
	var orgResponse orgResponse
	if err := dec.Decode(&orgResponse); err != nil {
		return nil, err
	}

	return orgResponse.Results, nil
}

// Projects returns all available Projects for the user
func (client *simpleClient) Projects() ([]Project, error) {
	resp, err := client.do(
		http.MethodGet,
		fmt.Sprintf("%s/api/public/v1.0/groups", client.atlasAPIBaseURL),
		nil,
		true,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch available Project IDs: %s", resp.Status)
	}

	dec := json.NewDecoder(resp.Body)
	var projectResponse projectResponse
	if err := dec.Decode(&projectResponse); err != nil {
		return nil, err
	}

	return projectResponse.Results, nil
}

// ProjectByID returns info of a Project for the user
func (client *simpleClient) ProjectByID(projectID string) (*Project, error) {
	var projectResponse Project
	err := client.SingleFetch(
		fmt.Sprintf("%s/api/public/v1.0/groups/%s", client.atlasAPIBaseURL, projectID),
		fmt.Sprintf("failed to find information for ProjectID [%s]", projectID),
		fmt.Sprintf("failed to fetch ProjectID [%s]", projectID),
		&projectResponse,
	)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return &projectResponse, nil
}

func (client *simpleClient) ProjectByName(projectName string) (*Project, error) {
	var projectResponse Project
	err := client.SingleFetch(
		fmt.Sprintf("%s/api/public/v1.0/groups/byName/%s", client.atlasAPIBaseURL, projectName),
		fmt.Sprintf("failed to find Project by name [%s]", projectName),
		fmt.Sprintf("failed to fetch Project by name [%s]", projectName),
		&projectResponse,
	)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return &projectResponse, nil
}

func (client *simpleClient) UserByName(userName string) (*User, error) {
	var userResponse User
	err := client.SingleFetch(
		fmt.Sprintf("%s/api/atlas/v1.0/users/byName/%s", client.atlasAPIBaseURL, userName),
		fmt.Sprintf("failed to find User by name [%s]", userName),
		fmt.Sprintf("failed to fetch User by name [%s]", userName),
		&userResponse,
	)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return &userResponse, nil
}

func (client *simpleClient) SingleFetch(url string, notFoundMsg string, failedMsg string, response interface{}) error {
	resp, err := client.do(
		http.MethodGet,
		fmt.Sprintf(url),
		nil,
		true,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf(notFoundMsg)
		}
		return fmt.Errorf("%s: %s", failedMsg, resp.Status)
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(response); err != nil {
		return err
	}
	return nil
}

func (client *simpleClient) do(
	method, url string, // nolint: unparam
	body interface{},
	needAuth bool, // nolint: unparam
) (*http.Response, error) {

	var bodyReader io.Reader
	if body != nil {
		md, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(md)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, errCommonServerError
	}

	if body != nil {
		req.Header.Add("Content-Type", string(utils.MediaTypeJSON))
	}

	req.Header.Add("User-Agent", "MongoDB-Stitch-CLI")

	cl := http.Client{}
	cl.Timeout = time.Second * 5
	if client.transport == nil {
		if needAuth {
			return nil, errors.New("expected to have auth context")
		}
		return cl.Do(req)
	}
	cl.Transport = client.transport

	resp, err := cl.Do(req)
	if err != nil {
		return nil, errCommonServerError
	}

	if resp.StatusCode == http.StatusUnauthorized {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to authenticate with MongoDB Cloud API")
	}

	return resp, nil
}

// DeleteDatabaseUser deletes the database user with the provided username
func (client *simpleClient) DeleteDatabaseUser(projectID, username string) error {
	resp, err := client.do(
		http.MethodDelete,
		fmt.Sprintf("%s/api/atlas/v1.0/groups/%s/databaseUsers/admin/%s",
			client.atlasAPIBaseURL,
			projectID,
			username,
		),
		nil,
		true,
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("error deleting database user '%s'", username)
	}
	return nil
}
