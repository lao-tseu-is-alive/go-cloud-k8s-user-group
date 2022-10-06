package users

import (
	"encoding/json"
	"fmt"
	"github.com/cristalhq/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/crypto"
	"log"
	"net/http"
	"time"
)

type Service struct {
	Log         *log.Logger
	Store       Storage
	JwtSecret   []byte
	JwtDuration int
}

// JwtCustomClaims are custom claims extending default ones.
type JwtCustomClaims struct {
	jwt.RegisteredClaims
	Id       int32  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
}

// UsersCreate will store the NewUser task in the store
/* to test it with curl you can try :
curl -s -XPOST -H "Content-Type: application/json" -H "Authorization: Bearer $token" \
-d '{"username":"cgil", "name":"Carlos GIL", "email":"c@gil.town", "password_hash":"4acf0b39d9c4766709a3689f553ac01ab550545ffa4544dfc0b2cea82fba02a3"}'  'http://localhost:8888/api/users'
*/
func (s Service) UsersCreate(ctx echo.Context) error {
	s.Log.Println("trace: entering CreateUser()")
	/* uncomment when jw is implemented
	// get the current user from JWT TOKEN
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*MyCustomJWTClaims)
	// IF USER IS NOT ADMIN RETURN 401 Unauthorized
	currentUserId := claims.ID
	if !s.Store.IsUserAdmin(currentUserId) {
		return echo.NewHTTPError(http.StatusUnauthorized, "current user has no admin privilege")
	}
	*/
	currentUserId := 1
	newUser := &User{
		Id:      0,
		Creator: int32(currentUserId),
	}
	if err := ctx.Bind(newUser); err != nil {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("CreateUser has invalid format [%v]", err))
	}
	if len(newUser.Name) < 1 {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprint("CreateUser name cannot be empty"))
	}
	if len(newUser.Name) < 6 {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprint("CreateUser name minLength is 5"))
	}
	if len(newUser.Username) < 1 {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprint("CreateUser username cannot be empty"))
	}
	if len(newUser.Username) < 6 {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprint("CreateUser username minLength is 5"))
	}
	if !crypto.ValidatePasswordHash(newUser.PasswordHash) {
		msg := fmt.Sprintf("CreateUser received invalid password hash in request body")
		s.Log.Printf(msg)
		return ctx.JSON(http.StatusBadRequest, msg)
	}
	s.Log.Printf("# CreateUser() newUser : %#v\n", newUser)
	userCreated, err := s.Store.Create(*newUser)
	if err != nil {
		msg := fmt.Sprintf("CreateUser had an error saving user:%#v, err:%#v", *newUser, err)
		s.Log.Printf(msg)
		return ctx.JSON(http.StatusBadRequest, msg)
	}
	s.Log.Printf("# CreateUser() User %#v\n", userCreated)
	return ctx.JSON(http.StatusCreated, userCreated)

}

// GetMaxId returns the greatest users id used by now
// curl -H "Content-Type: application/json" 'http://localhost:8888/users/maxid'
func (s Service) GetMaxId(ctx echo.Context) error {
	s.Log.Println("trace: entering GetMaxId()")
	var maxUserId int32 = 0
	maxUserId, _ = s.Store.GetMaxId()
	s.Log.Printf("# Exit GetMaxId() maxUserId: %d", maxUserId)
	return ctx.JSON(http.StatusOK, maxUserId)
}

// UsersGet will retrieve the User in the store and return then
// to test it with curl you can try :
// curl -s -H "Content-Type: application/json" -H "Authorization: Bearer $token" 'http://localhost:8888/api/users' |jq
func (s Service) UsersGet(ctx echo.Context, userId int32) error {
	s.Log.Printf("trace: entering GetUser(%d)", userId)
	if s.Store.Exist(userId) == false {
		msg := fmt.Sprintf("GetUser(%d) this id does not exist.", userId)
		s.Log.Printf(msg)
		return ctx.JSON(http.StatusNotFound, msg)
	}
	user, err := s.Store.Get(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("problem retrieving user :%v", err))
	}
	return ctx.JSON(http.StatusOK, user)
}

