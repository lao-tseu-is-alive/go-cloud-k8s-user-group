package users

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/config"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/crypto"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/database"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/version"
	"log"
	"time"
)

const (
	usersList = "SELECT id, name, email, username, create_time, is_admin, is_locked FROM go_users ORDER BY id;"
	usersGet  = `
SELECT id, name, email, username,
       password_hash, external_id, enterprise, phone, is_locked, is_admin,
       create_time, creator, last_modification_time, last_modification_user, 
       is_active, inactivation_time, inactivation_reason, comment, bad_password_count
FROM go_users WHERE id=$1;`

	usersExist  = "SELECT COUNT(*) FROM go_users WHERE id=$1"
	usersCount  = "SELECT COUNT(*) FROM go_users"
	usersMaxId  = "SELECT MAX(id) FROM go_users"
	usersCreate = "INSERT INTO go_users (name) VALUES($1) RETURNING id;"
	usersUpdate = `
UPDATE go_users
SET name                   = $1,
    email                  = $2,
    username               = $3,
    enterprise             = $4,
    phone                  = $5,
    is_locked              = $6,
    is_admin               = $7,
    last_modification_time = CURRENT_TIMESTAMP,
    last_modification_user = $8,
    is_active              = $9,
    inactivation_time      = $10,
    inactivation_reason    = $11,
    comment                = $12,
    password_hash          = $13,
    external_id   		   = $14    
WHERE id = $15;
`
	usersDelete      = "DELETE FROM go_users WHERE id = $1"
	usersCreateTable = `
CREATE TABLE IF NOT EXISTS go_users
(
    id        	serial    CONSTRAINT go_users_pk   primary key,
    name			text	not null	constraint go_user_unique_name	unique
        								constraint name_min_length check (length(btrim(name)) > 2),
    email			text	not null	constraint go_user_unique_email unique
										constraint email_min_length	check (length(btrim(email)) > 3),
    username		text	not null	constraint go_user_unique_username unique
										constraint username_min_length check (length(btrim(username)) > 2),
    password_hash	text	not null 	constraint password_hash_min_length check (length(btrim(password_hash)) > 30),
	external_id		text,
    enterprise		text,
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
comment on table go_users is 'go_user is the main table of the GO_USER microservice';
`
	insertAdminUser = `
INSERT INTO go_users (name, email, username, password_hash, is_admin, creator, comment) 
VALUES ('Initial Administrative Account','admin@example.com',$1,$2, true, 1, 'Initial Setup of Admin account')  RETURNING id;`
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
	if numberOfUsers > 0 {
		log.Printf("INFO: 'database contains %d records in «go_user»'", numberOfUsers)
		//TODO: update the password for admin user id=1 with actual env value
	} else {
		log.Printf("WARNING: '«go_user» contain %d records : creating go-admin user'", numberOfUsers)
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

// Create will store the new name in the store
func (db *PGX) Create(user NewUser) (*User, error) {
	db.log.Printf("info : Entering Create(%#v)", user)
	if len(user.Name) < 1 {
		return nil, errors.New("user name cannot be empty")
	}
	if len(user.Name) < 6 {
		return nil, errors.New("CreateUser name minLength is 5")
	}
	var lastInsertId int = 0
	err := db.Db.Conn.QueryRow(context.Background(), usersCreate, user.Name).Scan(&lastInsertId)
	if err != nil {
		db.log.Printf("error : Create(%v) unexpectedly failed. error : %v", user.Name, err)
		return nil, err
	}
	db.log.Printf("info : Create(%v) created with id : %v", user.Name, lastInsertId)

	// if we get to here all is good, so let's retrieve a fresh copy to send it back
	createdUser, err := db.Get(int32(lastInsertId))
	if err != nil {
		return nil, GetErrorF("error : users was created, but can not be retrieved", err)
	}
	return createdUser, nil
}

func (db *PGX) List(offset, limit int) ([]*User, error) {
	var res []*User

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
	db.log.Printf("info : Get(%d) entering...", id)
	if db.Exist(id) == true {
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
	db.log.Printf("info : Get(%d) id does not exist", id)
	return nil, errors.New("user with this id does not exist")
}

// GetMaxId returns the maximum value of users id existing in store.
func (db *PGX) GetMaxId() (int32, error) {
	existingMaxId, err := db.Db.GetQueryInt(usersMaxId)
	if err != nil {
		db.log.Printf("getMaxId() could not be retrieved from DB. failed db.Query err: %v", err)
		return 0, err
	}
	return int32(existingMaxId), nil
}

// Exist returns true only if a users with the specified id exists in store.
func (db *PGX) Exist(id int32) bool {
	count, err := db.Db.GetQueryInt(usersExist, id)
	if err != nil {
		db.log.Printf("exist(%d) could not be retrieved from DB. failed db.Query err: %v", id, err)
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
	count, err := db.Db.GetQueryInt(usersCount)
	if err != nil {
		db.log.Printf("count(*) could not be retrieved from DB. failed db.Query err: %v", err)
		return 0, err
	}
	return int32(count), nil
}

// Update the users stored in DB with given id and other information in struct
func (db *PGX) Update(id int32, user User) (*User, error) {
	if db.Exist(id) {
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
		rowsAffected, err = db.Db.ExecActionQuery(usersUpdate, user.Name, user.IsActive, user.LastModificationTime, id)

		if err != nil {
			return nil, GetErrorF("error : users could not be updated", err)
		}
		if rowsAffected < 1 {
			return nil, GetErrorF("error : users was not updated", err)
		}
		// if we get to here all is good, so let's retrieve a fresh copy to send it back
		updatedUser, err := db.Get(id)
		if err != nil {
			return nil, GetErrorF("error : users was updated, but can not be retrieved", err)
		}
		return updatedUser, nil
	}
	db.log.Printf("info : Update(%d) id does not exist", id)
	return nil, errors.New("user with this id does not exist")
}

// Delete the users stored in DB with given id
func (db *PGX) Delete(id int32) error {
	if db.Exist(id) {
		rowsAffected, err := db.Db.ExecActionQuery(usersDelete, id)
		if err != nil {
			return GetErrorF("error : users could not be deleted", err)
		}
		if rowsAffected < 1 {
			return GetErrorF("error : users was not deleted", err)
		}
		// if we get to here all is good
		return nil
	}
	db.log.Printf("info : Delete(%d) id does not exist", id)
	return errors.New("user with this id does not exist")
}

func (db *PGX) Close() {
	db.Db.Close()
}
