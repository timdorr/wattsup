-- migrate:up transaction:false
CREATE MATERIALIZED VIEW metrics_hourly
WITH (timescaledb.continuous) AS
SELECT 
  register_address,
  time_bucket(INTERVAL '1h', time) AS bucket,
  AVG(value) AS value
FROM metrics
GROUP BY register_address, bucket;

-- migrate:down transaction:false
DROP MATERIALIZED VIEW IF EXISTS metrics_hourly;
