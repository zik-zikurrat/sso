SELECT id, name, method, url, secure, created_at, updated_at
FROM services
WHERE ($1 IS NULL OR id = $1)
  AND ($2 IS NULL OR name = $2)
