INSERT INTO inbox (message_id, topic, processed_at)
VALUES (@message_id, @topic, @processed_at)
ON CONFLICT (message_id) DO NOTHING
RETURNING message_id;
