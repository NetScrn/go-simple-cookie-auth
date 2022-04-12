create table if not exists users
(
    id              bigint unsigned auto_increment,
    password_digest varchar(100),
    email           varchar(320),
    confirmed       bool default false,
    primary key (id),
    unique key email_unique_idx (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

create table if not exists sessions
(
    uuid       varchar(36),
    user_id    bigint unsigned not null,
    active     boolean         default false,
    expiry     timestamp,
    created_at timestamp       not null default CURRENT_TIMESTAMP,
    data       text,
    primary key (uuid),
    foreign key (user_id) references users(id) on delete restrict
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

