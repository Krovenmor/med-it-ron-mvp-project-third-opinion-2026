SELECT count(*)
FROM jobs j
JOIN cases c ON c.id = j.case_id
WHERE c.study_id = @study_id;
