SELECT DISTINCT recommendation_id
FROM bookings
WHERE case_id = @case_id;
