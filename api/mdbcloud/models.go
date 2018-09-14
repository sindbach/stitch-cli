package mdbcloud

import "time"

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
