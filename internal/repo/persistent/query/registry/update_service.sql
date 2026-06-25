UPDATE services
SET
    name = COALESCE($2, name),
    method = COALESCE($3, method),
    url = COALESCE($4, url),
    updated_at = now()
WHERE id = $1
RETURNING id;
