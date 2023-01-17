package users

const (
	// https://dev.to/yogski/dealing-with-enum-type-in-postgresql-1j3g
	orgunitType         = "CREATE TYPE  orgunit_type AS ENUM ('Entreprise', 'Direction', 'Service', 'Office', 'Bureau', 'Unité', 'Division');"
	orgunitTypeExist    = "select exists (select 1 from pg_type where typname = 'orgunit_type');"
	orgunitTypeList     = "SELECT UNNEST(enum_range(null::orgunit_type)) AS orgunit_type;"
	orgUnitsCount       = "SELECT COUNT(*) FROM go_orgunit"
	orgUnitsCreateTable = `
create table if not exists public.go_orgunit
(
    id                     integer              not null
        constraint pk_go_orgunit
            primary key,
    type                   orgunit_type         not null,
    name                   text                 not null
        constraint go_orgunit_unique_name unique
        constraint name_min_length check (length(btrim(name)) > 2),
    parent_id              integer,
    abbreviation           text                 not null
        constraint go_orgunit_unique_abbreviation unique
        constraint abbreviation_min_length check (length(btrim(abbreviation)) > 1),
    description            text,
    order_list             text                 not null
        constraint go_orgunit_unique_order_list unique
        constraint order_list_min_length check (length(btrim(order_list)) > 1),
    phone                  text,
    email                  text,
    create_time            timestamp            not null,
    creator                integer              not null,
    last_modification_time timestamp,
    last_modification_user integer,
    is_active              boolean default true not null,
    inactivation_time      timestamp,
    inactivation_reason    text,
    comment                text,
    guid                   uuid,
    full_name_de_norm      text
);

alter table public.go_orgunit   owner to go_cloud_k8s_user_group;
comment on table go_orgunit is 'go_orgunit contains the list of organisational units of the GO_USER microservice';
`
)
