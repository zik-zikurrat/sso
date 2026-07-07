INSERT INTO user_sessions (user_id, token_hash, expires_at)
VALUES ($1, $2, $3);
