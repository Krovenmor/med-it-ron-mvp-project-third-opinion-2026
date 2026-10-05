SELECT id, source_system, external_id, full_name, birth_date, sex, phone, email
FROM patients
WHERE source_system = @source_system
  AND external_id = @external_id;
