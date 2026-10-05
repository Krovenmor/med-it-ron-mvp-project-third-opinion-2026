SELECT type, actor, payload
FROM care.case_events
WHERE case_id = @case_id
ORDER BY id;
