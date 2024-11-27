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

