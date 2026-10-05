INSERT INTO patients (id, source_system, external_id, full_name, birth_date, sex, phone, email, updated_at)
VALUES (@id, @source_system, @external_id, @full_name, @birth_date, @sex, @phone, @email, @updated_at)
ON CONFLICT (id) DO UPDATE
SET full_name  = EXCLUDED.full_name,
    birth_date = EXCLUDED.birth_date,
    sex        = EXCLUDED.sex,
    phone      = EXCLUDED.phone,
    email      = EXCLUDED.email,
    updated_at = EXCLUDED.updated_at;
