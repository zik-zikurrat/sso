SELECT method, url, secure
FROM endpoints
WHERE service_id = $1;