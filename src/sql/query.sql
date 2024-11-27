-- name: GetTaskByID :one
SELECT * FROM tasks 
WHERE id = $1 LIMIT 1;

-- name: GetTasksByJobID :many
SELECT * FROM tasks
WHERE job_id = $1
ORDER BY updated_at DESC;

-- name: GetJobByID :one
SELECT * FROM jobs
WHERE id = $1 LIMIT 1;

-- name: GetJobsByConfigurationID :many
SELECT * FROM jobs
WHERE configuration_id = $1;

-- name: GetConfigurationByID :one
SELECT * FROM configurations
WHERE id = $1 LIMIT 1;

-- name: GetConfigurations :many
SELECT * FROM configurations
ORDER BY updated_at DESC;

-- name: AddConfiguration :one
INSERT INTO configurations (name, options)
VALUES ($1, $2)
RETURNING id;

-- name: UpdateConfigurationsByID :one
UPDATE configurations
SET 
    name = COALESCE(sqlc.narg(name), name),
    options = COALESCE(sqlc.narg(options), options)
WHERE id = $1
RETURNING *;

-- name: GetJobsToSchedule :many
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
RETURNING *;

-- name: AddJob :one
INSERT INTO jobs (configuration_id, name, description, start_time, end_time, status)
Values ($1, $2, $3, $4, $5, $6)
returning id;

-- name: UpdateJobByID :one
UPDATE jobs
SET 
    configuration_id = COALESCE(sqlc.narg(configuration_id), configuration_id),
    name = COALESCE(sqlc.narg(name), name),
    description = COALESCE(sqlc.narg(description), description),
    start_time = COALESCE(sqlc.narg(start_time), start_time),
    end_time = COALESCE(sqlc.narg(end_time), end_time),
    status = COALESCE(sqlc.narg(status), status)
WHERE id = $1
RETURNING *;