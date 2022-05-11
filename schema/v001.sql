create table org_department
(
    id         varchar(64) not null
        primary key,
    name       varchar(64) null,
    use_status bigint null,
    attr       bigint null,
    pid        varchar(64) null,
    super_pid  varchar(64) null,
    grade      bigint null,
    created_at bigint null,
    updated_at bigint null,
    deleted_at bigint null,
    created_by varchar(64) null,
    updated_by varchar(64) null,
    deleted_by varchar(64) null,
    tenant_id  varchar(64) null
);

create table org_use_columns
(
    id            varchar(64) not null
        primary key,
    column_id     varchar(64) null,
    viewer_status int(4) null,
    tenant_id     varchar(64) null
);

create table org_user
(
    id              varchar(64) not null
        primary key,
    name            varchar(64) null,
    phone           varchar(64) null,
    email           varchar(64) null,
    self_email      varchar(64) null,
    id_card         varchar(64) null,
    address         varchar(200) null,
    use_status      int(4) null,
    tenant_id       varchar(64) null,
    position        varchar(64) null,
    avatar          text null,
    job_number      text null,
    gender          int(4) null,
    source          varchar(64) null,
    password_status int(4) null,
    created_at      bigint null,
    updated_at      bigint null,
    deleted_at      bigint null,
    created_by      varchar(64) null,
    updated_by      varchar(64) null,
    deleted_by       varchar(64) null
);

create unique index org_user_email_uindex
    on org_user (email);


create unique index org_user_email_uindex
    on org_user (email);

create table org_user_account
(
    id         varchar(100) not null
        primary key,
    account    varchar(100) null,
    user_id    varchar(64) null,
    password   varchar(100) null,
    created_at bigint null,
    updated_at bigint null,
    deleted_at bigint null,
    created_by varchar(64) null,
    updated_by varchar(64) null,
    deleted_by varchar(64) null,
    tenant_id  varchar(64) null
);


create unique index org_user_account_account_uindex
    on org_user_account (account);

create table org_user_department_relation
(
    id      varchar(64) not null
        primary key,
    user_id varchar(64) null,
    dep_id  varchar(64) null,
    attr    varchar(64) null
);

create table org_user_table_columns
(
    id           varchar(64) not null
        primary key,
    name         varchar(64) null,
    columns_name varchar(64) null,
    types        varchar(64) null,
    len          bigint null,
    point_len    bigint null,
    attr         bigint null,
    status       bigint null,
    format       varchar(64) null,
    tenant_id    varchar(64) null,
    created_at   bigint null,
    updated_at   bigint null,
    deleted_at   bigint null,
    created_by   varchar(64) null,
    updated_by   varchar(64) null,
    deleted_by   varchar(64) null
);

create index columns_name
    on org_user_table_columns (columns_name);

create table org_user_tenant_relation
(
    id        varchar(64) not null
        primary key,
    user_id   varchar(64) null,
    tenant_id varchar(64) null,
    status    bigint      null
);

create table org_oct_use_columns
(
    id            varchar(64) not null
        primary key,
    column_id     varchar(64) null,
    created_at    bigint      null,
    updated_at    bigint      null,
    viewer_status int(4)      null,
    created_by    varchar(64) null,
    updated_bt    varchar(64) null,
    tenant_id     varchar(64) null
);

create table org_oct_user_table_columns
(
    id           varchar(64) not null
        primary key,
    name         varchar(64) null,
    columns_name varchar(64) null,
    types        varchar(64) null,
    len          bigint      null,
    point_len    bigint      null,
    attr         bigint      null,
    status       bigint      null,
    created_at   bigint      null,
    updated_at   bigint      null,
    deleted_at   bigint      null,
    created_by   varchar(64) null,
    updated_by   varchar(64) null,
    deleted_by   varchar(64) null,
    tenant_id    varchar(64) null,
    format       varchar(64) null
);

create index columns_name
    on org_oct_user_table_columns (columns_name);

create table org_user_leader_relation
(
    id        varchar(64)  not null
        primary key,
    user_id   varchar(64)  null,
    leader_id varchar(64)  null,
    attr      varchar(256) null
);

--  超管数据
INSERT INTO org_user (id, name, phone, email, self_email, id_card, address,
                      use_status, position, avatar, password_status, created_at,
                      updated_at, created_by, updated_by, tenant_id)
VALUES ('1', 'SuperAdmin', '13888886666', 'admin@yunify.com', 'admin@yunify.com', null, null, 1, null,
        null, 1, null, null, null, null, '1000');

INSERT INTO org_department (id, name, use_status, attr, pid, super_pid, grade, created_at, updated_at, deleted_at,
                            created_by, updated_by, deleted_by, tenant_id)
VALUES ('1', 'QCC', 1, 1, null, '1', 1, null, null, null, null, null, null, '1');

-- 密码 654321a..
INSERT INTO org_user_account (id, account, password, created_at, updated_at, deleted_at, created_by, updated_by,
                              deleted_by,
                              tenant_id)
VALUES ('1', 'admin@yunify.com', '24d04ec3c9f0e285791035a47ba3e61a', null, null, null, null, null, null, '1000');



INSERT INTO org_user_department_relation (id, user_id, dep_id, attr)
VALUES ('1', '1', '1', '直属领导');
