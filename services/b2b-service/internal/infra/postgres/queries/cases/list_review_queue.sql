SELECT c.id                AS case_id,
       c.urgency::text     AS urgency,
       c.modality,
       c.performed_at,
       c.received_at,
       p.external_id       AS patient_external_id,
       p.full_name         AS patient_full_name,
       p.birth_date        AS patient_birth_date,
       p.sex               AS patient_sex,
       count(r.id)         AS recommendations_total,
       count(r.mark)       AS recommendations_reviewed,
       (SELECT min(e.occurred_at)
        FROM case_events e
        WHERE e.case_id = c.id
          AND e.type = 'case_opened') AS opened_at
FROM cases c
JOIN patients p ON p.id = c.patient_id
LEFT JOIN recommendations r ON r.case_id = c.id
WHERE c.status = 'in_review'
GROUP BY c.id, p.id
ORDER BY c.urgency DESC NULLS LAST, c.received_at;
