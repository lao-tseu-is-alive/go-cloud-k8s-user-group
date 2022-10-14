package users

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/config"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/crypto"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/database"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/version"
	"log"
	"time"
)

var ErrUsernameNotFound = errors.New("username does not exist")

const (
	usersList = "SELECT id, name, email, username, creator, create_time, is_admin, is_locked, is_active FROM go_user ORDER BY id;"
	usersGet  = `
SELECT id, name, email, username,
       password_hash, external_id, orgunit_id, phone, is_locked, is_admin,
       create_time, creator, last_modification_time, last_modification_user, 
       is_active, inactivation_time, inactivation_reason, comment, bad_password_count
FROM go_user WHERE id=$1;`

	usersExist   = "SELECT COUNT(*) FROM go_user WHERE id=$1"
	usernameFind = "SELECT id FROM go_user WHERE username=$1;"
	usersCount   = "SELECT COUNT(*) FROM go_user"
	usersMaxId   = "SELECT MAX(id) FROM go_user"
	usersCreate  = `INSERT INTO go_user
(name, email, username, password_hash, external_id, orgunit_id, phone,
 is_locked, is_admin, create_time, creator, is_active, comment, bad_password_count)
VALUES ($1, $2, $3, $4, $5, $6, $7, false, $8, CURRENT_TIMESTAMP,$9,true,$10,0)
RETURNING id;`

	usersUpdate = `
UPDATE go_user
SET name                   = $1,
    email                  = $2,
    username               = $3,
    orgunit_id             = $4,
    phone                  = $5,
    is_locked              = $6,
    is_admin               = $7,
    last_modification_time = CURRENT_TIMESTAMP,
    last_modification_user = $8,
    is_active              = $9,
    inactivation_time      = $10,
    inactivation_reason    = $11,
    comment                = $12,
    external_id   		   = $13,
	bad_password_count	   = 0  -- we decide to reset counter every time an admin update users
WHERE id = $14;
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
    comment                = $6,
WHERE id = $7;`
)

type PGX struct {
	Db  *database.PgxDB
	log *log.Logger
}

