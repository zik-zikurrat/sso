SELECT s.id, s.name, s.created_at, s.updated_at,
       e.id, e.method, e.url, e.secure, e.created_at
FROM services AS s
    JOIN endpoints AS e ON s.id = e.service_id;
