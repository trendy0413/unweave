// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: queries.sql

package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

const BuildCreate = `-- name: BuildCreate :one
insert into unweave.build (project_id, builder_type, name, created_by, started_at)
values ($1, $2, $3, $4, case
                            when $5::timestamptz = '0001-01-01 00:00:00 UTC'::timestamptz
                                then now()
                            else $5::timestamptz end)
returning id
`

type BuildCreateParams struct {
	ProjectID   string    `json:"projectID"`
	BuilderType string    `json:"builderType"`
	Name        string    `json:"name"`
	CreatedBy   string    `json:"createdBy"`
	StartedAt   time.Time `json:"startedAt"`
}

func (q *Queries) BuildCreate(ctx context.Context, arg BuildCreateParams) (string, error) {
	row := q.db.QueryRowContext(ctx, BuildCreate,
		arg.ProjectID,
		arg.BuilderType,
		arg.Name,
		arg.CreatedBy,
		arg.StartedAt,
	)
	var id string
	err := row.Scan(&id)
	return id, err
}

const BuildGet = `-- name: BuildGet :one
select id, name, project_id, builder_type, status, created_by, created_at, started_at, finished_at, updated_at, meta_data
from unweave.build
where id = $1
`

func (q *Queries) BuildGet(ctx context.Context, id string) (UnweaveBuild, error) {
	row := q.db.QueryRowContext(ctx, BuildGet, id)
	var i UnweaveBuild
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ProjectID,
		&i.BuilderType,
		&i.Status,
		&i.CreatedBy,
		&i.CreatedAt,
		&i.StartedAt,
		&i.FinishedAt,
		&i.UpdatedAt,
		&i.MetaData,
	)
	return i, err
}

const BuildGetUsedBy = `-- name: BuildGetUsedBy :many
select s.id, s.name, s.node_id, s.region, s.created_by, s.created_at, s.ready_at, s.exited_at, s.status, s.project_id, s.ssh_key_id, s.connection_info, s.error, s.build_id, s.spec, s.commit_id, s.git_remote_url, s.command, n.provider
from (select id from unweave.build as ub where ub.id = $1) as b
         join unweave.session s
              on s.build_id = b.id
         join unweave.node as n on s.node_id = n.id
`

type BuildGetUsedByRow struct {
	ID             string               `json:"id"`
	Name           string               `json:"name"`
	NodeID         string               `json:"nodeID"`
	Region         string               `json:"region"`
	CreatedBy      string               `json:"createdBy"`
	CreatedAt      time.Time            `json:"createdAt"`
	ReadyAt        sql.NullTime         `json:"readyAt"`
	ExitedAt       sql.NullTime         `json:"exitedAt"`
	Status         UnweaveSessionStatus `json:"status"`
	ProjectID      string               `json:"projectID"`
	SshKeyID       sql.NullString       `json:"sshKeyID"`
	ConnectionInfo json.RawMessage      `json:"connectionInfo"`
	Error          sql.NullString       `json:"error"`
	BuildID        sql.NullString       `json:"buildID"`
	Spec           json.RawMessage      `json:"spec"`
	CommitID       sql.NullString       `json:"commitID"`
	GitRemoteUrl   sql.NullString       `json:"gitRemoteUrl"`
	Command        []string             `json:"command"`
	Provider       string               `json:"provider"`
}

