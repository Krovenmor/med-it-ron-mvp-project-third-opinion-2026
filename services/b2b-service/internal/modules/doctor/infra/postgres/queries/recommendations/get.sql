SELECT id, case_id, position, source, service_code, service_name, importance, rationale, guideline_ref,
       patient_text, already_booked, mark, reject_reason, reject_comment, marked_by, marked_at, created_at
FROM recommendations
WHERE id = @id
  AND case_id = @case_id;
