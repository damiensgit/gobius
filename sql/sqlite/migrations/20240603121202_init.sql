-- +goose Up
-- +goose StatementBegin
CREATE TABLE commitments (
    taskid TEXT PRIMARY KEY,
    commitment TEXT NOT NULL,
    validator TEXT NOT NULL,
    added TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE solutions (
    taskid TEXT PRIMARY KEY,
    cid BLOB NOT NULL,
    validator TEXT NOT NULL,
    added TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ipfs_cids (
    taskid TEXT PRIMARY KEY,
    cid BLOB NOT NULL,
    added TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tasks (
    taskid TEXT PRIMARY KEY,
    txhash TEXT NOT NULL,
    cumulativeGas FLOAT NOT NULL DEFAULT 0,
    status int not null default 0,
    claimtime int NOT NULL DEFAULT 0
);   

-- Index on status for tasks table as it's used as a filter
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_taskid_and_status ON tasks(taskid, status);

-- Index on claimtime for tasks table as it's used as a filter
CREATE INDEX idx_tasks_claimtime ON tasks(claimtime);

CREATE INDEX idx_solutions_date ON solutions(added);
CREATE INDEX idx_commitments_date ON commitments(added);

CREATE INDEX idx_solutions_validator ON solutions(validator);
CREATE INDEX idx_commitments_validator ON commitments(validator);

CREATE INDEX idx_ipfs_cids_date ON ipfs_cids(added);

-- +goose StatementEnd