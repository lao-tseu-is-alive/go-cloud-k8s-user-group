package users

import (
	"errors"
	"fmt"
	"log"
	"runtime"
)

const (
	countMetaUserServiceSQL = "SELECT COUNT(*) as num FROM go_metadata_db_schema WHERE service = 'GO_USER';"
	insertMetaUserService   = `INSERT INTO go_metadata_db_schema 
						(service, schema, table_name, version) 
				VALUES 	('GO_USER','public','go_user',$1) RETURNING id;`
)

// Storage is an interface to different implementation of persistence for Users
type Storage interface {
	// List returns the list of existing users with the given offset and limit.
	List(offset, limit int) ([]*User, error)
	// Get returns the users with the specified users ID.
	Get(id int32) (*User, error)
	// GetMaxId returns the maximum value of users id existing in store.
	GetMaxId() (int32, error)
	// Exist returns true only if a users with the specified id exists in store.
	Exist(id int32) bool
	// Count returns the total number of users.
	Count() (int32, error)
	// Create saves a new users in the storage.
	Create(user User) (*User, error)
	// Update updates the users with given ID in the storage.
	Update(id int32, user User) (*User, error)
	// Delete removes the users with given ID from the storage.
	Delete(id int32) error
	// FindUsername retrieves the user id for the given username or err if not found
	FindUsername(username string) (int32, error)
	// Close terminates properly the connection to the backend
	Close()
}

func GetStorageInstance(dbDriver, dbConnectionString string, log *log.Logger) (Storage, error) {
	var db Storage
	var err error
	switch dbDriver {
	case "postgres":
		db, err = NewPgxDB(dbConnectionString, runtime.NumCPU(), log)
		if err != nil {
			return nil, fmt.Errorf("error doing NewPgxDB(dbConnectionString : %w", err)
		}
	default:
		return nil, errors.New("unsupported DB driver type")
	}
	return db, nil
}

func GetErrorF(errMsg string, err error) error {
	return fmt.Errorf("%s [%v]", errMsg, err)
}
