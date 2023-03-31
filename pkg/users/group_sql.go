package users

const (
	groupsCount = "SELECT COUNT(*) FROM go_group"

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
