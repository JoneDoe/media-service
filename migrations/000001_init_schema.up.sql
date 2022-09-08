create table media
(
    uuid uuid not null
        constraint media_pk
            primary key,
    data json not null
);