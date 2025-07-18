-- Insert dummy events
INSERT INTO events (event_type, event_data, user_id) VALUES
('user_login', '{"ip": "192.168.1.100", "timestamp": "2024-01-15T10:30:00Z"}', 1),
('user_logout', '{"session_duration": 3600, "timestamp": "2024-01-15T11:30:00Z"}', 1),
('page_view', '{"page": "/dashboard", "referrer": "/login", "timestamp": "2024-01-15T10:31:00Z"}', 2),
('purchase', '{"product_id": 123, "amount": 99.99, "currency": "USD", "timestamp": "2024-01-15T12:00:00Z"}', 2),
('profile_update', '{"fields_changed": ["email", "first_name"], "timestamp": "2024-01-15T14:00:00Z"}', 4),
('password_reset', '{"method": "email", "timestamp": "2024-01-15T15:30:00Z"}', 3),
('user_registration', '{"signup_method": "email", "timestamp": "2024-01-15T09:00:00Z"}', 5),
('error_occurred', '{"error_code": "500", "message": "Internal server error", "timestamp": "2024-01-15T16:00:00Z"}', 1),
('api_call', '{"endpoint": "/api/v1/users", "method": "GET", "response_time": 250, "timestamp": "2024-01-15T17:00:00Z"}', 2),
('file_upload', '{"filename": "document.pdf", "size": 1024000, "timestamp": "2024-01-15T18:00:00Z"}', 4);
