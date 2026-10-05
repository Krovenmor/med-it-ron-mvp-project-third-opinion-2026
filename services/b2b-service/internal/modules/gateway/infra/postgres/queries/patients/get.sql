SELECT id, source_system, external_id, full_name, birth_date, sex, phone, email
FROM patients
WHERE id = @id;
