SELECT state, attempts
FROM gateway.jobs
WHERE key = @case_id
  AND kind = 'assess';
