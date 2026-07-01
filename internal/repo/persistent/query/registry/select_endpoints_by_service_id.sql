SELECT id, method, url, secure, created_at
FROM endpoints
WHERE service_id = $1;
