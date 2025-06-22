-- migrate:up transaction:false
CREATE MATERIALIZED VIEW metrics_daily
WITH (timescaledb.continuous) AS
SELECT 
  register_address,
  time_bucket('1h', time) AS bucket,
  AVG(value) AS value
FROM metrics
WHERE time > date_trunc('day', now())
GROUP BY register_address, bucket;

-- migrate:down transaction:false
DROP MATERIALIZED VIEW IF EXISTS metrics_daily;
