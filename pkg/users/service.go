package users

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type Service struct {
	Log   *log.Logger
	Store Storage
}

func (s Service) ChangeUserPassword(ctx echo.Context, userId int32) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetLogin(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) LoginUser(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetResetPasswordEmail(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) SendResetPassword(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetResetPasswordToken(ctx echo.Context, resetPasswordToken string) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) ResetPassword(ctx echo.Context, resetPasswordToken string) error {
	//TODO implement me
	panic("implement me")
}

type ErrorService struct {
	Err    error  `json:"err" `
	Status int    `json:"status" `
	Msg    string `json:"msg" `
}

func (e *ErrorService) Error() string {
	return fmt.Sprintf("Status[%d] %s. error: %v", e.Status, e.Msg, e.Err)
}

// GetMaxId returns the greatest users id used by now
// curl -H "Content-Type: application/json" 'http://localhost:8080/users/maxid'
func (s Service) GetMaxId(ctx echo.Context) error {
	s.Log.Println("# Entering GetMaxId()")
	var maxUserId int32 = 0
	maxUserId, _ = s.Store.GetMaxId()
	s.Log.Printf("# Exit GetMaxId() maxUserId: %d", maxUserId)
	return ctx.JSON(http.StatusOK, maxUserId)
}

func (s Service) GetUser(ctx echo.Context, userId int32) error {
	s.Log.Printf("# Entering GetUser(%d)", userId)
	if s.Store.Exist(userId) == false {
		return ctx.JSON(http.StatusNotFound, ErrorService{
			Err:    errors.New("not found"),
			Status: http.StatusNotFound,
			Msg:    fmt.Sprintf("user id : %d does not exist", userId),
		})
	}
	user, err := s.Store.Get(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("problem retrieving user :%v", err))
	}
	return ctx.JSON(http.StatusOK, user)
}

//GetUsers will retrieve all Users in the store and return then
//to test it with curl you can try :
//curl -H "Content-Type: application/json" 'http://localhost:8080/users' |json_pp
func (s Service) GetUsers(ctx echo.Context, params GetUsersParams) error {
	s.Log.Printf("# Entering GetUsers() %v", params)
	list, err := s.Store.List(0, 100)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("there was a problem when calling store.List :%v", err))
	}
	return ctx.JSON(http.StatusOK, list)
}

//CreateUser will store the NewUser task in the store
//to test it with curl you can try :
//curl -XPOST -H "Content-Type: application/json" -d '{"task":"learn Linux"}'  'http://localhost:8080/users'
//curl -XPOST -H "Content-Type: application/json" -d '{"task":""}'  'http://localhost:8080/users'
func (s Service) CreateUser(ctx echo.Context) error {
	s.Log.Println("# Entering CreateUser()")
	newUser := &NewUser{}
	if err := ctx.Bind(newUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("CreateUser has invalid format [%v]", err))
	}
	if len(newUser.Name) < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprint("CreateUser task cannot be empty"))
	}
	if len(newUser.Name) < 6 {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprint("CreateUser task minLength is 5"))
	}
	s.Log.Printf("# CreateUser() newUser : %#v\n", newUser)
	userCreated, err := s.Store.Create(*newUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("problem saving new user :%v", err))
	}
	s.Log.Printf("# CreateUser() User %#v\n", userCreated)
	return ctx.JSON(http.StatusCreated, userCreated)

}

// UpdateUser will store the modified information in the store for the given userId
// curl -v -XPUT -H "Content-Type: application/json" -d '{"id": 3, "task":"learn Linux", "completed": true}'  'http://localhost:8080/users/3'
// curl -v -XPUT -H "Content-Type: application/json" -d '{"id": 3, "task":"learn Linux", "completed": false}'  'http://localhost:8080/users/3'
func (s Service) UpdateUser(ctx echo.Context, userId int32) error {
	s.Log.Printf("# Entering UpdateUser(%d)", userId)
	if s.Store.Exist(userId) == false {
		return ctx.JSON(http.StatusNotFound, ErrorService{
			Err:    errors.New("not found"),
			Status: http.StatusNotFound,
			Msg:    fmt.Sprintf("user id : %d does not exist", userId),
		})
	}
	t := new(User)
	if err := ctx.Bind(t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("UpdateUser has invalid format [%v]", err))
	}
	if len(t.Name) < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprint("CreateUser task cannot be empty"))
	}
	//refuse an attempt to modify a userId (in url) with a different id in the body !
	if t.Id != userId {
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf("UpdateUser id : [%d] and posted Id [%d] cannot differ ", userId, t.Id))
	}

	updatedUser, err := s.Store.Update(userId, *t)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("problem updating user :%v", err))
	}
	return ctx.JSON(http.StatusOK, updatedUser)
}

// DeleteUser will remove the given userID entry from the store, and if not present will return 400 Bad Request
//curl -v -XDELETE -H "Content-Type: application/json" 'http://localhost:8080/users/3' ->  204 No Content if present and delete it
//curl -v -XDELETE -H "Content-Type: application/json" 'http://localhost:8080/users/93333' -> 400 Bad Request
func (s Service) DeleteUser(ctx echo.Context, userId int32) error {
	s.Log.Printf("# Entering DeleteUser(%d)", userId)
	if s.Store.Exist(userId) == false {
		return ctx.JSON(http.StatusNotFound, ErrorService{
			Err:    errors.New("not found"),
			Status: http.StatusNotFound,
			Msg:    fmt.Sprintf("user id : %d does not exist", userId),
		})
	} else {
		err := s.Store.Delete(userId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("problem deleting user :%v", err))
		}
		return ctx.NoContent(http.StatusNoContent)
	}
}
