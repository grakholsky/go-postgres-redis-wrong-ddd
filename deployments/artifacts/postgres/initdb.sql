create table users
(
    id    serial       not null
        constraint users_pkey
            primary key,
    name  varchar(255) not null,
    email varchar(255) not null
);

create unique index users_email_uindex
    on users (email);

create table devices
(
    id      serial       not null
        constraint devices_pkey
            primary key,
    name    varchar(255) not null,
    user_id integer      not null
        constraint devices_user_id_fk
            references users
            on delete cascade
);

create table device_metrics
(
    id          serial  not null
        constraint device_metrics_pkey
            primary key,
    device_id   integer not null
        constraint device_metrics_device_id_fk
            references devices
            on delete cascade,
    metric_1    integer,
    metric_2    integer,
    metric_3    integer,
    metric_4    integer,
    metric_5    integer,
    local_time  timestamp,
    server_time timestamp default now()
);

create table device_alerts
(
    id        serial  not null
        constraint device_alerts_pkey
            primary key,
    device_id integer not null
        constraint device_alerts_devices_id_fk
            references devices
            on delete cascade,
    message   text
);