SELECT r.urgency,
       r.modality,
       r.received_at,
       r.assessed_at,
       r.review_due_at,
       r.confirmed_at,
       (SELECT min(e.occurred_at)
        FROM case_events e
        WHERE e.case_id = r.case_id
          AND e.to_status = 'notified') AS notified_at,
       (SELECT min(e.occurred_at)
        FROM case_events e
        WHERE e.case_id = r.case_id
          AND e.to_status = 'completed') AS completed_at,
       EXISTS (SELECT 1 FROM plan_items i WHERE i.case_id = r.case_id) AS accepted,
       first_booking.created_at AS booked_at,
       first_booking.channel    AS booking_channel
FROM routes r
LEFT JOIN LATERAL (SELECT b.created_at, b.channel
                   FROM bookings b
                   WHERE b.case_id = r.case_id
                   ORDER BY b.created_at
                   LIMIT 1) first_booking ON true
WHERE r.received_at >= @from
  AND r.received_at < @to;
