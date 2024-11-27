CREATE TYPE message_type AS ENUM ('Unknown', 'JobStartRequest', 'JobStart', 'JobStopRequest', 'JobStop', 'ExperimentStartRequest', 'ExperimentStart', 'ExperimentStopRequest', 'ExperimentStop');
CREATE TYPE task_type AS ENUM ('Unknown', 'DeletePod', 'EvictPod', 'TerminatePod', 'PodCrashLoop', 'ScaleDeployment', 'UpdateDeployment', 'RollbackDeployment', 'DrainNode', 'TerminateNode', 'SimulateLatency', 'SimulateNetworkLoss', 'SimulateNetworkPartition', 'OomKillPod', 'CPUThrottling', 'MemoryPressure', 'SimulateAppCrash', 'SimulateDiskFailure', 'RotateSecrets');
CREATE TYPE task_status AS ENUM ('Unknown', 'Pending', 'Scheduled', 'Running', 'Completed', 'Failed', 'Cancelled', 'TimedOut', 'Retrying', 'Skipped');
CREATE TYPE job_status AS ENUM ('Unknown', 'Pending', 'Running', 'Completed');


CREATE TABLE jobs (
    id UUID PRIMARY KEY,
    configuration_id UUID REFERENCES configurations(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    start_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ,
    status job_status,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE TABLE tasks (
    id UUID PRIMARY KEY,
    job_id UUID REFERENCES jobs(id) ON DELETE CASCADE,
    type task_type NOT NULL,
    status task_status,
    scheduled_at TIMESTAMPTZ,
    timeout INT,
    details JSONB,
    results JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE TABLE messages (
    id UUID PRIMARY KEY,
    response_id UUID REFERENCES messages(id) ON DELETE SET NULL,
    type message_type,
    timestamp TIMESTAMPTZ NOT NULL,
    source VARCHAR(255) NOT NULL,
    contents JSONB
);


CREATE TABLE logs (
    id UUID PRIMARY KEY,
    job_id UUID REFERENCES jobs(id) ON DELETE CASCADE,
    timestamp TIMESTAMPTZ DEFAULT NOW(),
    log_message TEXT
);

CREATE TABLE configurations (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    key VARCHAR(255) NOT NULL,
    value JSONB,   
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);
