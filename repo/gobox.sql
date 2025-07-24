create table user (
    id integer primary key,
    name text not null unique,
    registered_at datetime not null default (unixepoch('now'))
);

create table public_key (
    id integer primary key,
    user_id integer not null,
    pem text not null unique,
    added_at datetime not null default (unixepoch('now')),
    foreign key (user_id) references user (id)
);
