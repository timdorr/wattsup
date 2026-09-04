-- migrate:up
SELECT add_retention_policy('metrics', drop_after => INTERVAL '6 months');

-- migrate:down
SELECT remove_retention_policy('metrics');
