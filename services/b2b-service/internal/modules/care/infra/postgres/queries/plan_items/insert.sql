INSERT INTO plan_items (id, case_id, position, service_code, service_name, patient_text, mark)
VALUES (@id, @case_id, @position, @service_code, @service_name, @patient_text, @mark)
ON CONFLICT (id) DO NOTHING;
