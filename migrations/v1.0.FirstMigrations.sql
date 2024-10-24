create table process
(
    id          serial       not null
        constraint process_pk
            primary key,
    origin      varchar(20)  not null,
    origin_name varchar(255) not null
);

create table transactions
(
    id         integer          not null,
    amount     double precision not null,
    day        smallint         not null,
    month      smallint,
    process_id integer          not null
        constraint transactions_process_id_fk
            references process,
    constraint transactions_pk
        primary key (id, process_id)
);
