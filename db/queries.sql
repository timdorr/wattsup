-- name: CreateMetrics :copyfrom
INSERT INTO metrics (
  "time", device_id, register_address, "value"
) VALUES (
  $1, $2, $3, $4
);
