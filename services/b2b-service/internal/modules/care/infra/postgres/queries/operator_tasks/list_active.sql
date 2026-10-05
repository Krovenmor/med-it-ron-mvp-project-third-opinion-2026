SELECT t.id, t.case_id, t.reason, t.status, t.assignee, t.due_at, t.next_call_at, t.attempts,
       t.created_at, t.updated_at, t.closed_at,
       r.patient_id,
       r.status                          AS route_status,
       r.urgency,
       r.modality,
       r.received_at,
       r.confirmed_at,
       COALESCE(p.external_id, '')       AS patient_external_id,
       COALESCE(p.full_name, '')         AS patient_full_name,
       COALESCE(p.phone, '')             AS patient_phone,
       (SELECT i.service_name
        FROM plan_items i
        WHERE i.case_id = r.case_id
          AND i.declined_at IS NULL
          AND NOT EXISTS (SELECT 1 FROM bookings b WHERE b.recommendation_id = i.id)
        ORDER BY (i.mark = 'critical') DESC, i.position
        LIMIT 1) AS offer_service
FROM operator_tasks t
JOIN routes r ON r.case_id = t.case_id
LEFT JOIN patients p ON p.id = r.patient_id
WHERE t.status IN ('new', 'in_progress', 'no_answer', 'callback')
  AND (@assignee = '' OR t.assignee = @assignee)
ORDER BY (r.urgency = 'emergency') DESC, (t.due_at < @now) DESC, t.due_at;
