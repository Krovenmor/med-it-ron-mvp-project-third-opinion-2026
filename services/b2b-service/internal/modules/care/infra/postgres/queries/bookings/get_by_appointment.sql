SELECT id, case_id, recommendation_id, appointment_id, scheduled_at, channel, created_at
FROM bookings
WHERE appointment_id = @appointment_id;
