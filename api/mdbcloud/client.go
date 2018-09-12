package mdbcloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sindbach/stitch-cli/utils"

	"github.com/edaniels/digest"
)

var errCommonServerError = fmt.Errorf("an unexpected server error has occurred")

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

type databaseUserResponse struct {
	Results []DatabaseUser `json:"results"`
}

type DatabaseUser struct {
	DatabaseName string         `json:"databaseName"`
	Username     string         `json:"username"`
	Roles        []DatabaseRole `json:"roles"`
}

type DatabaseRole struct {
	DatabaseName   string `json:"databaseName"`
	CollectionName string `json:"collectionName"`
	Name           string `json:"roleName"`
}

type orgResponse struct {
	Results []Org `json:"results"`
}

// Org represents a mongodb atlas organizations
type Org struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

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

type clusterResponse struct {
	Results []Cluster `json:"results"`
}

type Cluster struct {
	Autoscaling              ClusterAutoScaling     `json:"autoScaling"`
	BackupEnabled            bool                   `json:"backupEnabled"`
	BIConnector              ClusterBIConnector     `json:"biConnector"`
	ClusterType              string                 `json:"clusterType"`
	DiskSizeGB               float32                `json:"diskSizeGB"`
	EncryptionAtRestProvider string                 `json:"encryptionAtRestProvider"`
	GroupID                  string                 `json:"groupId"`
	ID                       string                 `json:"id"`
	MongoDBVersion           string                 `json:"mongoDBVersion"`
	MongoDBMajorVersion      string                 `json:"MongoDBMajorVersion"`
	MongoURI                 string                 `json:"mongoURI"`
	MongoURIUpdated          string                 `json:"mongoURIUpdated"`
	MongoURIWithOptions      string                 `json:"mongoURIWithOptions"`
	Name                     string                 `json:"name"`
	NumShards                int32                  `json:"numShards"`
	Paused                   bool                   `json:"paused"`
	ProviderBackupEnabled    bool                   `json:"providerBackupEnabled"`
	ProviderSettings         ClusterProviderSetting `json:"providerSettings"`
	ReplicationFactor        int32                  `json:"replicationFactor"`
	StateName                string                 `json:"stateName"`
}

type ClusterAutoScaling struct {
	DiskGBEnabled bool
}

type ClusterBIConnector struct {
	Enabled        bool   `json:"enabled"`
	ReadPreference string `json:""`
}

type ClusterProviderSetting struct {
	Name                string  `json:"providerName"`
	BackingProviderName string  `json:"backingProviderName"`
	RegionName          string  `json:"regionName"`
	InstanceSizeName    string  `json:"instanceSizeName"`
	DiskIOPS            float32 `json:"diskIOPS"`
	DiskTypeName        string  `json:"diskTypeName"`
	EncryptEBSVolume    bool    `json:"encryptEBSVolume"`
}

type processResponse struct {
	Results []Process `json:"results"`
}

// Process represents a mongodb atlas process
type Process struct {
	ID             string    `json:"id"`
	Hostname       string    `json:"hostname"`
	Port           int32     `json:"port"`
	TypeName       string    `json:"typeName"`
	ReplicasetName string    `json:"replicaSetName"`
	ShardName      string    `json:"shardName"`
	Version        string    `json:"version"`
	Lastping       time.Time `json:"lastPing"`
	Created        time.Time `json:"created"`
}

