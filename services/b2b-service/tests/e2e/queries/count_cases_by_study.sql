SELECT count(*)
FROM gateway.intakes
WHERE study_id = @study_id;
