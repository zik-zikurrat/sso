SELECT id, login, email, password_hash, role, created_at, updated_at
FROM users
WHERE ($1 IS NULL OR id = $1)
  AND ($2 IS NULL OR login = $2)
  AND ($3 IS NULL OR email = $3);
