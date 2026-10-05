SELECT n.id, n.case_id, n.recipient, n.channel, n.kind, n.text, n.created_at,
       r.urgency,
       COALESCE(p.external_id, '') AS patient_external_id,
       COALESCE(p.full_name, '')   AS patient_full_name
FROM notifications n
JOIN routes r ON r.case_id = n.case_id
LEFT JOIN patients p ON p.id = r.patient_id
ORDER BY n.created_at DESC
LIMIT @limit;
