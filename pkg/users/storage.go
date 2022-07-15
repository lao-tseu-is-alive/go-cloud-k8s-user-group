package users

import (
	"errors"
	"fmt"
	"log"
	"runtime"
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
	Create(todo NewUser) (*User, error)
	// Update updates the users with given ID in the storage.
	Update(id int32, todo User) (*User, error)
	// Delete removes the users with given ID from the storage.
	Delete(id int32) error
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
			return nil, fmt.Errorf("error doing NewPgxDB(dbConnectionString : %v", err)
		}
	default:
		return nil, errors.New("unsupported DB driver type")
	}
	return db, nil
}

func GetErrorF(errMsg string, err error) error {
	return errors.New(fmt.Sprintf("%s [%v]", errMsg, err))
}
