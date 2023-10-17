create table if not exists Files(
    id int primary key generated always as identity,
    name text not null,
    size bigint not null,
    mod_time timestamptz not null,
    is_directory boolean not null,
    parent_id int not null,
    path text,
    data_oid oid
);

insert into Files(name, size, mod_time, is_directory, parent_id)
values ('root', 0, now(), true, 0);