// Client provides access to the MongoDB Cloud Manager APIs
type Client interface {
	WithAuth(username string, apiKey string) Client
	Orgs() ([]Org, error)
	OrgByID(string) (*Org, error)
	UserByName(string) (*User, error)
	Projects() ([]Project, error)
	ProjectByOrgID(string) ([]Project, error)
	ProjectByID(string) (*Project, error)
	ProjectByName(string) (*Project, error)
	ClustersByProjectID(string) ([]Cluster, error)
	DBUsersByProjectID(string) ([]DatabaseUser, error)
	ProcessByProjectID(string) ([]Process, error)
	LogByProcessID(string, string, string) error
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
		0,
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

func (client *simpleClient) OrgByID(id string) (*Org, error) {
	var response Org
	err := client.SingleFetch(
		fmt.Sprintf("%s/api/atlas/v1.0/orgs/%s", client.atlasAPIBaseURL, id),
		fmt.Sprintf("failed to find Org by ID [%s]", id),
		fmt.Sprintf("failed to fetch Org by ID [%s]", id),
		&response,
	)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return &response, nil
}

// Projects returns all available Projects for the user
func (client *simpleClient) Projects() ([]Project, error) {
	resp, err := client.do(
		http.MethodGet,
		fmt.Sprintf("%s/api/public/v1.0/groups", client.atlasAPIBaseURL),
		nil,
		true,
		0,
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
	var response Project
	err := client.SingleFetch(
		fmt.Sprintf("%s/api/public/v1.0/groups/%s", client.atlasAPIBaseURL, projectID),
		fmt.Sprintf("failed to find information for ProjectID [%s]", projectID),
		fmt.Sprintf("failed to fetch ProjectID [%s]", projectID),
		&response,
	)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return &response, nil
}

func (client *simpleClient) ProjectByOrgID(orgID string) ([]Project, error) {
	var response projectResponse
	err := client.SingleFetch(
		fmt.Sprintf("%s/api/atlas/v1.0/orgs/%s/groups", client.atlasAPIBaseURL, orgID),
		fmt.Sprintf("failed to find Projects using Org ID [%s]", orgID),
		fmt.Sprintf("failed to fetch Projects using Org ID [%s]", orgID),
		&response,
	)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return response.Results, nil
}

func (client *simpleClient) ProjectByName(projectName string) (*Project, error) {
	var response Project
	err := client.SingleFetch(
		fmt.Sprintf("%s/api/public/v1.0/groups/byName/%s", client.atlasAPIBaseURL, projectName),
		fmt.Sprintf("failed to find Project by name [%s]", projectName),
		fmt.Sprintf("failed to fetch Project by name [%s]", projectName),
		&response,
	)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return &response, nil
}

func (client *simpleClient) ClustersByProjectID(projectID string) ([]Cluster, error) {
	var response clusterResponse
	err := client.SingleFetch(
		fmt.Sprintf("%s/api/atlas/v1.0/groups/%s/clusters", client.atlasAPIBaseURL, projectID),
		fmt.Sprintf("failed to find Clusters using Project ID [%s]", projectID),
		fmt.Sprintf("failed to fetch Clusters using Project ID [%s]", projectID),
		&response,
	)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return response.Results, nil
}

func (client *simpleClient) DBUsersByProjectID(projectID string) ([]DatabaseUser, error) {
	var response databaseUserResponse
	err := client.SingleFetch(
		fmt.Sprintf("%s/api/atlas/v1.0/groups/%s/databaseUsers", client.atlasAPIBaseURL, projectID),
		fmt.Sprintf("failed to find DB users using Project ID [%s]", projectID),
		fmt.Sprintf("failed to fetch DB users using Project ID [%s]", projectID),
		&response,
	)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return response.Results, nil
}

func (client *simpleClient) ProcessByProjectID(projectID string) ([]Process, error) {
	var response processResponse
	err := client.SingleFetch(
		fmt.Sprintf("%s/api/atlas/v1.0/groups/%s/processes", client.atlasAPIBaseURL, projectID),
		fmt.Sprintf("failed to find processes using Project ID [%s]", projectID),
		fmt.Sprintf("failed to fetch Project using Project ID [%s]", projectID),
		&response,
	)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return response.Results, nil
}

func (client *simpleClient) LogByProcessID(projectID string, processID string, logType string) error {
	err := client.SingleDownload(
		fmt.Sprintf("%s/api/atlas/v1.0/groups/%s/clusters/%s/logs/%s", client.atlasAPIBaseURL, projectID, processID, logType),
		fmt.Sprintf("failed to find log using Process ID [%s]", processID),
		fmt.Sprintf("failed to fetch log using Process ID [%s]", processID),
		fmt.Sprintf("log_%s_%s", processID, logType),
	)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}

func (client *simpleClient) UserByName(userName string) (*User, error) {
	var response User
	err := client.SingleFetch(
		fmt.Sprintf("%s/api/atlas/v1.0/users/byName/%s", client.atlasAPIBaseURL, userName),
		fmt.Sprintf("failed to find User by name [%s]", userName),
		fmt.Sprintf("failed to fetch User by name [%s]", userName),
		&response,
	)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return &response, nil
}

func (client *simpleClient) SingleDownload(url string, notFoundMsg string, failedMsg string, filename string) error {
	output, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer output.Close()

	resp, err := client.do(
		http.MethodGet,
		fmt.Sprintf(url),
		nil,
		true,
		time.Second*300,
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
	fmt.Println("copying")
	// Writer the body to file
	_, err = io.Copy(output, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func (client *simpleClient) SingleFetch(url string, notFoundMsg string, failedMsg string, response interface{}) error {
	resp, err := client.do(
		http.MethodGet,
		fmt.Sprintf(url),
		nil,
		true,
		0,
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
	timeout time.Duration,
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
	if timeout != 0 {
		cl.Timeout = timeout
	} else {
		cl.Timeout = time.Second * 5
	}
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
		0,
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
