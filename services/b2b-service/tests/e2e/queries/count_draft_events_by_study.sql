SELECT count(*)
FROM case_events e
JOIN cases c ON c.id = e.case_id
WHERE c.study_id = @study_id
  AND e.to_status = 'draft';