func (q *Queries) BuildGetUsedBy(ctx context.Context, id string) ([]BuildGetUsedByRow, error) {
	rows, err := q.db.QueryContext(ctx, BuildGetUsedBy, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []BuildGetUsedByRow
	for rows.Next() {
		var i BuildGetUsedByRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.NodeID,
			&i.Region,
			&i.CreatedBy,
			&i.CreatedAt,
			&i.ReadyAt,
			&i.ExitedAt,
			&i.Status,
			&i.ProjectID,
			&i.SshKeyID,
			&i.ConnectionInfo,
			&i.Error,
			&i.BuildID,
			&i.Spec,
			&i.CommitID,
			&i.GitRemoteUrl,
			pq.Array(&i.Command),
			&i.Provider,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const BuildUpdate = `-- name: BuildUpdate :exec
update unweave.build
set status      = $2,
    meta_data   = $3,
    started_at  = coalesce(
            nullif($4::timestamptz, '0001-01-01 00:00:00 UTC'::timestamptz),
            started_at),
    finished_at = coalesce(
            nullif($5::timestamptz, '0001-01-01 00:00:00 UTC'::timestamptz),
            finished_at)
where id = $1
`

type BuildUpdateParams struct {
	ID         string             `json:"id"`
	Status     UnweaveBuildStatus `json:"status"`
	MetaData   json.RawMessage    `json:"metaData"`
	StartedAt  time.Time          `json:"startedAt"`
	FinishedAt time.Time          `json:"finishedAt"`
}

func (q *Queries) BuildUpdate(ctx context.Context, arg BuildUpdateParams) error {
	_, err := q.db.ExecContext(ctx, BuildUpdate,
		arg.ID,
		arg.Status,
		arg.MetaData,
		arg.StartedAt,
		arg.FinishedAt,
	)
	return err
}

const MxSessionGet = `-- name: MxSessionGet :one

select s.id,
       s.name,
       s.status,
       s.node_id,
       n.provider,
       s.region,
       s.created_at,
       s.connection_info,
       ssh_key.name       as ssh_key_name,
       ssh_key.public_key,
       ssh_key.created_at as ssh_key_created_at
from unweave.session as s
         join unweave.ssh_key on s.ssh_key_id = ssh_key.id
         join unweave.node as n on s.node_id = n.id
where s.id = $1
`

type MxSessionGetRow struct {
	ID              string               `json:"id"`
	Name            string               `json:"name"`
	Status          UnweaveSessionStatus `json:"status"`
	NodeID          string               `json:"nodeID"`
	Provider        string               `json:"provider"`
	Region          string               `json:"region"`
	CreatedAt       time.Time            `json:"createdAt"`
	ConnectionInfo  json.RawMessage      `json:"connectionInfo"`
	SshKeyName      string               `json:"sshKeyName"`
	PublicKey       string               `json:"publicKey"`
	SshKeyCreatedAt time.Time            `json:"sshKeyCreatedAt"`
}

// -----------------------------------------------------------------
// The queries below return data in the format expected by the API.
// -----------------------------------------------------------------
func (q *Queries) MxSessionGet(ctx context.Context, id string) (MxSessionGetRow, error) {
	row := q.db.QueryRowContext(ctx, MxSessionGet, id)
	var i MxSessionGetRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Status,
		&i.NodeID,
		&i.Provider,
		&i.Region,
		&i.CreatedAt,
		&i.ConnectionInfo,
		&i.SshKeyName,
		&i.PublicKey,
		&i.SshKeyCreatedAt,
	)
	return i, err
}

const MxSessionsGet = `-- name: MxSessionsGet :many
select s.id,
       s.name,
       s.status,
       s.node_id,
       n.provider,
       s.region,
       s.created_at,
       s.connection_info,
       ssh_key.name       as ssh_key_name,
       ssh_key.public_key,
       ssh_key.created_at as ssh_key_created_at
from unweave.session as s
         join unweave.ssh_key on s.ssh_key_id = ssh_key.id
         join unweave.node as n on s.node_id = n.id
where s.project_id = $1
`

type MxSessionsGetRow struct {
	ID              string               `json:"id"`
	Name            string               `json:"name"`
	Status          UnweaveSessionStatus `json:"status"`
	NodeID          string               `json:"nodeID"`
	Provider        string               `json:"provider"`
	Region          string               `json:"region"`
	CreatedAt       time.Time            `json:"createdAt"`
	ConnectionInfo  json.RawMessage      `json:"connectionInfo"`
	SshKeyName      string               `json:"sshKeyName"`
	PublicKey       string               `json:"publicKey"`
	SshKeyCreatedAt time.Time            `json:"sshKeyCreatedAt"`
}

func (q *Queries) MxSessionsGet(ctx context.Context, projectID string) ([]MxSessionsGetRow, error) {
	rows, err := q.db.QueryContext(ctx, MxSessionsGet, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MxSessionsGetRow
	for rows.Next() {
		var i MxSessionsGetRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Status,
			&i.NodeID,
			&i.Provider,
			&i.Region,
			&i.CreatedAt,
			&i.ConnectionInfo,
			&i.SshKeyName,
			&i.PublicKey,
			&i.SshKeyCreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const NodeCreate = `-- name: NodeCreate :exec
select unweave.insert_node(
               $1,
               $2,
               $3,
               $4 :: jsonb,
               $5,
               $6,
               $7 :: text[]
           )
`

type NodeCreateParams struct {
	ID        string          `json:"id"`
	Provider  string          `json:"provider"`
	Region    string          `json:"region"`
	Metadata  json.RawMessage `json:"metadata"`
	Status    string          `json:"status"`
	OwnerID   string          `json:"ownerID"`
	SshKeyIds []string        `json:"sshKeyIds"`
}

func (q *Queries) NodeCreate(ctx context.Context, arg NodeCreateParams) error {
	_, err := q.db.ExecContext(ctx, NodeCreate,
		arg.ID,
		arg.Provider,
		arg.Region,
		arg.Metadata,
		arg.Status,
		arg.OwnerID,
		pq.Array(arg.SshKeyIds),
	)
	return err
}

const NodeStatusUpdate = `-- name: NodeStatusUpdate :exec
update unweave.node
set status = $2,
    ready_at = coalesce($3, ready_at),
    terminated_at = coalesce($4, terminated_at)
where id = $1
`

type NodeStatusUpdateParams struct {
	ID           string       `json:"id"`
	Status       string       `json:"status"`
	ReadyAt      sql.NullTime `json:"readyAt"`
	TerminatedAt sql.NullTime `json:"terminatedAt"`
}

func (q *Queries) NodeStatusUpdate(ctx context.Context, arg NodeStatusUpdateParams) error {
	_, err := q.db.ExecContext(ctx, NodeStatusUpdate,
		arg.ID,
		arg.Status,
		arg.ReadyAt,
		arg.TerminatedAt,
	)
	return err
}

const ProjectGet = `-- name: ProjectGet :one
select id, default_build_id
from unweave.project
where id = $1
`

func (q *Queries) ProjectGet(ctx context.Context, id string) (UnweaveProject, error) {
	row := q.db.QueryRowContext(ctx, ProjectGet, id)
	var i UnweaveProject
	err := row.Scan(&i.ID, &i.DefaultBuildID)
	return i, err
}

const SSHKeyAdd = `-- name: SSHKeyAdd :exec
insert into unweave.ssh_key (owner_id, name, public_key)
values ($1, $2, $3)
`

type SSHKeyAddParams struct {
	OwnerID   string `json:"ownerID"`
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
}

func (q *Queries) SSHKeyAdd(ctx context.Context, arg SSHKeyAddParams) error {
	_, err := q.db.ExecContext(ctx, SSHKeyAdd, arg.OwnerID, arg.Name, arg.PublicKey)
	return err
}

const SSHKeyGetByName = `-- name: SSHKeyGetByName :one
select id, name, owner_id, created_at, public_key, is_active
from unweave.ssh_key
where name = $1
  and owner_id = $2
`

type SSHKeyGetByNameParams struct {
	Name    string `json:"name"`
	OwnerID string `json:"ownerID"`
}

func (q *Queries) SSHKeyGetByName(ctx context.Context, arg SSHKeyGetByNameParams) (UnweaveSshKey, error) {
	row := q.db.QueryRowContext(ctx, SSHKeyGetByName, arg.Name, arg.OwnerID)
	var i UnweaveSshKey
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.OwnerID,
		&i.CreatedAt,
		&i.PublicKey,
		&i.IsActive,
	)
	return i, err
}

const SSHKeyGetByPublicKey = `-- name: SSHKeyGetByPublicKey :one
select id, name, owner_id, created_at, public_key, is_active
from unweave.ssh_key
where public_key = $1
  and owner_id = $2
`

type SSHKeyGetByPublicKeyParams struct {
	PublicKey string `json:"publicKey"`
	OwnerID   string `json:"ownerID"`
}

func (q *Queries) SSHKeyGetByPublicKey(ctx context.Context, arg SSHKeyGetByPublicKeyParams) (UnweaveSshKey, error) {
	row := q.db.QueryRowContext(ctx, SSHKeyGetByPublicKey, arg.PublicKey, arg.OwnerID)
	var i UnweaveSshKey
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.OwnerID,
		&i.CreatedAt,
		&i.PublicKey,
		&i.IsActive,
	)
	return i, err
}

const SSHKeysGet = `-- name: SSHKeysGet :many
select id, name, owner_id, created_at, public_key, is_active
from unweave.ssh_key
where owner_id = $1
`

func (q *Queries) SSHKeysGet(ctx context.Context, ownerID string) ([]UnweaveSshKey, error) {
	rows, err := q.db.QueryContext(ctx, SSHKeysGet, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UnweaveSshKey
	for rows.Next() {
		var i UnweaveSshKey
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.OwnerID,
			&i.CreatedAt,
			&i.PublicKey,
			&i.IsActive,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const SessionCreate = `-- name: SessionCreate :one
insert into unweave.session (node_id, created_by, project_id, ssh_key_id,
                             region, name, connection_info, commit_id, git_remote_url, command, build_id)
values ($1, $2, $3, (select id
                     from unweave.ssh_key as ssh_keys
                     where ssh_keys.name = $11
                       and owner_id = $2), $4, $5, $6, $7, $8, $9, $10)
returning id
`

type SessionCreateParams struct {
	NodeID         string          `json:"nodeID"`
	CreatedBy      string          `json:"createdBy"`
	ProjectID      string          `json:"projectID"`
	Region         string          `json:"region"`
	Name           string          `json:"name"`
	ConnectionInfo json.RawMessage `json:"connectionInfo"`
	CommitID       sql.NullString  `json:"commitID"`
	GitRemoteUrl   sql.NullString  `json:"gitRemoteUrl"`
	Command        []string        `json:"command"`
	BuildID        sql.NullString  `json:"buildID"`
	SshKeyName     string          `json:"sshKeyName"`
}

func (q *Queries) SessionCreate(ctx context.Context, arg SessionCreateParams) (string, error) {
	row := q.db.QueryRowContext(ctx, SessionCreate,
		arg.NodeID,
		arg.CreatedBy,
		arg.ProjectID,
		arg.Region,
		arg.Name,
		arg.ConnectionInfo,
		arg.CommitID,
		arg.GitRemoteUrl,
		pq.Array(arg.Command),
		arg.BuildID,
		arg.SshKeyName,
	)
	var id string
	err := row.Scan(&id)
	return id, err
}

const SessionGet = `-- name: SessionGet :one
select id, name, node_id, region, created_by, created_at, ready_at, exited_at, status, project_id, ssh_key_id, connection_info, error, build_id, spec, commit_id, git_remote_url, command
from unweave.session
where id = $1
`

func (q *Queries) SessionGet(ctx context.Context, id string) (UnweaveSession, error) {
	row := q.db.QueryRowContext(ctx, SessionGet, id)
	var i UnweaveSession
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.NodeID,
		&i.Region,
		&i.CreatedBy,
		&i.CreatedAt,
		&i.ReadyAt,
		&i.ExitedAt,
		&i.Status,
		&i.ProjectID,
		&i.SshKeyID,
		&i.ConnectionInfo,
		&i.Error,
		&i.BuildID,
		&i.Spec,
		&i.CommitID,
		&i.GitRemoteUrl,
		pq.Array(&i.Command),
	)
	return i, err
}

const SessionGetAllActive = `-- name: SessionGetAllActive :many
select id, name, node_id, region, created_by, created_at, ready_at, exited_at, status, project_id, ssh_key_id, connection_info, error, build_id, spec, commit_id, git_remote_url, command
from unweave.session
where status = 'initializing'
   or status = 'running'
`

func (q *Queries) SessionGetAllActive(ctx context.Context) ([]UnweaveSession, error) {
	rows, err := q.db.QueryContext(ctx, SessionGetAllActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UnweaveSession
	for rows.Next() {
		var i UnweaveSession
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.NodeID,
			&i.Region,
			&i.CreatedBy,
			&i.CreatedAt,
			&i.ReadyAt,
			&i.ExitedAt,
			&i.Status,
			&i.ProjectID,
			&i.SshKeyID,
			&i.ConnectionInfo,
			&i.Error,
			&i.BuildID,
			&i.Spec,
			&i.CommitID,
			&i.GitRemoteUrl,
			pq.Array(&i.Command),
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const SessionSetError = `-- name: SessionSetError :exec
update unweave.session
set status = 'error'::unweave.session_status,
    error  = $2
where id = $1
`

type SessionSetErrorParams struct {
	ID    string         `json:"id"`
	Error sql.NullString `json:"error"`
}

func (q *Queries) SessionSetError(ctx context.Context, arg SessionSetErrorParams) error {
	_, err := q.db.ExecContext(ctx, SessionSetError, arg.ID, arg.Error)
	return err
}

const SessionStatusUpdate = `-- name: SessionStatusUpdate :exec
update unweave.session
set status = $2,
    ready_at = coalesce($3, ready_at),
    exited_at = coalesce($4, exited_at)
where id = $1
`

type SessionStatusUpdateParams struct {
	ID       string               `json:"id"`
	Status   UnweaveSessionStatus `json:"status"`
	ReadyAt  sql.NullTime         `json:"readyAt"`
	ExitedAt sql.NullTime         `json:"exitedAt"`
}

func (q *Queries) SessionStatusUpdate(ctx context.Context, arg SessionStatusUpdateParams) error {
	_, err := q.db.ExecContext(ctx, SessionStatusUpdate,
		arg.ID,
		arg.Status,
		arg.ReadyAt,
		arg.ExitedAt,
	)
	return err
}

const SessionUpdateConnectionInfo = `-- name: SessionUpdateConnectionInfo :exec
update unweave.session
set connection_info = $2
where id = $1
`

type SessionUpdateConnectionInfoParams struct {
	ID             string          `json:"id"`
	ConnectionInfo json.RawMessage `json:"connectionInfo"`
}

func (q *Queries) SessionUpdateConnectionInfo(ctx context.Context, arg SessionUpdateConnectionInfoParams) error {
	_, err := q.db.ExecContext(ctx, SessionUpdateConnectionInfo, arg.ID, arg.ConnectionInfo)
	return err
}

const SessionsGet = `-- name: SessionsGet :many
select session.id, ssh_key.name as ssh_key_name, session.status
from unweave.session
         left join unweave.ssh_key
                   on ssh_key.id = session.ssh_key_id
where project_id = $1
order by unweave.session.created_at desc
limit $2 offset $3
`

type SessionsGetParams struct {
	ProjectID string `json:"projectID"`
	Limit     int32  `json:"limit"`
	Offset    int32  `json:"offset"`
}

type SessionsGetRow struct {
	ID         string               `json:"id"`
	SshKeyName sql.NullString       `json:"sshKeyName"`
	Status     UnweaveSessionStatus `json:"status"`
}

func (q *Queries) SessionsGet(ctx context.Context, arg SessionsGetParams) ([]SessionsGetRow, error) {
	rows, err := q.db.QueryContext(ctx, SessionsGet, arg.ProjectID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SessionsGetRow
	for rows.Next() {
		var i SessionsGetRow
		if err := rows.Scan(&i.ID, &i.SshKeyName, &i.Status); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
