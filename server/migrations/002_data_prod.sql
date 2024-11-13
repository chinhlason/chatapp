-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL,
    username varchar(20),
    password varchar(20),
    is_online BOOLEAN DEFAULT FALSE,
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS friends (
    id SERIAL,
    id_user INT REFERENCES users(id),
    id_friend INT REFERENCES users(id),
    status varchar(20),
    interaction_at timestamp
);

CREATE TABLE IF NOT EXISTS rooms (
    id SERIAL,
    name varchar(20),
    interaction_at timestamp,
    PRIMARY KEY(id)
);


CREATE TABLE IF NOT EXISTS roles (
    id SERIAL,
    name varchar(20),
    PRIMARY KEY(id)
);

INSERT INTO roles (name) VALUES ('OWNER');
INSERT INTO roles (name) VALUES ('MEMBER');

CREATE TABLE IF NOT EXISTS user_in_room (
    id SERIAL,
    id_user INT REFERENCES users(id),
    id_room INT REFERENCES rooms(id),
    id_role INT REFERENCES roles(id),
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS messages (
    id SERIAL,
    create_at timestamp,
    id_sender INT REFERENCES users(id),
    id_receiver INT REFERENCES rooms(id),
    content TEXT,
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS message_read_status (
    id_message INT REFERENCES messages(id),
    id_receiver INT REFERENCES users(id),
    is_read BOOLEAN,
    read_at timestamp,
    PRIMARY KEY(id_message, id_receiver)
);

CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL,
    create_at timestamp,
    id_sender INT REFERENCES users(id),
    id_receiver INT REFERENCES users(id),
    content TEXT
);

-- +goose Down
DROP TABLE message_read_status;
DROP TABLE messages;
DROP TABLE user_in_room;
DROP TABLE notifications;
DROP TABLE roles;
DROP TABLE friends;
DROP TABLE users;
DROP TABLE rooms;