// UsersList will retrieve all Users in the store and return then
// to test it with curl you can try :
// curl -s -H "Content-Type: application/json" -H "Authorization: Bearer $token" 'http://localhost:8888/api/users' |jq
func (s Service) UsersList(ctx echo.Context, params UsersListParams) error {
	s.Log.Printf("trace: entering GetUsers() params:%v", params)
	list, err := s.Store.List(0, 100)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("there was a problem when calling store.List :%v", err))
	}
	return ctx.JSON(http.StatusOK, list)
}

// UsersDelete will remove the given userID entry from the store, and if not present will return 400 Bad Request
// curl -v -XDELETE -H "Content-Type: application/json" -H "Authorization: Bearer $token" 'http://localhost:8888/api/users/3' ->  204 No Content if present and delete it
// curl -v -XDELETE -H "Content-Type: application/json"  -H "Authorization: Bearer $token" 'http://localhost:8888/users/93333' -> 400 Bad Request
func (s Service) UsersDelete(ctx echo.Context, userId int32) error {
	s.Log.Printf("trace: entering DeleteUser(%d)", userId)
	/* uncomment when jw is implemented
	// get the current user from JWT TOKEN
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*MyCustomJWTClaims)
	// IF USER IS NOT ADMIN RETURN 401 Unauthorized
	currentUserId := claims.ID
	if !s.Store.IsUserAdmin(currentUserId) {
		return echo.NewHTTPError(http.StatusUnauthorized, "current user has no admin privilege")
	}
	*/
	if userId == 1 {
		msg := fmt.Sprintln("DeleteUser cannot delete the original admin")
		s.Log.Printf(msg)
		return echo.NewHTTPError(http.StatusBadRequest, msg)
	}
	if s.Store.Exist(userId) == false {
		msg := fmt.Sprintf("DeleteUser(%d) cannot delete this id, it does not exist !", userId)
		s.Log.Printf(msg)
		return ctx.JSON(http.StatusNotFound, msg)
	} else {
		err := s.Store.Delete(userId)
		if err != nil {
			msg := fmt.Sprintf("DeleteUser(%d) got an error: %#v ", userId, err)
			s.Log.Printf(msg)
			return echo.NewHTTPError(http.StatusInternalServerError, msg)
		}
		return ctx.NoContent(http.StatusNoContent)
	}
}

// UsersUpdate will store the modified information in the store for the given userId
// curl -v -XPUT -H "Content-Type: application/json" -d '{"id": 3, "task":"learn Linux", "completed": true}'  'http://localhost:8888/users/3'
// curl -v -XPUT -H "Content-Type: application/json" -d '{"id": 3, "task":"learn Linux", "completed": false}'  'http://localhost:8888/users/3'
func (s Service) UsersUpdate(ctx echo.Context, userId int32) error {
	s.Log.Printf("trace: entering UpdateUser(%d)", userId)
	/* uncomment when jw is implemented
	// get the current user from JWT TOKEN
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*MyCustomJWTClaims)
	// IF USER IS NOT ADMIN RETURN 401 Unauthorized
	currentUserId := claims.ID
	if !s.Store.IsUserAdmin(currentUserId) {
		return echo.NewHTTPError(http.StatusUnauthorized, "current user has no admin privilege")
	}
	*/
	if s.Store.Exist(userId) == false {
		msg := fmt.Sprintf("UpdateUser(%d) cannot modify this id, it does not exist.", userId)
		s.Log.Printf(msg)
		return ctx.JSON(http.StatusNotFound, msg)
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
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("UpdateUser got problem updating user :%v", err))
	}
	return ctx.JSON(http.StatusOK, updatedUser)
}

// UsersChangePassword allows a user to change it's own password
// (PUT /api/users/{userId}/changepassword)
func (s Service) UsersChangePassword(ctx echo.Context, userId int32) error {
	s.Log.Printf("trace: entering ChangeUserPassword(%d)", userId)
	//TODO implement me
	panic("implement me")
}

/////////////////////////// HANDLERS WITHOUT JWT

// GetLogin allows client to do a preflight prepare for a login
// (GET /login)
func (s Service) GetLogin(ctx echo.Context) error {
	s.Log.Println("trace: entering GetLogin()")
	return ctx.JSON(http.StatusOK, "you must post login credentials")
}

