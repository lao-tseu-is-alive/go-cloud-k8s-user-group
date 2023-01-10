package users

const (
	groupsCount       = "SELECT COUNT(*) FROM go_group"
	groupsCreateTable = `
CREATE TABLE IF NOT EXISTS go_group
(
    id        	serial    CONSTRAINT go_group_pk   primary key,
    name			text	not null	constraint go_group_unique_name	unique
        constraint name_min_length check (length(btrim(name)) > 2),
    create_time		timestamp default now() not null,
    creator			integer	not null,
    last_modification_time	timestamp,
    last_modification_user	integer,
    is_active				boolean default true not null,
    inactivation_time		timestamp,
    inactivation_reason    	text,
    comment                	text
);
comment on table go_group is 'go_group contains the list of groups of the GO_USER microservice';
`

	groupCreate = `INSERT INTO go_group (name, create_time, creator, is_active, comment)
					VALUES ($1, CURRENT_TIMESTAMP,$2,true,$3) RETURNING id;`
	groupGet = `SELECT id, name, create_time, creator, last_modification_time, last_modification_user, is_active,
	inactivation_time, inactivation_reason, comment FROM go_group WHERE id=$1`
	groupList   = "SELECT id, name, is_active FROM go_group ORDER BY id"
	groupDelete = "DELETE FROM go_group WHERE id = $1"
	groupUpdate = `UPDATE go_group
SET name                   = $1,
    last_modification_time = CURRENT_TIMESTAMP,
    last_modification_user = $2,
    is_active              = $3,
    inactivation_time      = $4,
    inactivation_reason    = $5,
    comment                = $6
WHERE id = $7;`
)
