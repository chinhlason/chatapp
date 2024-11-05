-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL,
    username varchar(20),
    password varchar(20),
    PRIMARY KEY(id)
);

INSERT INTO users (username, password) VALUES ('user1', 'user1');
INSERT INTO users (username, password) VALUES ('user2', 'user2');

CREATE TABLE rooms (
    id SERIAL,
    name varchar(20),
    PRIMARY KEY(id)
);

INSERT INTO rooms (name) VALUES ('room1');

CREATE TABLE user_in_room (
    id SERIAL,
    id_user INT REFERENCES users(id),
    id_room INT REFERENCES rooms(id),
    PRIMARY KEY(id)
);

INSERT INTO user_in_room (id_user, id_room) VALUES (1, 1);

CREATE TABLE IF NOT EXISTS messages (
    id SERIAL,
    create_at timestamp,
    id_sender INT REFERENCES users(id),
    id_receiver INT REFERENCES rooms(id),
    content TEXT
);

INSERT INTO messages (create_at, id_sender, id_receiver, content) VALUES ('2021-01-01 00:00:00', 1, 1, 'first message');

-- +goose Down
DROP TABLE messages;
DROP TABLE user_in_room;
DROP TABLE users;
DROP TABLE rooms;