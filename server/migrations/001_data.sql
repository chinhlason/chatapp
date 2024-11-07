-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL,
    username varchar(20),
    password varchar(20),
    PRIMARY KEY(id)
);

INSERT INTO users (username, password) VALUES ('user1', 'user1');
INSERT INTO users (username, password) VALUES ('user2', 'user2');
INSERT INTO users (username, password) VALUES ('user3', 'user3');
INSERT INTO users (username, password) VALUES ('user4', 'user4');
INSERT INTO users (username, password) VALUES ('user5', 'user5');
INSERT INTO users (username, password) VALUES ('user6', 'user6');


-- insert 20 friends user
INSERT INTO users (username, password) VALUES ('friend1', 'friend1');
INSERT INTO users (username, password) VALUES ('friend2', 'friend2');
INSERT INTO users (username, password) VALUES ('friend3', 'friend3');
INSERT INTO users (username, password) VALUES ('friend4', 'friend4');
INSERT INTO users (username, password) VALUES ('friend5', 'friend5');
INSERT INTO users (username, password) VALUES ('friend6', 'friend6');
INSERT INTO users (username, password) VALUES ('friend7', 'friend7');
INSERT INTO users (username, password) VALUES ('friend8', 'friend8');
INSERT INTO users (username, password) VALUES ('friend9', 'friend9');
INSERT INTO users (username, password) VALUES ('friend10', 'friend10');
INSERT INTO users (username, password) VALUES ('friend11', 'friend11');
INSERT INTO users (username, password) VALUES ('friend12', 'friend12');
INSERT INTO users (username, password) VALUES ('friend13', 'friend13');
INSERT INTO users (username, password) VALUES ('friend14', 'friend14');
INSERT INTO users (username, password) VALUES ('friend15', 'friend15');
INSERT INTO users (username, password) VALUES ('friend16', 'friend16');
INSERT INTO users (username, password) VALUES ('friend17', 'friend17');
INSERT INTO users (username, password) VALUES ('friend18', 'friend18');
INSERT INTO users (username, password) VALUES ('friend19', 'friend19');
INSERT INTO users (username, password) VALUES ('friend20', 'friend20');

CREATE TABLE IF NOT EXISTS friends (
    id SERIAL,
    id_user INT REFERENCES users(id),
    id_friend INT REFERENCES users(id),
    status varchar(20),
    interaction_at timestamp
);

INSERT INTO friends (id_user, id_friend, status) VALUES (1, 2, 'ACCEPTED');
INSERT INTO friends (id_user, id_friend, status) VALUES (1, 3, 'PENDING');
-- INSERT INTO friends (id_user, id_friend, status) VALUES (1, 4, 'PENDING');
INSERT INTO friends (id_user, id_friend, status) VALUES (1, 5, 'PENDING');
INSERT INTO friends (id_user, id_friend, status) VALUES (1, 6, 'PENDING');
INSERT INTO friends (id_user, id_friend, status, interaction_at) VALUES (1, 4, 'ACCEPTED', '2021-01-01 00:00:00');
INSERT INTO friends (id_user, id_friend, status, interaction_at) VALUES (1, 7, 'ACCEPTED', '2021-01-01 00:00:00');
INSERT INTO friends (id_user, id_friend, status, interaction_at) VALUES (1, 8, 'ACCEPTED', '2021-01-01 00:00:00');
INSERT INTO friends (id_user, id_friend, status, interaction_at) VALUES (1, 9, 'ACCEPTED', '2021-01-01 00:00:00');
INSERT INTO friends (id_user, id_friend, status, interaction_at) VALUES (1, 10, 'ACCEPTED', '2021-01-01 00:00:07');
INSERT INTO friends (id_user, id_friend, status, interaction_at) VALUES (1, 11, 'ACCEPTED', '2021-01-01 00:00:05');
INSERT INTO friends (id_user, id_friend, status, interaction_at) VALUES (1, 12, 'ACCEPTED', '2021-01-01 00:00:01');
INSERT INTO friends (id_user, id_friend, status, interaction_at) VALUES (1, 13, 'ACCEPTED', '2021-01-01 00:00:10');
INSERT INTO friends (id_user, id_friend, status, interaction_at) VALUES (1, 14, 'ACCEPTED', '2021-01-01 00:00:50');
INSERT INTO friends (id_user, id_friend, status, interaction_at) VALUES (1, 15, 'ACCEPTED', '2021-01-01 00:00:00');
INSERT INTO friends (id_user, id_friend, status, interaction_at) VALUES (1, 16, 'ACCEPTED', '2021-01-01 00:00:00');
INSERT INTO friends (id_user, id_friend, status, interaction_at) VALUES (1, 17, 'ACCEPTED', '2021-01-01 00:00:00');
INSERT INTO friends (id_user, id_friend, status, interaction_at) VALUES (1, 18, 'ACCEPTED', '2021-01-01 00:00:00');
INSERT INTO friends (id_user, id_friend, status, interaction_at) VALUES (1, 19, 'ACCEPTED', '2021-01-01 00:00:00');
INSERT INTO friends (id_user, id_friend, status, interaction_at) VALUES (1, 20, 'ACCEPTED', '2021-01-01 00:00:04');
INSERT INTO friends (id_user, id_friend, status, interaction_at) VALUES (1, 21, 'ACCEPTED', '2021-01-01 00:00:00');
INSERT INTO friends (id_user, id_friend, status, interaction_at) VALUES (1, 22, 'ACCEPTED', '2021-01-01 00:00:00');
INSERT INTO friends (id_user, id_friend, status, interaction_at) VALUES (1, 23, 'ACCEPTED', '2021-01-01 00:00:03');
INSERT INTO friends (id_user, id_friend, status, interaction_at) VALUES (1, 24, 'ACCEPTED', '2021-01-01 00:00:00');
INSERT INTO friends (id_user, id_friend, status, interaction_at) VALUES (1, 25, 'ACCEPTED', '2021-01-01 00:00:02');
INSERT INTO friends (id_user, id_friend, status, interaction_at) VALUES (1, 26, 'ACCEPTED', '2021-01-01 00:00:01');

CREATE TABLE IF NOT EXISTS rooms (
    id SERIAL,
    name varchar(20),
    PRIMARY KEY(id)
);

INSERT INTO rooms (name) VALUES ('room1');

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

INSERT INTO user_in_room (id_user, id_room, id_role) VALUES (1, 1, 1);
INSERT INTO user_in_room (id_user, id_room, id_role) VALUES (2, 1, 1);

CREATE TABLE IF NOT EXISTS messages (
    id SERIAL,
    create_at timestamp,
    id_sender INT REFERENCES users(id),
    id_receiver INT REFERENCES rooms(id),
    content TEXT
);

INSERT INTO messages (create_at, id_sender, id_receiver, content) VALUES ('2021-01-01 00:00:00', 1, 1, 'first message');

CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL,
    create_at timestamp,
    id_sender INT REFERENCES users(id),
    id_receiver INT REFERENCES users(id),
    content TEXT
);

-- +goose Down
DROP TABLE messages;
DROP TABLE user_in_room;
DROP TABLE notifications;
DROP TABLE roles;
DROP TABLE friends;
DROP TABLE users;
DROP TABLE rooms;