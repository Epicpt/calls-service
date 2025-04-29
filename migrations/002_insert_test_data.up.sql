INSERT INTO "users" (username, password_hash) VALUES
    ('testuser1', 'hashedpassword1'),
    ('testuser2', 'hashedpassword2');

INSERT INTO "calls" (client_name, phone_number, description, status, user_id) VALUES
    ('Иван Иванов', '+79031234567', 'Не работает телефон', 'открыта', 1),
    ('Мария Петрова', '8-903-123-45-67', 'Ошибка в приложении', 'открыта', 2);