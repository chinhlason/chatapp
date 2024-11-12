-- +goose Up
CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL,
                                     username varchar(20),
    password varchar(20),
    is_online BOOLEAN DEFAULT FALSE,
    PRIMARY KEY(id)
    );

INSERT INTO users (username, password, is_online) VALUES ('user1', 'user1', false);
INSERT INTO users (username, password, is_online) VALUES ('user2', 'user2', false);
INSERT INTO users (username, password, is_online) VALUES ('user3', 'user3', false);
INSERT INTO users (username, password, is_online) VALUES ('user4', 'user4', false);
INSERT INTO users (username, password, is_online) VALUES ('user5', 'user5', false);
INSERT INTO users (username, password, is_online) VALUES ('user6', 'user6', false);


-- insert 20 friends user
INSERT INTO users (username, password, is_online) VALUES ('friend1', 'friend1', false);
INSERT INTO users (username, password, is_online) VALUES ('friend2', 'friend2', false);
INSERT INTO users (username, password, is_online) VALUES ('friend3', 'friend3', false);
INSERT INTO users (username, password, is_online) VALUES ('friend4', 'friend4', false);
INSERT INTO users (username, password, is_online) VALUES ('friend5', 'friend5', false);
INSERT INTO users (username, password, is_online) VALUES ('friend6', 'friend6', false);
INSERT INTO users (username, password, is_online) VALUES ('friend7', 'friend7', false);
INSERT INTO users (username, password, is_online) VALUES ('friend8', 'friend8', false);
INSERT INTO users (username, password, is_online) VALUES ('friend9', 'friend9', false);
INSERT INTO users (username, password, is_online) VALUES ('friend10', 'friend10', false);
INSERT INTO users (username, password, is_online) VALUES ('friend11', 'friend11', false);
INSERT INTO users (username, password, is_online) VALUES ('friend12', 'friend12', false);
INSERT INTO users (username, password, is_online) VALUES ('friend13', 'friend13', false);
INSERT INTO users (username, password, is_online) VALUES ('friend14', 'friend14', false);
INSERT INTO users (username, password, is_online) VALUES ('friend15', 'friend15', false);
INSERT INTO users (username, password, is_online) VALUES ('friend16', 'friend16', false);
INSERT INTO users (username, password, is_online) VALUES ('friend17', 'friend17', false);
INSERT INTO users (username, password, is_online) VALUES ('friend18', 'friend18', false);
INSERT INTO users (username, password, is_online) VALUES ('friend19', 'friend19', false);
INSERT INTO users (username, password, is_online) VALUES ('friend20', 'friend20', false);

CREATE TABLE IF NOT EXISTS friends (
                                       id SERIAL,
                                       id_user INT REFERENCES users(id),
    id_friend INT REFERENCES users(id),
    status varchar(20),
    interaction_at timestamp
    );

INSERT INTO friends (id_user, id_friend, status) VALUES (2, 1, 'ACCEPTED');
INSERT INTO friends (id_user, id_friend, status) VALUES (4, 1, 'ACCEPTED');
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
INSERT INTO rooms (name) VALUES ('room2');

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
INSERT INTO user_in_room (id_user, id_room, id_role) VALUES (1, 2, 2);
INSERT INTO user_in_room (id_user, id_room, id_role) VALUES (4, 2, 2);

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

-- insert 20 messages
INSERT INTO messages (create_at, id_sender, id_receiver, content) VALUES ('2021-01-01 00:00:00', 1, 1, 'first message');
INSERT INTO messages (create_at, id_sender, id_receiver, content) VALUES ('2021-01-01 00:00:00', 1, 1, 'first message2');
INSERT INTO messages (create_at, id_sender, id_receiver, content) VALUES ('2021-01-01 00:00:00', 1, 1, 'first message3');
INSERT INTO messages (create_at, id_sender, id_receiver, content) VALUES ('2021-01-01 00:00:00', 1, 1, 'first message4');
INSERT INTO messages (create_at, id_sender, id_receiver, content) VALUES ('2021-01-01 00:00:00', 1, 1, 'first message5');
INSERT INTO messages (create_at, id_sender, id_receiver, content) VALUES ('2021-01-01 00:00:00', 1, 1, 'first message6');
INSERT INTO messages (create_at, id_sender, id_receiver, content) VALUES ('2021-01-01 00:00:00', 1, 1, 'first message7');
INSERT INTO messages (create_at, id_sender, id_receiver, content) VALUES ('2021-01-01 00:00:00', 1, 1, 'first message8');
INSERT INTO messages (create_at, id_sender, id_receiver, content) VALUES ('2021-01-01 00:00:00', 1, 1, 'first message9');
INSERT INTO messages (create_at, id_sender, id_receiver, content) VALUES ('2021-01-01 00:00:00', 2, 1, 'first message10');
INSERT INTO messages (create_at, id_sender, id_receiver, content) VALUES ('2021-01-01 00:00:00', 2, 1, 'first message11');
INSERT INTO messages (create_at, id_sender, id_receiver, content) VALUES ('2021-01-01 00:00:00', 2, 1, 'first message12');
INSERT INTO messages (create_at, id_sender, id_receiver, content) VALUES ('2021-01-01 00:00:00', 2, 1, 'first message13');
INSERT INTO messages (create_at, id_sender, id_receiver, content) VALUES ('2021-01-01 00:00:00', 2, 1, 'first message14');
INSERT INTO messages (create_at, id_sender, id_receiver, content) VALUES ('2021-01-01 00:00:00', 2, 1, 'first message15');
INSERT INTO messages (create_at, id_sender, id_receiver, content) VALUES ('2021-01-01 00:00:00', 2, 1, 'first message16');
INSERT INTO messages (create_at, id_sender, id_receiver, content) VALUES ('2021-01-01 00:00:00', 2, 1, 'first message17');
INSERT INTO messages (create_at, id_sender, id_receiver, content) VALUES ('2021-01-01 00:00:00', 1, 1, 'first message18');
INSERT INTO messages (create_at, id_sender, id_receiver, content) VALUES ('2021-01-01 00:00:00', 1, 1, 'first message19');
INSERT INTO messages (create_at, id_sender, id_receiver, content) VALUES ('2021-01-01 00:00:00', 1, 1, 'first message20');



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