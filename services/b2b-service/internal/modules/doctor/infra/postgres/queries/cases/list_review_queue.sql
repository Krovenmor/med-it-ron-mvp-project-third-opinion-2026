SELECT c.id                  AS case_id,
       c.urgency::text       AS urgency,
       c.modality,
       c.performed_at,
       c.received_at,
       c.patient_id,
       c.patient_external_id,
       c.patient_birth_date,
       c.patient_sex,
       count(r.id)           AS recommendations_total,
       count(r.mark)         AS recommendations_reviewed,
       (SELECT min(e.occurred_at)
        FROM case_events e
        WHERE e.case_id = c.id
          AND e.type = 'case_opened') AS opened_at
FROM cases c
LEFT JOIN recommendations r ON r.case_id = c.id
WHERE c.status = 'in_review'
GROUP BY c.id
ORDER BY c.urgency DESC, c.received_at;
