INSERT INTO patients (source_system, external_id, full_name, birth_date, sex, phone, email, created_at, updated_at)
VALUES (@source_system, @external_id, @full_name, @birth_date, @sex, @phone, @email, @now, @now)
ON CONFLICT (source_system, external_id) DO UPDATE
SET full_name  = EXCLUDED.full_name,
    birth_date = EXCLUDED.birth_date,
    sex        = EXCLUDED.sex,
    phone      = EXCLUDED.phone,
    email      = EXCLUDED.email,
    updated_at = EXCLUDED.updated_at
RETURNING id;
