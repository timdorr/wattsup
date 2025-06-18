-- name: CreateMetric :exec
INSERT INTO metrics (
  time, device_name, device_id, register_name, register_address, value
) VALUES (
  $1, $2, $3, $4, $5, $6
);