func NewPgxDB(dbConnectionString string, maxConnectionsInPool int, log *log.Logger) (Storage, error) {
	var psql PGX
	var successOrFailure = "OK"

	pgxPool, err := database.GetPgxConn(dbConnectionString, maxConnectionsInPool, log)
	if err != nil {
		successOrFailure = "FAILED"
		log.Printf("Connecting to database : %s \n", successOrFailure)
		return nil, errors.New(fmt.Sprintf("error connecting to database. err : %s", err))
	}
	log.Printf("INFO: 'Connection to database : %s'\n", successOrFailure)
	var numberOfServicesSchema int
	errMetaTable := pgxPool.Conn.QueryRow(context.Background(), countMetaUserServiceSQL).Scan(&numberOfServicesSchema)
	if errMetaTable != nil {
		log.Printf("WARNING: problem counting the rows in metadata table : %v", errMetaTable)
		return nil, errors.New("unable to countMetaUserServiceSQL")
	}
	if numberOfServicesSchema < 1 {
		log.Printf("WARNING: database does not contain this service in the metadata table, will try to insert it  ! ")
		var lastInsertId int = 0
		err := pgxPool.Conn.QueryRow(context.Background(), insertMetaUserService, version.VERSION).Scan(&lastInsertId)
		if err != nil {
			log.Printf("ERROR: insertMetaUserService ver:(%s) unexpectedly failed. error : %v", version.VERSION, err)
			return nil, errors.New("unable to insert metadata for the service table «go_user» ")
		}
		log.Printf("INFO: insertMetaUserService ver:(%s) created with id : %d", version.VERSION, lastInsertId)
	}
	var numberOfUsers int
	errUsersTable := pgxPool.Conn.QueryRow(context.Background(), usersCount).Scan(&numberOfUsers)
	if errUsersTable != nil {
		log.Printf("WARNING: problem counting the rows in «go_user» table : %v", errUsersTable)
		log.Printf("WARNING: 'the database does not contain the table «go_user»  wil try to create it!'")
		commandTag, err := pgxPool.Conn.Exec(context.Background(), usersCreateTable)
		if err != nil {
			log.Printf("ERROR: problem creating the «go_user» table : %v", err)
			return nil, errors.New("unable to create the table «go_user» ")
		}
		log.Printf("SUCCESS: «go_user» table was created rows affected : %v", int(commandTag.RowsAffected()))
	}
	var numberOfGroups int
	errGroupsTable := pgxPool.Conn.QueryRow(context.Background(), groupsCount).Scan(&numberOfGroups)
	if errGroupsTable != nil {
		log.Printf("WARNING: problem counting the rows in «go_group» table : %v", errUsersTable)
		log.Printf("WARNING: 'the database does not contain the table «go_group»  wil try to create it!'")
		commandTag, err := pgxPool.Conn.Exec(context.Background(), groupsCreateTable)
		if err != nil {
			log.Printf("ERROR: problem creating the «go_group» table : %v", err)
			return nil, errors.New("unable to create the table «go_group» ")
		}
		log.Printf("SUCCESS: «go_group» table was created rows affected : %v", int(commandTag.RowsAffected()))
	}

	adminUser := config.GetAdminUserFromFromEnv("admin")
	adminPassword, err := config.GetAdminPasswordFromFromEnv()
	if err != nil {
		log.Printf("ERROR: 'GetAdminPasswordFromFromEnv returned error : %v'", err)
		return nil, errors.New("unable to retrieve a valid admin password GetAdminPasswordFromFromEnv")
	}
	var lastInsertId int = 0
	passwordHash := crypto.Sha256Hash(adminPassword)
	goHash, err := crypto.HashAndSalt(passwordHash)
	if err != nil {
		log.Printf("ERROR: crypto.HashAndSalt unexpectedly failed. error : %v", err)
		return nil, errors.New("unable to calculate hash for the admin password ")
	}

	if numberOfUsers > 0 {
		log.Printf("INFO: 'database contains %d records in «go_user»'", numberOfUsers)
		commandTag, err := pgxPool.Conn.Exec(context.Background(), updateAdminUser, adminUser, goHash)
		if err != nil {
			log.Printf("ERROR: updateAdminUser adminUser:(%s) hash : %s unexpectedly failed. error : %v", adminUser, goHash, err)
			return nil, errors.New("unable to update adminUser in table «go_user» ")
		}
		log.Printf("INFO: 'update %d row with admin user %s  in «go_user»'", int(commandTag.RowsAffected()), adminUser)
	} else {
		log.Printf("WARNING: '«go_user» contain %d records : creating initial admin user: %s'", numberOfUsers, adminUser)
		err = pgxPool.Conn.QueryRow(context.Background(), insertAdminUser, adminUser, goHash).Scan(&lastInsertId)
		if err != nil {
			log.Printf("ERROR: insertAdminUser adminUser:(%s) hash : %s unexpectedly failed. error : %v", adminUser, goHash, err)
			return nil, errors.New("unable to insert adminUser in table «go_user» ")
		}
		log.Printf("INFO: insertAdminUser adminUser:(%s) created with id : %d", adminUser, lastInsertId)
	}
	psql.Db = pgxPool
	psql.log = log
	return &psql, err
}

// Create will store the new user in the store
func (db *PGX) Create(u User) (*User, error) {
	db.log.Printf("trace : entering Create(%q,%q)", u.Name, u.Username)
	if len(u.Name) < 1 {
		return nil, errors.New("user name cannot be empty")
	}
	if len(u.Name) < 6 {
		return nil, errors.New("CreateUser name minLength is 5")
	}
	var lastInsertId int = 0

	goHash, err := crypto.HashAndSalt(u.PasswordHash)
	if err != nil {
		db.log.Printf("error : Create(%q) had an error doing crypto.HashAndSalt. error : %v", u.Name, err)
		return nil, err
	}
	u.PasswordHash = goHash
	err = db.Db.Conn.QueryRow(context.Background(), usersCreate,
		u.Name, u.Email, u.Username, u.PasswordHash, &u.ExternalId, &u.OrgunitId, &u.Phone, //$1-$7
		u.IsAdmin, u.Creator, &u.Comment).Scan(&lastInsertId)
	if err != nil {
		db.log.Printf("error : Create(%q) unexpectedly failed. error : %v", u.Name, err)
		return nil, err
	}
	db.log.Printf("info : Create(%q) created with id : %v", u.Name, lastInsertId)

	// if we get to here all is good, so let's retrieve a fresh copy to send it back
	createdUser, err := db.Get(int32(lastInsertId))
	if err != nil {
		return nil, GetErrorF("error : users was created, but can not be retrieved", err)
	}
	return createdUser, nil
}

