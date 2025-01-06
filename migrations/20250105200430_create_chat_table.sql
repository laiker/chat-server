-- +goose Up
create table chat (
   id serial primary key,
   created_at timestamp not null default now()
);

-- +goose Down
drop table chat;
