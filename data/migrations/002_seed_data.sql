-- +goose Up

INSERT INTO keys (value) VALUES
('KEY_001'),
('KEY_002');

INSERT INTO terminals (serial_number, address, name) VALUES
('TERM001', 'Moscow, Lenina 1', 'Terminal 1'),
('TERM002', 'Moscow, Tverskaya 10', 'Terminal 2');

INSERT INTO users (login, name, password_hash, is_admin) VALUES
('admin', 'Administrator', '$2a$10$hmcNT5tK9HQVlOAg5FNe2erXIkAHnK1iAYEzDafSbOov9rQ0PtLwO', 1),
('user1', 'Test User', '$2a$10$Bbs9F2P.1lynI0AuUPkuL.da7JofFoGilRk5BgIEknV4./2RyQTsG', 0);

INSERT INTO cards (card_number, balance, is_blocked, owner_name, key_id) VALUES
('CARD001', 500.0, 0, 'Ivan Ivanov', 1),
('CARD002', 150.0, 0, 'Petr Petrov', 1),
('CARD003', 0.0, 1, 'Blocked User', 2);

INSERT INTO transactions (amount, card_id, terminal_id) VALUES
(50.0, 1, 1),
(25.0, 2, 2);

-- +goose Down

DELETE FROM transactions;
DELETE FROM cards;
DELETE FROM users;
DELETE FROM terminals;
DELETE FROM keys;