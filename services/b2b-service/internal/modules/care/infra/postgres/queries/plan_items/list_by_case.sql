SELECT id, case_id, position, service_code, service_name, patient_text, mark, declined_at, decline_reason,
       decline_comment
FROM plan_items
WHERE case_id = @case_id
ORDER BY position;
