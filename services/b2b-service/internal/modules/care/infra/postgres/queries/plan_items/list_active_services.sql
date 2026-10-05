SELECT DISTINCT i.service_code, i.service_name
FROM plan_items i
JOIN routes r ON r.case_id = i.case_id
WHERE r.patient_id = @patient_id
  AND r.case_id <> @exclude_case_id
  AND r.status IN ('confirmed', 'notified', 'booked');