// List will retrieve all users in the store
func (db *PGX) List(offset, limit int) ([]*UserList, error) {
	db.log.Println("trace : entering List()")
	var res []*UserList

	err := pgxscan.Select(context.Background(), db.Db.Conn, &res, usersList)
	if err != nil {
		db.log.Printf("error : List pgxscan.Select unexpectedly failed, error : %v", err)
		return nil, err
	}
	if res == nil {
		db.log.Println("info : List returned no results ")
		return nil, errors.New("records not found")
	}

	return res, nil
}

func (db *PGX) Get(id int32) (*User, error) {
	db.log.Printf("trace : entering Get(%d)", id)
	res := &User{}
	err := pgxscan.Get(context.Background(), db.Db.Conn, res, usersGet, id)
	if err != nil {
		db.log.Printf("error : Get(%d) pgxscan.Select unexpectedly failed, error : %v", id, err)
		return nil, err
	}
	if res == nil {
		db.log.Printf("info : Get(%d) returned no results ", id)
		return nil, errors.New("records not found")
	}
	return res, nil
}

// GetMaxId returns the maximum value of users id existing in store.
func (db *PGX) GetMaxId() (int32, error) {
	db.log.Println("trace : entering GetMaxId()")
	existingMaxId, err := db.Db.GetQueryInt(usersMaxId)
	if err != nil {
		db.log.Printf("getMaxId() could not be retrieved from DB. failed db.Query err: %v", err)
		return 0, err
	}
	return int32(existingMaxId), nil
}

// Exist returns true only if a users with the specified id exists in store.
func (db *PGX) Exist(id int32) bool {
	db.log.Printf("trace : entering Exist(%d)", id)
	count, err := db.Db.GetQueryInt(usersExist, id)
	if err != nil {
		db.log.Printf("error: Exist(%d) could not be retrieved from DB. failed db.Query err: %v", id, err)
		return false
	}
	if count > 0 {
		db.log.Printf("info : Exist(%d) id does exist  count:%v", id, count)
		return true
	} else {
		db.log.Printf("info : Exist(%d) id does not exist count:%v", id, count)
		return false
	}
}

// Count returns the number of users stored in DB
func (db *PGX) Count() (int32, error) {
	db.log.Println("trace : entering Count()")
	count, err := db.Db.GetQueryInt(usersCount)
	if err != nil {
		db.log.Printf("error: Count() could not be retrieved from DB. failed db.Query err: %v", err)
		return 0, err
	}
	return int32(count), nil
}

// Update the users stored in DB with given id and other information in struct
func (db *PGX) Update(id int32, user User) (*User, error) {
	db.log.Printf("trace : entering Update(%d)", id)
	// first check business rules for name field
	if len(user.Name) < 1 {
		return nil, errors.New("user name cannot be empty")
	}
	if len(user.Name) < 6 {
		return nil, errors.New("CreateUser name minLength is 5")
	}
	var rowsAffected int = 0
	var err error
	now := time.Now()
	user.LastModificationTime = &now
	if user.IsActive == false {
		user.InactivationTime = &now
	} else {
		user.InactivationTime = nil
	}
	if len(user.PasswordHash) > 0 {
		err := db.ResetPassword(id, user.PasswordHash, int(*user.LastModificationUser))
		if err != nil {
			return nil, GetErrorF("error: Update() reset password failed", err)
		}
		db.log.Printf("info : password change successful for user id [%d]", id)
	}
	db.log.Printf("info : just before Update(%+v)", user)
	rowsAffected, err = db.Db.ExecActionQuery(usersUpdate,
		user.Name, user.Email, user.Username,
		&user.OrgunitId, &user.Phone, user.IsLocked, user.IsAdmin, user.LastModificationUser,
		user.IsActive, &user.InactivationTime, &user.InactivationReason, &user.Comment, user.ExternalId, id)
	if err != nil {
		return nil, GetErrorF("error: Update() query failed", err)
	}
	if rowsAffected < 1 {
		return nil, GetErrorF("error : Update() no row modified", err)
	}
	// if we get to here all is good, so let's retrieve a fresh copy to send it back
	updatedUser, err := db.Get(id)
	if err != nil {
		return nil, GetErrorF("error : Update() user updated, but can not be retrieved", err)
	}
	return updatedUser, nil
}

