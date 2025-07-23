create table user (
    id integer primary key,
    name text not null,
    registration integer not null default (unixepoch('now'))
);

create table public_key (
    id integer primary key,
    user_id integer not null,
    key_pem text not null,
    foreign key (user_id) references user (id)
);
