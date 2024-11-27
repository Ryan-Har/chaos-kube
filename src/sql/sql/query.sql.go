// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package sql

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addJob = `-- name: AddJob :one
INSERT INTO jobs (configuration_id, name, description, start_time, end_time, status)
Values ($1, $2, $3, $4, $5, $6)
returning id
`

type AddJobParams struct {
	ConfigurationID pgtype.UUID
	Name            string
	Description     pgtype.Text
	StartTime       pgtype.Timestamptz
	EndTime         pgtype.Timestamptz
	Status          NullJobStatus
}

func (q *Queries) AddJob(ctx context.Context, arg AddJobParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, addJob,
		arg.ConfigurationID,
		arg.Name,
		arg.Description,
		arg.StartTime,
		arg.EndTime,
		arg.Status,
	)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}

const getConfigurationByID = `-- name: GetConfigurationByID :one
SELECT id, name, options, created_at, updated_at FROM configurations
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetConfigurationByID(ctx context.Context, id pgtype.UUID) (Configuration, error) {
	row := q.db.QueryRow(ctx, getConfigurationByID, id)
	var i Configuration
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Options,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getConfigurations = `-- name: GetConfigurations :many
SELECT id, name, options, created_at, updated_at FROM configurations
ORDER BY updated_at DESC
`

func (q *Queries) GetConfigurations(ctx context.Context) ([]Configuration, error) {
	rows, err := q.db.Query(ctx, getConfigurations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Configuration
	for rows.Next() {
		var i Configuration
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Options,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getJobByID = `-- name: GetJobByID :one
SELECT id, configuration_id, name, description, start_time, end_time, status, created_at, updated_at FROM jobs
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetJobByID(ctx context.Context, id pgtype.UUID) (Job, error) {
	row := q.db.QueryRow(ctx, getJobByID, id)
	var i Job
	err := row.Scan(
		&i.ID,
		&i.ConfigurationID,
		&i.Name,
		&i.Description,
		&i.StartTime,
		&i.EndTime,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getJobsByConfigurationID = `-- name: GetJobsByConfigurationID :many
SELECT id, configuration_id, name, description, start_time, end_time, status, created_at, updated_at FROM jobs
WHERE configuration_id = $1
`

func (q *Queries) GetJobsByConfigurationID(ctx context.Context, configurationID pgtype.UUID) ([]Job, error) {
	rows, err := q.db.Query(ctx, getJobsByConfigurationID, configurationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Job
	for rows.Next() {
		var i Job
		if err := rows.Scan(
			&i.ID,
			&i.ConfigurationID,
			&i.Name,
			&i.Description,
			&i.StartTime,
			&i.EndTime,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getJobsToSchedule = `-- name: GetJobsToSchedule :many
WITH cte AS (
    SELECT id
    FROM jobs
    WHERE status = 'Pending'
    AND start_time BETWEEN NOW() AND NOW() + INTERVAL '5 minutes'
    FOR UPDATE NOWAIT
)
UPDATE jobs
SET status = 'Running', updated_at = NOW()
WHERE id IN (SELECT id FROM cte)
RETURNING id, configuration_id, name, description, start_time, end_time, status, created_at, updated_at
`

func (q *Queries) GetJobsToSchedule(ctx context.Context) ([]Job, error) {
	rows, err := q.db.Query(ctx, getJobsToSchedule)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Job
	for rows.Next() {
		var i Job
		if err := rows.Scan(
			&i.ID,
			&i.ConfigurationID,
			&i.Name,
			&i.Description,
			&i.StartTime,
			&i.EndTime,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTaskByID = `-- name: GetTaskByID :one
SELECT id, job_id, type, status, scheduled_at, timeout, details, results, created_at, updated_at FROM tasks 
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTaskByID(ctx context.Context, id pgtype.UUID) (Task, error) {
	row := q.db.QueryRow(ctx, getTaskByID, id)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.JobID,
		&i.Type,
		&i.Status,
		&i.ScheduledAt,
		&i.Timeout,
		&i.Details,
		&i.Results,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getTasksByJobID = `-- name: GetTasksByJobID :many
SELECT id, job_id, type, status, scheduled_at, timeout, details, results, created_at, updated_at FROM tasks
WHERE job_id = $1
ORDER BY updated_at DESC
`

func (q *Queries) GetTasksByJobID(ctx context.Context, jobID pgtype.UUID) ([]Task, error) {
	rows, err := q.db.Query(ctx, getTasksByJobID, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.JobID,
			&i.Type,
			&i.Status,
			&i.ScheduledAt,
			&i.Timeout,
			&i.Details,
			&i.Results,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateJobByID = `-- name: UpdateJobByID :one
UPDATE jobs
SET 
    configuration_id = COALESCE($2, configuration_id),
    name = COALESCE($3, name),
    description = COALESCE($4, description),
    start_time = COALESCE($5, start_time),
    end_time = COALESCE($6, end_time),
    status = COALESCE($7, status)
WHERE id = $1
RETURNING id, configuration_id, name, description, start_time, end_time, status, created_at, updated_at
`

type UpdateJobByIDParams struct {
	ID              pgtype.UUID
	ConfigurationID pgtype.UUID
	Name            pgtype.Text
	Description     pgtype.Text
	StartTime       pgtype.Timestamptz
	EndTime         pgtype.Timestamptz
	Status          NullJobStatus
}

func (q *Queries) UpdateJobByID(ctx context.Context, arg UpdateJobByIDParams) (Job, error) {
	row := q.db.QueryRow(ctx, updateJobByID,
		arg.ID,
		arg.ConfigurationID,
		arg.Name,
		arg.Description,
		arg.StartTime,
		arg.EndTime,
		arg.Status,
	)
	var i Job
	err := row.Scan(
		&i.ID,
		&i.ConfigurationID,
		&i.Name,
		&i.Description,
		&i.StartTime,
		&i.EndTime,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