// Delete the users stored in DB with given id
func (db *PGX) Delete(id int32) error {
	db.log.Printf("trace : entering Delete(%d)", id)
	rowsAffected, err := db.Db.ExecActionQuery(usersDelete, id)
	if err != nil {
		return GetErrorF("error : user could not be deleted", err)
	}
	if rowsAffected < 1 {
		return GetErrorF("error : user was not deleted", err)
	}
	// if we get to here all is good
	return nil
}

// FindUsername retrieves the user id for the given username or err if not found.
func (db *PGX) FindUsername(username string) (int32, error) {
	db.log.Printf("trace : entering FindUsername(%s)", username)
	idUser, err := db.Db.GetQueryInt(usernameFind, username)
	if err != nil {
		if err == pgx.ErrNoRows {
			db.log.Printf("warning: FindUsername(%s) did not find any rows with this username", username)
			return 0, ErrUsernameNotFound
		}
		db.log.Printf("error: FindUsername(%s) could not be retrieved from DB. failed db.Query err: %v", username, err)
		return 0, err
	}
	db.log.Printf("info : FindUsername(%s) id was found :%v", username, idUser)
	return int32(idUser), nil
}
func (db *PGX) Close() {
	db.Db.Close()
}

// IsUserAdmin returns true if the user with the specified id has the is_admin attribute set to true
func (db *PGX) IsUserAdmin(idUser int32) bool {
	var isAdmin bool
	isAdmin, err := db.Db.GetQueryBool("SELECT is_admin FROM go_user WHERE id = $1", idUser)
	if err != nil {
		db.log.Printf("error: IsUserAdmin(%d) could not be retrieved from DB. failed db.Query err: %v", idUser, err)
		return false
	}
	return isAdmin
}

// IsUserActive returns true if the user with the specified id has the is_active attribute set to true
func (db *PGX) IsUserActive(idUser int32) bool {
	var isActive bool
	isActive, err := db.Db.GetQueryBool("SELECT is_active FROM go_user WHERE id = $1", idUser)
	if err != nil {
		db.log.Printf("error: IsUserActive(%d) could not be retrieved from DB. failed db.Query err: %v", idUser, err)
		return false
	}
	return isActive
}

// ResetPassword the user password stored in DB for given id
func (db *PGX) ResetPassword(id int32, passwordHash string, idCurrentUser int) error {
	db.log.Printf("trace : entering ResetPassword(%d)", id)
	goHash, err := crypto.HashAndSalt(passwordHash)
	if err != nil {
		db.log.Printf("error : ResetPassword(%d) had an error doing crypto.HashAndSalt. error : %v", id, err)
		return err
	}
	rowsAffected, err := db.Db.ExecActionQuery(
		"UPDATE go_user SET is_locked=false, password_hash = $1, last_modification_time = now(), last_modification_user = $2  WHERE id = $3;",
		goHash, idCurrentUser, id)
	if err != nil {
		return GetErrorF("error : could not reset password for user ", err)
	}
	if rowsAffected < 1 {
		return GetErrorF("error : password fo user was not reset", err)
	}
	// if we get to here all is good
	return nil
}

