-- migrate:up
SELECT add_continuous_aggregate_policy('metrics_hourly',
  start_offset => INTERVAL '2 hour',
  end_offset => NULL,
  schedule_interval => INTERVAL '1m');

-- migrate:down
SELECT remove_continuous_aggregate_policy('metrics_hourly');
