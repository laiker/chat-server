-- +goose Up
create table chat_log (
   id serial primary key,
   name text not null,
   entity_id int null,
   created_at timestamp not null default now()
);

-- +goose Down
drop table chat_log;