// LoginUser allows client to try to authenticate, and then receive a valid JWT
// curl -X POST -H "Content-Type: application/json" -d '{"username": "go-admin", "password_hash": "your_pwd_hash" }'  http://localhost:8888/login
// with the received token you can try : curl  -H "Authorization: Bearer $token "  http://localhost:8888/restricted
func (s Service) LoginUser(ctx echo.Context) error {
	s.Log.Println("trace: entering LoginUser()")

	uLogin := new(UserLogin)
	if err := ctx.Bind(uLogin); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user login or json format in request body")
	}
	idUser, err := s.Store.FindUsername(uLogin.Username)
	if err != nil {
		if err == ErrUsernameNotFound {
			s.Log.Printf("LoginUser(%s) username was not found in DB.", uLogin.Username)
			return ctx.JSON(http.StatusUnauthorized, "username not found")
		}
		msg := fmt.Sprintf("LoginUser(%s) s.Store.FindUsername got an error: %#v ", uLogin.Username, err)
		s.Log.Printf(msg)
		ctx.JSON(http.StatusNotFound, msg)
		return echo.NewHTTPError(http.StatusBadRequest, msg)
	}
	user, err := s.Store.Get(idUser)
	if err != nil {
		msg := fmt.Sprintf("LoginUser(%s) s.Store.Get(%d) got an error: %#v ", uLogin.Username, idUser, err)
		s.Log.Printf(msg)
		return echo.NewHTTPError(http.StatusInternalServerError, msg)
	}
	if !user.IsActive {
		msg := fmt.Sprintf("LoginUser(%s) this user id (%d) is not active anymore", uLogin.Username, idUser)
		s.Log.Printf(msg)
		return echo.NewHTTPError(http.StatusUnauthorized, msg)
	}
	if user.IsLocked {
		msg := fmt.Sprintf("LoginUser(%s) this user id (%d) is locked", uLogin.Username, idUser)
		s.Log.Printf(msg)
		return echo.NewHTTPError(http.StatusUnauthorized, msg)
	}
	if !crypto.ComparePasswords(user.PasswordHash, uLogin.PasswordHash) {
		msg := fmt.Sprintf("LoginUser(%s) ComparePasswords failed for user id (%d). wrong password !", uLogin.Username, idUser)
		s.Log.Printf(msg)
		//TODO increment bad password count and if max bad passwords reached then lock this account
		return echo.NewHTTPError(http.StatusUnauthorized, msg)
	}

	// Set custom claims
	claims := &JwtCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        "",
			Audience:  nil,
			Issuer:    "",
			Subject:   "",
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Minute * time.Duration(s.JwtDuration))},
			IssuedAt:  nil,
			NotBefore: nil,
		},
		Id:       user.Id,
		Name:     user.Name,
		Email:    string(user.Email),
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}

	// Create token with claims
	signer, _ := jwt.NewSignerHS(jwt.HS512, s.JwtSecret)
	builder := jwt.NewBuilder(signer)
	token, err := builder.Build(claims)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("LoginUser(%s) succesfull login for user id (%d)", uLogin.Username, idUser)
	s.Log.Printf(msg)
	return ctx.JSON(http.StatusOK, echo.Map{
		"token": token.String(),
	})
}

func (s Service) GetResetPasswordEmail(ctx echo.Context) error {
	s.Log.Println("trace: entering GetResetPasswordEmail()")
	//TODO implement me
	panic("implement me")
}

func (s Service) SendResetPassword(ctx echo.Context) error {
	s.Log.Println("trace: entering SendResetPassword()")
	//TODO implement me
	panic("implement me")
}

func (s Service) GetResetPasswordToken(ctx echo.Context) error {
	s.Log.Println("trace: entering GetResetPasswordToken()")
	//TODO implement me
	panic("implement me")
}

func (s Service) ResetPassword(ctx echo.Context) error {
	s.Log.Println("trace: entering ResetPassword()")
	//TODO implement me
	panic("implement me")
}

func (s Service) GetStatus(ctx echo.Context) error {
	s.Log.Println("trace: entering GetStatus()")
	u := ctx.Get("jwtdata").(*jwt.Token)
	claims := JwtCustomClaims{}
	err := u.DecodeClaims(&claims)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	username := claims.Username
	idUser := claims.Id
	res, err := json.Marshal(claims)
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, "JWT User Data Could Not Be Marshaled To Json")
	}
	s.Log.Printf("info: GetStatus(user:%s, id:%d)", username, idUser)
	return ctx.JSONBlob(http.StatusOK, res)
}
