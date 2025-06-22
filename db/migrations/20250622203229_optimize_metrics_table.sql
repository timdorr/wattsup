-- migrate:up

ALTER TABLE metrics SET (
  tsdb.segmentby = 'register_address',
  tsdb.chunk_interval='1 day'
);

CREATE INDEX metrics_register_address_time_idx ON metrics (register_address, time DESC);


-- migrate:down
ALTER TABLE metrics RESET (
  tsdb.partition_column='time',
  tsdb.segmentby='device_id', 
  tsdb.chunk_interval='7 days'
);

DROP INDEX IF EXISTS metrics_register_address_time_idx;
