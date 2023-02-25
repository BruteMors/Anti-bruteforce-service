create table blacklist (
                           id serial PRIMARY KEY not null,
                           prefix varchar not null,
                           mask varchar not null,
                           CREATED_AT timestamp DEFAULT now()
);

create table whitelist (
                           id serial PRIMARY KEY not null,
                           prefix varchar not null,
                           mask varchar not null,
                           CREATED_AT timestamp DEFAULT now()
)