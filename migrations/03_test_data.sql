-- Тестовый пользователь
INSERT INTO users (username, email, password_hash) VALUES
('test_user', 'test@example.com', '$2a$10$X7Q8Y9Z0A1B2C3D4E5F6G7H8I9J0K1L2M3N4O5P6Q7R8S9T0U1V2W3X4Y5Z');

-- Тестовый счет
INSERT INTO accounts (user_id, balance) VALUES
(1, 10000.00);

-- Тестовая карта (зашифрованные данные)
INSERT INTO cards (account_id, card_number, expiry_date, cvv_hash, hmac) VALUES
(1, 
 pgp_sym_encrypt('1234567890123456', 'card_secret_key'),
 pgp_sym_encrypt('12/25', 'card_secret_key'),
 '$2a$10$X7Q8Y9Z0A1B2C3D4E5F6G7H8I9J0K1L2M3N4O5P6Q7R8S9T0U1V2W3X4Y5Z',
 pgp_sym_encrypt('1234567890123456', 'hmac_secret_key')
);

-- Тестовая транзакция
INSERT INTO transactions (from_account_id, to_account_id, amount, type, status) VALUES
(1, 1, 1000.00, 'DEPOSIT', 'COMPLETED');

-- Тестовый кредит
INSERT INTO credits (user_id, account_id, amount, interest_rate, term_months, status) VALUES
(1, 1, 50000.00, 15.00, 12, 'ACTIVE');

-- Тестовый график платежей
INSERT INTO payment_schedules (credit_id, payment_date, amount, principal, interest, status) VALUES
(1, CURRENT_DATE + INTERVAL '1 month', 4512.92, 4166.67, 346.25, 'PENDING'); 