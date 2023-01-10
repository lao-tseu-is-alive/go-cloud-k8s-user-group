package users

const (
	usersList = "SELECT id, name, email, username, creator, create_time, is_admin, is_locked, is_active FROM go_user ORDER BY id;"
	usersGet  = `
SELECT id, name, email, username,
       password_hash, external_id, orgunit_id, groups_id, phone, is_locked, is_admin,
       create_time, creator, last_modification_time, last_modification_user, 
       is_active, inactivation_time, inactivation_reason, comment, bad_password_count
FROM go_user WHERE id=$1;`

	usersExist   = "SELECT COUNT(*) FROM go_user WHERE id=$1"
	usernameFind = "SELECT id FROM go_user WHERE username=$1;"
	usersCount   = "SELECT COUNT(*) FROM go_user"
	usersMaxId   = "SELECT MAX(id) FROM go_user"
	usersCreate  = `INSERT INTO go_user
(name, email, username, password_hash, external_id, orgunit_id, groups_id, phone,
 is_locked, is_admin, create_time, creator, is_active, comment, bad_password_count)
VALUES ($1, $2, $3, $4, $5, $6, $7, false, $8, CURRENT_TIMESTAMP,$9,true,$10,0)
RETURNING id;`

	usersUpdate = `
UPDATE go_user
SET name                   = $1,
    email                  = $2,
    username               = $3,
	external_id   		   = $4,
    orgunit_id             = $5,
    groups_id              = $6,
    phone                  = $7,
    is_locked              = $8,
    is_admin               = $9,
    last_modification_time = CURRENT_TIMESTAMP,
    last_modification_user = $10,
    is_active              = $11,
    inactivation_time      = $12,
    inactivation_reason    = $13,
    comment                = $14, 
	bad_password_count	   = 0  -- we decide to reset counter every time an admin update users
WHERE id = $15;
`
	usersDelete      = "DELETE FROM go_user WHERE id = $1"
	usersCreateTable = `
CREATE TABLE IF NOT EXISTS go_user
(
    id        	serial    CONSTRAINT go_user_pk   primary key,
    name			text	not null	constraint go_user_unique_name	unique
        								constraint name_min_length check (length(btrim(name)) > 2),
    email			text	not null	constraint go_user_unique_email unique
										constraint email_min_length	check (length(btrim(email)) > 3),
    username		text	not null	constraint go_user_unique_username unique
										constraint username_min_length check (length(btrim(username)) > 2),
    password_hash	text	not null 	constraint password_hash_min_length check (length(btrim(password_hash)) > 30),
	external_id		int,
    orgunit_id		int,
    groups_id		int [],
    phone			text,
    is_locked		boolean   default false not null,
    is_admin		boolean   default false not null,
    create_time		timestamp default now() not null,
    creator			integer	not null,
    last_modification_time	timestamp,
    last_modification_user	integer,
    is_active				boolean default true not null,
    inactivation_time		timestamp,
    inactivation_reason    	text,
    comment                	text,
    bad_password_count     	integer default 0 not null
);
comment on table go_user is 'go_user is the main table of the GO_USER microservice';
`
	insertAdminUser = `
INSERT INTO go_user (name, email, username, password_hash, is_admin, creator, comment) 
VALUES ('Administrative Account','admin@example.com',$1,$2, true, 1, 'Initial setup of Admin account')  RETURNING id;`

	updateAdminUser = `
UPDATE go_user
SET username               = $1,
    password_hash		 = $2,
    is_locked              = false,
    is_admin               = true,
    last_modification_time = CURRENT_TIMESTAMP,
    last_modification_user = 1,
    is_active              = true, 
	bad_password_count	   = 0  	-- we decide to reset counter 
WHERE id = 1;
`
)
