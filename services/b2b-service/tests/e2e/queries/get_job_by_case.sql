SELECT state, attempts
FROM jobs
WHERE case_id = @case_id;
