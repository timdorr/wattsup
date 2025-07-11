-- migrate:up
CREATE TABLE metrics (
   "time" TIMESTAMPTZ,
   device_id INTEGER,
   register_address INTEGER,
   value INTEGER
) WITH (
  tsdb.hypertable,
  tsdb.partition_column='time',
  tsdb.segmentby='device_id', 
  tsdb.orderby='time DESC'
);

-- migrate:down
DROP TABLE metrics;
