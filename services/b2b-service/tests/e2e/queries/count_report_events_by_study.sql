SELECT count(*)
FROM gateway.jobs j
JOIN gateway.intakes i ON j.key = i.case_id::text
WHERE i.study_id = @study_id
  AND j.kind = 'deliver_event'
  AND j.payload ->> 'topic' = 'gateway.report_received';
