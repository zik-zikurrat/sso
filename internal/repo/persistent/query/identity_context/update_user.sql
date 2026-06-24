UPDATE users
SET
    login = COALESCE($2, login),
    email = COALESCE($3, email),
    password_hash = COALESCE($4, password_hash),
    updated_at = now()
WHERE id = $1
RETURNING id;
