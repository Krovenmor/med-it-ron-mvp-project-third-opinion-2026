INSERT INTO recommendations (case_id, position, source, service_code, service_name, importance, rationale,
                             guideline_ref, patient_text, already_booked, mark, marked_by, marked_at, created_at)
SELECT @case_id, COALESCE(MAX(position), 0) + 1, @source, @service_code, @service_name, @importance, @rationale,
       @guideline_ref, @patient_text, @already_booked, @mark, @marked_by, @marked_at, @created_at
FROM recommendations
WHERE case_id = @case_id
RETURNING id, case_id, position, source, service_code, service_name, importance, rationale, guideline_ref,
       patient_text, already_booked, mark, reject_reason, reject_comment, marked_by, marked_at, created_at;
