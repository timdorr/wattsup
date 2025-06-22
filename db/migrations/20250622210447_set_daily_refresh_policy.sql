-- migrate:up
SELECT add_continuous_aggregate_policy('metrics_daily',
  start_offset => INTERVAL '1 day',
  end_offset => NULL,
  schedule_interval => INTERVAL '5m');

-- migrate:down
