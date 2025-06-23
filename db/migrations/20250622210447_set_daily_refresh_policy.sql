-- migrate:up
SELECT add_continuous_aggregate_policy('metrics_hourly',
  start_offset => INTERVAL '1 hour',
  end_offset => NULL,
  schedule_interval => INTERVAL '5m');

-- migrate:down
