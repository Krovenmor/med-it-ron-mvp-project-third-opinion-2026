INSERT INTO notifications (case_id, recipient, channel, kind, text, created_at)
VALUES (@case_id, @recipient, @channel, @kind, @text, @created_at);
