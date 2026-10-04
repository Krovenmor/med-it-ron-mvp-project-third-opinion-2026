SELECT DISTINCT r.service_code, r.service_name
FROM recommendations r
JOIN cases c ON c.id = r.case_id
WHERE c.patient_id = @patient_id
  AND c.id <> @exclude_case_id
  AND c.status IN ('confirmed', 'notified', 'booked');