// CreateGroup saves a new group in the storage.
func (db *PGX) CreateGroup(g Group) (*Group, error) {
	db.log.Printf("trace : entering CreateGroup(%q)", g.Name)
	if len(g.Name) < 1 {
		return nil, errors.New("group.name cannot be empty")
	}
	if len(g.Name) < 6 {
		return nil, errors.New("group.name minLength is 5")
	}
	var lastInsertId int = 0

	err := db.Db.Conn.QueryRow(context.Background(), groupCreate,
		g.Name, g.Creator, &g.Comment).Scan(&lastInsertId)
	if err != nil {
		db.log.Printf("error : CreateGroup(%q) unexpectedly failed. error : %v", g.Name, err)
		return nil, err
	}
	db.log.Printf("info : CreateGroup(%q) created with id : %v", g.Name, lastInsertId)

	// if we get to here all is good, so let's retrieve a fresh copy to send it back
	created, err := db.GetGroup(int32(lastInsertId))
	if err != nil {
		return nil, GetErrorF("error : group was created, but can not be retrieved", err)
	}
	return created, nil
}

// GetGroup returns the group with the specified users ID.
func (db *PGX) GetGroup(id int32) (*Group, error) {
	db.log.Printf("trace : entering GetGroup(%d)", id)
	res := &Group{}
	err := pgxscan.Get(context.Background(), db.Db.Conn, res, groupGet, id)
	if err != nil {
		db.log.Printf("error : GetGroup(%d) Select unexpectedly failed, error : %v", id, err)
		return nil, err
	}
	if res == nil {
		db.log.Printf("info : GetGroup(%d) returned no results ", id)
		return nil, errors.New("records not found")
	}
	return res, nil
}

// ListGroup will retrieve all groups in the store
func (db *PGX) ListGroup(offset, limit int) ([]*GroupList, error) {
	db.log.Println("trace : entering ListGroup()")
	var res []*GroupList

	err := pgxscan.Select(context.Background(), db.Db.Conn, &res, groupList)
	if err != nil {
		db.log.Printf("error : ListGroup Select unexpectedly failed, error : %v", err)
		return nil, err
	}
	if res == nil {
		db.log.Println("info : ListGroup query returned no results ")
		return nil, errors.New("records not found")
	}

	return res, nil
}

// DeleteGroup removes the group with given ID from the storage.
func (db *PGX) DeleteGroup(id int32) error {
	db.log.Printf("trace : entering DeleteGroup(%d)", id)
	rowsAffected, err := db.Db.ExecActionQuery(groupDelete, id)
	if err != nil {
		return GetErrorF("error : group could not be deleted", err)
	}
	if rowsAffected < 1 {
		return GetErrorF("error : group was not deleted", err)
	}
	return nil
}

// UpdateGroup updates the group with given ID in the storage.
func (db *PGX) UpdateGroup(id int32, group Group) (*Group, error) {
	db.log.Printf("trace : entering UpdateGroup(%d)", id)
	// first check business rules for name field
	if len(group.Name) < 1 {
		return nil, errors.New("user name cannot be empty")
	}
	if len(group.Name) < 6 {
		return nil, errors.New("CreateUser name minLength is 5")
	}
	var rowsAffected int = 0
	var err error
	now := time.Now()
	group.LastModificationTime = &now
	if group.IsActive == false {
		group.InactivationTime = &now
	} else {
		group.InactivationTime = nil
	}
	db.log.Printf("info : just before UpdateGroup(%+v)", group)
	rowsAffected, err = db.Db.ExecActionQuery(groupUpdate,
		group.Name, group.LastModificationUser, group.IsActive, &group.InactivationTime,
		&group.InactivationReason, &group.Comment, id)
	if err != nil {
		return nil, GetErrorF("error: UpdateGroup() query failed", err)
	}
	if rowsAffected < 1 {
		return nil, GetErrorF("error : UpdateGroup() no row modified", err)
	}
	// if we get to here all is good, so let's retrieve a fresh copy to send it back
	updatedUser, err := db.GetGroup(id)
	if err != nil {
		return nil, GetErrorF("error : Update() group updated, but can not be retrieved", err)
	}
	return updatedUser, nil
}
