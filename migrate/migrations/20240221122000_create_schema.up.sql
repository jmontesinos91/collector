SET statement_timeout = 0;

--bun:split

create table if not exists traffic
(
    id         uuid primary key unique,
    request    varchar(256)             not null,
    imei       varchar(256)             not null,
    ip         varchar(256)             not null,
    alarm      varchar(256)             not null,
    created_at timestamp with time zone not null default current_timestamp,
    updated_at timestamp with time zone not null default current_timestamp
);
