UPDATE services
SET
    name = COALESCE($2, name),
    updated_at = now()
WHERE id = $1;
