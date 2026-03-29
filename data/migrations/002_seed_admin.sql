-- +goose Up

INSERT INTO users (login, name, password_hash, is_admin)
VALUES ('admin', 'Admin', 'PASTE_BCRYPT_HASH_HERE', 1);

-- +goose Down

DELETE FROM users WHERE login = 'admin';