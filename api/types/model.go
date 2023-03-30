package types

import (
	"time"

	"github.com/rs/zerolog"
)

type SessionStatus string

const (
	RuntimeProviderKey               = "Provider"
	StatusInitializing SessionStatus = "initializing"
	StatusRunning      SessionStatus = "running"
	StatusTerminated   SessionStatus = "terminated"
	StatusError        SessionStatus = "error"
)

type NoOpLogHook struct{}

func (d NoOpLogHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {}

var NewErrLogHook = func() zerolog.Hook { return NoOpLogHook{} }

// Provider is the platform that the node is spawned on. This is where the user
// code runs
type Provider string

func (r Provider) String() string {
	return string(r)
}

const (
	LambdaLabsProvider Provider = "lambdalabs"
	UnweaveProvider    Provider = "unweave"
)

func (r Provider) DisplayName() string {
	switch r {
	case LambdaLabsProvider:
		return "LambdaLabs"
	case UnweaveProvider:
		return "Unweave"
	default:
		return "Unknown"
	}
}

type Build struct {
	BuildID     string     `json:"buildID"`
	Name        string     `json:"name"`
	ProjectID   string     `json:"projectID"`
	Status      string     `json:"status"`
	BuilderType string     `json:"builderType"`
	CreatedAt   time.Time  `json:"createdAt"`
	StartedAt   *time.Time `json:"startedAt,omitempty"`
	FinishedAt  *time.Time `json:"finishedAt,omitempty"`
}

type LogEntry struct {
	TimeStamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

type NodeSpecs struct {
	VCPUs int `json:"vCPUs"`
	// Memory is the RAM in GB
	Memory int `json:"memory"`
	// GPUMemory is the GPU RAM in GB
	GPUMemory *int `json:"gpuMemory"`
}

type NodeType struct {
	ID          string    `json:"id"`
	Name        *string   `json:"name"`
	Price       *int      `json:"price"`
	Regions     []string  `json:"regions"`
	Description *string   `json:"description"`
	Provider    Provider  `json:"provider"`
	Specs       NodeSpecs `json:"specs"`
}

type Node struct {
	ID       string        `json:"id"`
	TypeID   string        `json:"typeID"`
	Region   string        `json:"region"`
	KeyPair  SSHKey        `json:"sshKeyPair"`
	Status   SessionStatus `json:"status"`
	Provider Provider      `json:"provider"`
}

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SSHKey struct {
	Name      string     `json:"name"`
	PublicKey *string    `json:"publicKey,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
}

type ConnectionInfo struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
}

type Session struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	SSHKeys    []SSHKey        `json:"sshKeys"`
	Connection *ConnectionInfo `json:"connection,omitempty"`
	Status     SessionStatus   `json:"status"`
	NodeID     string          `json:"nodeID"`
}

type Exec struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	SSHKey     SSHKey          `json:"sshKey"`
	Connection *ConnectionInfo `json:"connection,omitempty"`
	Status     SessionStatus   `json:"status"`
	CreatedAt  *time.Time      `json:"createdAt,omitempty"`
	NodeTypeID string          `json:"nodeTypeID"`
	Region     string          `json:"region"`
	Provider   Provider        `json:"provider"`
	Ctx        ExecCtx         `json:"ctx"`
}

type ExecCtx struct {
	Command []string `json:"command"`
	BuildID *string  `json:"buildID,omitempty"`
}
