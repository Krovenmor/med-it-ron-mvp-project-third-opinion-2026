INSERT INTO bookings (case_id, recommendation_id, appointment_id, scheduled_at, channel, created_at)
VALUES (@case_id, @recommendation_id, @appointment_id, @scheduled_at, @channel, @created_at)
ON CONFLICT (appointment_id) DO NOTHING
RETURNING id, case_id, recommendation_id, appointment_id, scheduled_at, channel, created_at;
