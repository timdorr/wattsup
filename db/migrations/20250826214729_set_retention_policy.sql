-- migrate:up
SELECT add_retention_policy('metrics', INTERVAL '6 months');

-- migrate:down
SELECT remove_retention_policy('metrics');
