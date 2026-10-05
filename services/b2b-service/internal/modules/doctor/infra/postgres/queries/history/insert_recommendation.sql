INSERT INTO recommendations (id, case_id, position, source, service_code, service_name, importance, rationale,
                             guideline_ref, patient_text, already_booked, mark, reject_reason, reject_comment,
                             marked_by, marked_at, created_at)
VALUES (@id, @case_id, @position, @source, @service_code, @service_name, @importance, @rationale,
        @guideline_ref, @patient_text, @already_booked, NULLIF(@mark, ''), NULLIF(@reject_reason, ''),
        NULLIF(@reject_comment, ''), NULLIF(@marked_by, ''), @marked_at, @created_at);
