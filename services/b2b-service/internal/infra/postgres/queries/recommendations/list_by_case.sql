SELECT id, case_id, position, source, service_code, service_name, importance, rationale, guideline_ref,
       patient_text, already_booked, created_at
FROM recommendations
WHERE case_id = @case_id
ORDER BY position;
