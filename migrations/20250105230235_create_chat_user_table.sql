-- +goose Up
create table chat_user (
    chat_id INT NOT NULL,
    user_id INT NOT NULL,
    PRIMARY KEY (chat_id, user_id)
);

-- +goose Down
drop table chat_user;
