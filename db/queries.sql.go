// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: queries.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const ProjectCreate = `-- name: ProjectCreate :one
insert into unweave.projects (name, owner_id)
values ($1, $2)
returning id
`

type ProjectCreateParams struct {
	Name    string    `json:"name"`
	OwnerID uuid.UUID `json:"ownerID"`
}

func (q *Queries) ProjectCreate(ctx context.Context, arg ProjectCreateParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, ProjectCreate, arg.Name, arg.OwnerID)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const ProjectGet = `-- name: ProjectGet :one
select id, name, owner_id, created_at
from unweave.projects
where id = $1
`

func (q *Queries) ProjectGet(ctx context.Context, id uuid.UUID) (UnweaveProject, error) {
	row := q.db.QueryRowContext(ctx, ProjectGet, id)
	var i UnweaveProject
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.OwnerID,
		&i.CreatedAt,
	)
	return i, err
}

const SSHKeyAdd = `-- name: SSHKeyAdd :exec
insert INTO unweave.ssh_keys (owner_id, name, public_key)
values ($1, $2, $3)
`

type SSHKeyAddParams struct {
	OwnerID   uuid.UUID `json:"ownerID"`
	Name      string    `json:"name"`
	PublicKey string    `json:"publicKey"`
}

func (q *Queries) SSHKeyAdd(ctx context.Context, arg SSHKeyAddParams) error {
	_, err := q.db.ExecContext(ctx, SSHKeyAdd, arg.OwnerID, arg.Name, arg.PublicKey)
	return err
}

const SSHKeyGetByName = `-- name: SSHKeyGetByName :one
select id, name, owner_id, created_at, public_key
from unweave.ssh_keys
where name = $1
`

func (q *Queries) SSHKeyGetByName(ctx context.Context, name string) (UnweaveSshKey, error) {
	row := q.db.QueryRowContext(ctx, SSHKeyGetByName, name)
	var i UnweaveSshKey
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.OwnerID,
		&i.CreatedAt,
		&i.PublicKey,
	)
	return i, err
}

const SSHKeyGetByPublicKey = `-- name: SSHKeyGetByPublicKey :one
select id, name, owner_id, created_at, public_key
from unweave.ssh_keys
where public_key = $1
`

func (q *Queries) SSHKeyGetByPublicKey(ctx context.Context, publicKey string) (UnweaveSshKey, error) {
	row := q.db.QueryRowContext(ctx, SSHKeyGetByPublicKey, publicKey)
	var i UnweaveSshKey
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.OwnerID,
		&i.CreatedAt,
		&i.PublicKey,
	)
	return i, err
}

const SSHKeysGet = `-- name: SSHKeysGet :many
select id, name, owner_id, created_at, public_key
from unweave.ssh_keys
where owner_id = $1
`

func (q *Queries) SSHKeysGet(ctx context.Context, ownerID uuid.UUID) ([]UnweaveSshKey, error) {
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
insert into unweave.sessions (node_id, created_by, project_id, runtime, ssh_key_id)
values ($1, $2, $3, $4, (select id
                         from unweave.ssh_keys
                         where name = $5 and owner_id = $2))
returning id
`

type SessionCreateParams struct {
	NodeID     string    `json:"nodeID"`
	CreatedBy  uuid.UUID `json:"createdBy"`
	ProjectID  uuid.UUID `json:"projectID"`
	Runtime    string    `json:"runtime"`
	SshKeyName string    `json:"sshKeyName"`
}

func (q *Queries) SessionCreate(ctx context.Context, arg SessionCreateParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, SessionCreate,
		arg.NodeID,
		arg.CreatedBy,
		arg.ProjectID,
		arg.Runtime,
		arg.SshKeyName,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const SessionGet = `-- name: SessionGet :one
select id, node_id, created_by, created_at, ready_at, exited_at, status, project_id, runtime, ssh_key_id
from unweave.sessions
where id = $1
`

func (q *Queries) SessionGet(ctx context.Context, id uuid.UUID) (UnweaveSession, error) {
	row := q.db.QueryRowContext(ctx, SessionGet, id)
	var i UnweaveSession
	err := row.Scan(
		&i.ID,
		&i.NodeID,
		&i.CreatedBy,
		&i.CreatedAt,
		&i.ReadyAt,
		&i.ExitedAt,
		&i.Status,
		&i.ProjectID,
		&i.Runtime,
		&i.SshKeyID,
	)
	return i, err
}

const SessionSetTerminated = `-- name: SessionSetTerminated :exec
update unweave.sessions
set status = unweave.session_status('terminated')
where id = $1
`

func (q *Queries) SessionSetTerminated(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, SessionSetTerminated, id)
	return err
}

const SessionsGet = `-- name: SessionsGet :many
select sessions.id, ssh_keys.name as ssh_key_name, sessions.status
from unweave.sessions
         left join unweave.ssh_keys
                   on ssh_keys.id = sessions.ssh_key_id
where project_id = $1
order by unweave.sessions.created_at desc
limit $2 offset $3
`

type SessionsGetParams struct {
	ProjectID uuid.UUID `json:"projectID"`
	Limit     int32     `json:"limit"`
	Offset    int32     `json:"offset"`
}

type SessionsGetRow struct {
	ID         uuid.UUID            `json:"id"`
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
