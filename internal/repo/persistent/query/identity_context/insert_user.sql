INSERT INTO users (email, password_hash, login, role)
VALUES ($1, $2, $3, $4)
RETURNING id;
