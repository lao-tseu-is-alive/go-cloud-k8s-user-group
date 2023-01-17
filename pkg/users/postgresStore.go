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

type PGX struct {
	Db  *database.PgxDB
	log *log.Logger
}

// NewPgxDB will instantiate a new storage of type postgres and ensure schema exist
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
	isOrgUnitTypeAlreadyThere := false
	err = pgxPool.Conn.QueryRow(context.Background(), orgunitTypeExist).Scan(&isOrgUnitTypeAlreadyThere)
	if err != nil {
		log.Printf("error : checking orgunitTypeExist unexpectedly failed. args : (%v), error : %v\n", orgunitTypeExist, err)
		return nil, errors.New("unable to check if type «orgunit_type» already exists")
	}
	if isOrgUnitTypeAlreadyThere != true {
		commandTag, err := pgxPool.Conn.Exec(context.Background(), orgunitType)
		if err != nil {
			log.Printf("ERROR: problem creating the «orgunit_type» type : %v", err)
			return nil, errors.New("unable to create the type «orgunit_type» ")
		}
		log.Printf("SUCCESS: «orgunit_type» type was created, rows affected : %v", int(commandTag.RowsAffected()))
	}
	var numberOfOrgUnits int
	errOrgUnitsTable := pgxPool.Conn.QueryRow(context.Background(), orgUnitsCount).Scan(&numberOfOrgUnits)
	if errOrgUnitsTable != nil {
		commandTag, err := pgxPool.Conn.Exec(context.Background(), orgUnitsCreateTable)
		if err != nil {
			log.Printf("ERROR: problem creating the «go_orgunit» table : %v", err)
			return nil, errors.New("unable to create the table «go_orgunit» ")
		}
		log.Printf("SUCCESS: «go_org_unit» table was created rows affected : %v", int(commandTag.RowsAffected()))
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
		log.Printf("INFO: 'will run updateAdminUser for adminUser:%s  in «go_user»'", adminUser)
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
	var lastInsertId int = 0

	goHash, err := crypto.HashAndSalt(u.PasswordHash)
	if err != nil {
		db.log.Printf("error : Create(%q) had an error doing crypto.HashAndSalt. error : %v", u.Name, err)
		return nil, err
	}
	u.PasswordHash = goHash
	err = db.Db.Conn.QueryRow(context.Background(), usersCreate,
		u.Name, u.Email, u.Username, u.PasswordHash, &u.ExternalId, &u.OrgunitId, &u.GroupsId, &u.Phone, //$1-$7
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

// Get will retrieve one user with given id
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
	const usersUpdate = `
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
	rowsAffected, err = db.Db.ExecActionQuery(usersUpdate,
		user.Name, user.Email, user.Username,
		&user.ExternalId, &user.OrgunitId, &user.GroupsId, &user.Phone, user.IsLocked, user.IsAdmin, user.LastModificationUser,
		user.IsActive, &user.InactivationTime, &user.InactivationReason, &user.Comment, id)
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

// Close will close the database store
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
		// return an empty array
		return make([]*GroupList, 1), nil
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
