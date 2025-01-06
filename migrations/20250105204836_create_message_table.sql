-- +goose Up
create table message (
    id serial primary key,
    chat_id int not null,
    user_id int not null,
    message text not null,
    created_at timestamp not null default now()
);

-- +goose Down
drop table message;
