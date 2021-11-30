create table if not exists users
(
    id                   bigint unsigned auto_increment,
    password_digest      varchar(100),
    email                varchar(320),
    confirmed            bool default false,
    primary key (id),
    unique key email_unique_idx (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

create table if not exists sessions
(
    uuid  varchar(36),
    state tinyint not null check (state in (0, 1)),
    data  text
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

