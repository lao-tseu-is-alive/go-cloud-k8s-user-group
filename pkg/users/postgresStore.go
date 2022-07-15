package users

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/database"
	"log"
	"time"
)

const (
	usersList       = "SELECT id, name, completed, created_at, completed_at FROM users ORDER BY id;"
	usersGet        = "SELECT id, name, completed, created_at, completed_at FROM users WHERE id=$1;"
	usersCompleted  = "SELECT completed FROM users WHERE id=$1"
	usersExist      = "SELECT COUNT(*) FROM users WHERE id=$1"
	usersCount      = "SELECT COUNT(*) FROM users"
	usersMaxId      = "SELECT MAX(id) FROM users"
	usersCreate     = "INSERT INTO users (name) VALUES($1) RETURNING id;"
	usersUpdate     = "UPDATE users SET name=$1, completed=$2, completed_at=$3 WHERE id=$4"
	usersUpdateName = "UPDATE users SET name=$1 WHERE id=$2"
	usersDelete     = "DELETE FROM users WHERE id = $1"
	insertAdminUser = `INSERT INTO go_user (name, email, username, password_hash, is_admin, creator, comment) 
	VALUES ($1,$2,$3,$4, true, 1, 'Initial initDBSchema Setup Admin account');`
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
	} else {
		log.Printf("INFO: 'Connection to database : %s'\n", successOrFailure)
		var numberOfUsers int
		if errTodosTable := pgxPool.Conn.QueryRow(context.Background(), usersCount).Scan(&numberOfUsers); errTodosTable != nil {
			log.Printf("ERROR: 'the database does not contain the table «users» !'")
			return nil, errors.New("database does not contain the table «users»")
		}
		log.Printf("INFO: 'database contains %d records in users'", numberOfUsers)
	}

	psql.Db = pgxPool
	psql.log = log
	return &psql, err
}

//Create will store the new name in the store
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
		res := &User{
			Completed:   false,
			CompletedAt: nil,
			CreatedAt:   nil,
			Id:          0,
			Name:        "",
		}
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
		updateAll := true
		var rowsAffected int = 0
		var err error
		now := time.Now()
		// implements basic Business Rules !
		// let's first check if name was already completed in DB
		alreadyCompleted, _ := db.Db.GetQueryBool(usersCompleted, id)
		switch user.Completed {
		case true:
			if alreadyCompleted == false {
				// this name was not completed, now it is, so update CompletedAt
				user.CompletedAt = &now
			}
		case false:
			if alreadyCompleted == true {
				// this task was completed, but user changed it to not completed so update with nil
				user.CompletedAt = nil
			}
		default:
			// in all other cases the values of Completed and CompletedAt fields should not be changed in DB
			// so here let's update only the Name field
			rowsAffected, err = db.Db.ExecActionQuery(usersUpdateName, user.Name, id)
			updateAll = false
		}
		if updateAll {
			rowsAffected, err = db.Db.ExecActionQuery(usersUpdate, user.Name, user.Completed, user.CompletedAt, id)
		}
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
