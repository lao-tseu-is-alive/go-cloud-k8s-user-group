package users

import (
	"fmt"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/database"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/golog"
)

const (
	countMetaUserServiceSQL = "SELECT COUNT(*) as num FROM go_metadata_db_schema WHERE service = 'GO_USER';"
	insertMetaUserService   = `INSERT INTO go_metadata_db_schema 
						(service, schema, table_name, version) 
				VALUES 	('GO_USER','public','go_user',$1) RETURNING id;`
)

// Storage is an interface to different implementation of persistence for Users/Groups
type Storage interface {
	// List returns the list of existing users with the given offset and limit.
	List(offset, limit int) ([]*UserList, error)
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
	// IsUserAdmin returns true if the user with the specified id has the is_admin attribute set to true
	IsUserAdmin(id int32) bool
	// IsUserActive returns true if the user with the specified id has the is_active attribute set to true
	IsUserActive(id int32) bool
	// ResetPassword updates the password hash for the user with the specified id
	ResetPassword(id int32, passwordHash string, idCurrentUser int) error
	// CreateGroup saves a new group in the storage.
	CreateGroup(group Group) (*Group, error)
	// UpdateGroup updates the group with given ID in the storage.
	UpdateGroup(id int32, group Group) (*Group, error)
	// DeleteGroup removes the group with given ID from the storage.
	DeleteGroup(id int32) error
	// ListGroup returns the list of active groups with the given offset and limit.
	ListGroup(offset, limit int) ([]*GroupList, error)
	// GetGroup returns the group with the specified users ID.
	GetGroup(id int32) (*Group, error)
}

func GetStorageInstanceOrPanic(dbDriver string, db database.DB, l golog.MyLogger) Storage {
	var store Storage
	var err error
	switch dbDriver {
	case "pgx":
		store = NewPgxDB(db, l)
		if err != nil {
			l.Fatal("error doing NewPgxDB(pgConn : %w", err)
		}

	default:
		panic("unsupported DB driver type")
	}
	return store
}

func GetErrorF(errMsg string, err error) error {
	return fmt.Errorf("%s [%v]", errMsg, err)
}
