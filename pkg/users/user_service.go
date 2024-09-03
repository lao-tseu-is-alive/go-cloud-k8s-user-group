package users

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/database"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/goHttpEcho"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/golog"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/crypto"
	"net/http"
)

type Service struct {
	Logger golog.MyLogger
	DbConn database.DB
	Store  Storage
	Server *goHttpEcho.Server
}

// UserCreate will store a new User in the store
/* to test it with curl you can try :
curl -s -XPOST -H "Content-Type: application/json" -H "Authorization: Bearer $token" \
-d '{"username":"cgil", "name":"Carlos GIL", "email":"c@gil.town", "password_hash":"4acf0b39d9c4766709a3689f553ac01ab550545ffa4544dfc0b2cea82fba02a3"}'  'http://localhost:8888/api/users'
*/
func (s Service) UserCreate(ctx echo.Context) error {
	goHttpEcho.TraceRequest("UserCreate", ctx.Request(), s.Logger)
	// get the current user from JWT TOKEN
	claims := s.Server.JwtCheck.GetJwtCustomClaimsFromContext(ctx)
	currentUserId := claims.User.UserId
	s.Logger.Info("in UserCreate : currentUserId: %d", currentUserId)
	// IF USER IS NOT ADMIN RETURN 401 Unauthorized
	if !s.Store.IsUserAdmin(int32(currentUserId)) {
		return echo.NewHTTPError(http.StatusUnauthorized, "current user has no admin privilege")
	}
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
	if len(newUser.Name) < 5 {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("CreateUser name minLength is 5 not (%d)", len(newUser.Name)))
	}
	if len(newUser.Username) < 1 {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprint("CreateUser username cannot be empty"))
	}
	if len(newUser.Username) < 5 {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprint("CreateUser username minLength is 5"))
	}
	if !crypto.ValidatePasswordHash(newUser.PasswordHash) {
		msg := fmt.Sprintf("CreateUser received invalid password hash in request body")
		s.Logger.Info(msg)
		return ctx.JSON(http.StatusBadRequest, msg)
	}
	s.Logger.Info("# CreateUser() newUser : %#v\n", newUser)
	userCreated, err := s.Store.Create(*newUser)
	if err != nil {
		msg := fmt.Sprintf("CreateUser had an error saving user:%#v, err:%#v", *newUser, err)
		s.Logger.Info(msg)
		return ctx.JSON(http.StatusBadRequest, msg)
	}
	s.Logger.Info("# CreateUser() User %#v\n", userCreated)
	return ctx.JSON(http.StatusCreated, userCreated)

}

// GetMaxId returns the greatest users id used by now
// curl -H "Content-Type: application/json" 'http://localhost:8888/users/maxid'
func (s Service) GetMaxId(ctx echo.Context) error {
	s.Logger.Debug("trace: entering GetMaxId()")
	var maxUserId int32 = 0
	maxUserId, _ = s.Store.GetMaxId()
	s.Logger.Info("# Exit GetMaxId() maxUserId: %d", maxUserId)
	return ctx.JSON(http.StatusOK, maxUserId)
}

// UserGet will retrieve the User in the store and return then
// to test it with curl you can try :
// curl -s -H "Content-Type: application/json" -H "Authorization: Bearer $token" 'http://localhost:8888/api/users' |jq
func (s Service) UserGet(ctx echo.Context, userId int32) error {
	goHttpEcho.TraceRequest("UserGet", ctx.Request(), s.Logger)
	// get the current user from JWT TOKEN
	claims := s.Server.JwtCheck.GetJwtCustomClaimsFromContext(ctx)
	currentUserId := int32(claims.User.UserId)
	// IF USER IS NOT ADMIN RETURN 401 Unauthorized IF GET is for ANOTHER id
	if !s.Store.IsUserAdmin(currentUserId) && currentUserId != userId {
		return ctx.JSON(http.StatusUnauthorized, "current user has no admin privilege")
	}
	if s.Store.Exist(userId) == false {
		msg := fmt.Sprintf("UserGet(%d) this id does not exist.", userId)
		s.Logger.Info(msg)
		return ctx.JSON(http.StatusNotFound, msg)
	}
	user, err := s.Store.Get(userId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("problem retrieving user :%v", err))
	}
	return ctx.JSON(http.StatusOK, user)
}

// UserList will retrieve all Users in the store and return then
// to test it with curl you can try :
// curl -s -H "Content-Type: application/json" -H "Authorization: Bearer $token" 'http://localhost:8888/api/users' |jq
func (s Service) UserList(ctx echo.Context, params UserListParams) error {
	goHttpEcho.TraceRequest("UserList", ctx.Request(), s.Logger)
	// get the current user from JWT TOKEN
	claims := s.Server.JwtCheck.GetJwtCustomClaimsFromContext(ctx)
	currentUserId := int32(claims.User.UserId)
	s.Logger.Info("in UserCreate : currentUserId: %d", currentUserId)
	list, err := s.Store.List(0, 100)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("there was a problem when calling store.List :%v", err))
	}
	return ctx.JSON(http.StatusOK, list)
}

// UserDelete will remove the given userID entry from the store, and if not present will return 400 Bad Request
// curl -v -XDELETE -H "Content-Type: application/json" -H "Authorization: Bearer $token" 'http://localhost:8888/api/users/3' ->  204 No Content if present and delete it
// curl -v -XDELETE -H "Content-Type: application/json"  -H "Authorization: Bearer $token" 'http://localhost:8888/users/93333' -> 400 Bad Request
func (s Service) UserDelete(ctx echo.Context, userId int32) error {
	handlerName := "UserChangePassword"
	goHttpEcho.TraceRequest(handlerName, ctx.Request(), s.Logger)
	// get the current user from JWT TOKEN
	claims := s.Server.JwtCheck.GetJwtCustomClaimsFromContext(ctx)
	currentUserId := int32(claims.User.UserId)
	s.Logger.Info("in UserDelete : currentUserId: %d", currentUserId)
	// IF USER IS NOT ADMIN RETURN 401 Unauthorized
	if !s.Store.IsUserAdmin(currentUserId) {
		return ctx.JSON(http.StatusUnauthorized, "current user has no admin privilege")
	}
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
		msg := fmt.Sprintf("%s cannot delete the original admin", handlerName)
		s.Logger.Info(msg)
		return ctx.JSON(http.StatusBadRequest, msg)
	}
	if s.Store.Exist(userId) == false {
		msg := fmt.Sprintf("%s(%d) cannot delete this id, it does not exist !", handlerName, userId)
		s.Logger.Info(msg)
		return ctx.JSON(http.StatusNotFound, msg)
	} else {
		err := s.Store.Delete(userId)
		if err != nil {
			msg := fmt.Sprintf("%s(%d) got an error: %#v ", handlerName, userId, err)
			s.Logger.Info(msg)
			return ctx.JSON(http.StatusInternalServerError, msg)
		}
		return ctx.NoContent(http.StatusNoContent)
	}
}

// UserUpdate will store the modified information in the store for the given userId
// curl -v -XPUT -H "Content-Type: application/json" -d '{"id": 3, "task":"learn Linux", "completed": true}'  'http://localhost:8888/users/3'
// curl -v -XPUT -H "Content-Type: application/json" -d '{"id": 3, "task":"learn Linux", "completed": false}'  'http://localhost:8888/users/3'
func (s Service) UserUpdate(ctx echo.Context, userId int32) error {
	goHttpEcho.TraceRequest("UserUpdate", ctx.Request(), s.Logger)
	// get the current user from JWT TOKEN
	claims := s.Server.JwtCheck.GetJwtCustomClaimsFromContext(ctx)
	currentUserId := int32(claims.User.UserId)
	s.Logger.Info("in UserUpdate : currentUserId: %d", currentUserId)
	// IF USER IS NOT ADMIN RETURN 401 Unauthorized
	if !s.Store.IsUserAdmin(currentUserId) {
		return ctx.JSON(http.StatusUnauthorized, "current user has no admin privilege")
	}
	if s.Store.Exist(userId) == false {
		msg := fmt.Sprintf("UpdateUser(%d) cannot modify this id, it does not exist.", userId)
		s.Logger.Info(msg)
		return ctx.JSON(http.StatusNotFound, msg)
	}
	t := new(User)
	if err := ctx.Bind(t); err != nil {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("UpdateUser has invalid format [%v]", err))
	}
	t.LastModificationUser = &currentUserId
	if len(t.Name) < 1 {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprint("UsersUpdate name cannot be empty"))
	}
	if len(t.Name) < 5 {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprint("UsersUpdate name minLength is 5"))
	}
	if len(t.Username) < 1 {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprint("UsersUpdate username cannot be empty"))
	}
	if len(t.Username) < 5 {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprint("UsersUpdate username minLength is 5"))
	}
	//refuse an attempt to modify a userId (in url) with a different id in the body !
	if t.Id != userId {
		return ctx.JSON(http.StatusBadRequest,
			fmt.Sprintf("UpdateUser id : [%d] and posted Id [%d] cannot differ ", userId, t.Id))
	}

	updatedUser, err := s.Store.Update(userId, *t)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("UpdateUser got problem updating user :%v", err))
	}
	return ctx.JSON(http.StatusOK, updatedUser)
}

// UserChangePassword allows a user to change it's own password
// (PUT /api/users/{userId}/changepassword)
func (s Service) UserChangePassword(ctx echo.Context, userId int32) error {
	handlerName := "UserChangePassword"
	goHttpEcho.TraceRequest(handlerName, ctx.Request(), s.Logger)
	// get the current user from JWT TOKEN
	claims := s.Server.JwtCheck.GetJwtCustomClaimsFromContext(ctx)
	currentUserId := int32(claims.User.UserId)
	s.Logger.Info("in %s : currentUserId: %d", handlerName, currentUserId)
	UserLogin := new(UserLogin)
	if err := ctx.Bind(UserLogin); err != nil {
		return ctx.JSON(http.StatusBadRequest, "invalid user login or json format in request body")
	}
	if (currentUserId != 1) && (userId != currentUserId) {
		msg := fmt.Sprintf("%s(%d) cannot change password of another user id (%d)", handlerName, currentUserId, userId)
		s.Logger.Info(msg)
		return ctx.JSON(http.StatusUnauthorized, msg)
	}
	if s.Store.Exist(userId) == false {
		msg := fmt.Sprintf("%s(%d) cannot delete this id, it does not exist !", handlerName, userId)
		s.Logger.Info(msg)
		return ctx.JSON(http.StatusNotFound, msg)
	} else {
		err := s.Store.ResetPassword(userId, UserLogin.PasswordHash, int(currentUserId))
		if err != nil {
			msg := fmt.Sprintf("%s(%d) got an error: %#v ", handlerName, userId, err)
			s.Logger.Info(msg)
			return ctx.JSON(http.StatusInternalServerError, msg)
		}
		return ctx.JSON(http.StatusOK, "password changed")
	}
}

/////////////////////////// HANDLERS WITHOUT JWT

// GetLogin allows client to do a preflight prepare for a login
// (GET /login)
func (s Service) GetLogin(ctx echo.Context) error {
	s.Logger.Debug("trace: entering GetLogin()")
	return ctx.JSON(http.StatusOK, "you must post login credentials")
}

// LoginUser allows client to try to authenticate, and then receive a valid JWT
// curl -X POST -H "Content-Type: application/json" -d '{"username": "go-admin", "password_hash": "your_pwd_hash" }'  http://localhost:8888/login
// with the received token you can try : curl  -H "Authorization: Bearer $token "  http://localhost:8888/restricted
func (s Service) LoginUser(ctx echo.Context) error {
	goHttpEcho.TraceRequest("login", ctx.Request(), s.Logger)
	//TODO: check if redirect_uri is passed as parameter in url, if it is then at the end do a redirect to this uri with the response
	uLogin := new(UserLogin)
	if err := ctx.Bind(uLogin); err != nil {
		return ctx.JSON(http.StatusBadRequest, "invalid user login or json format in request body")
	}
	idUser, err := s.Store.FindUsername(uLogin.Username)
	if err != nil {
		if errors.Is(err, ErrUsernameNotFound) {
			s.Logger.Info("LoginUser(%s) username was not found in DB.", uLogin.Username)
			return ctx.JSON(http.StatusUnauthorized, "username not found")
		}
		msg := fmt.Sprintf("LoginUser(%s) s.Store.FindUsername got an error: %#v ", uLogin.Username, err)
		s.Logger.Info(msg)
		return ctx.JSON(http.StatusBadRequest, msg)
	}
	user, err := s.Store.Get(idUser)
	if err != nil {
		msg := fmt.Sprintf("LoginUser(%s) s.Store.Get(%d) got an error: %#v ", uLogin.Username, idUser, err)
		s.Logger.Info(msg)
		return ctx.JSON(http.StatusInternalServerError, msg)
	}
	if !user.IsActive {
		msg := fmt.Sprintf("LoginUser(%s) this user id (%d) is not active anymore", uLogin.Username, idUser)
		s.Logger.Info(msg)
		return ctx.JSON(http.StatusUnauthorized, msg)
	}
	if user.IsLocked {
		msg := fmt.Sprintf("LoginUser(%s) this user id (%d) is locked", uLogin.Username, idUser)
		s.Logger.Info(msg)
		return ctx.JSON(http.StatusUnauthorized, msg)
	}
	if !crypto.ComparePasswords(user.PasswordHash, uLogin.PasswordHash) {
		msg := fmt.Sprintf("LoginUser(%s) ComparePasswords failed for user id (%d). wrong password !", uLogin.Username, idUser)
		s.Logger.Info(msg)
		//TODO increment bad password count and if max bad passwords reached then lock this account
		return ctx.JSON(http.StatusUnauthorized, msg)
	}

	externalId := 0
	if user.ExternalId != nil {
		externalId = int(*user.ExternalId)
	}

	userInfo := &goHttpEcho.UserInfo{
		UserId:     int(user.Id),
		ExternalId: externalId,
		Name:       user.Name,
		Email:      string(user.Email),
		Login:      user.Username,
		IsAdmin:    user.IsAdmin,
	}
	token, err := s.Server.JwtCheck.GetTokenFromUserInfo(userInfo)
	if err != nil {
		errGetUInfFromLogin := fmt.Sprintf("Error getting jwt token from user info: %v", err)
		s.Logger.Error(errGetUInfFromLogin)
		return ctx.JSON(http.StatusInternalServerError, errGetUInfFromLogin)
	}
	// Prepare the response
	response := map[string]string{
		"token": token.String(),
	}
	s.Logger.Info("LoginUser(%s) successful login", uLogin.Username)
	return ctx.JSON(http.StatusOK, response)
	/*
		// Set custom claims
		claims := &JwtCustomClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ID:        "",
				Audience:  nil,
				Issuer:    "",
				Subject:   "",
				ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Minute * time.Duration(s.JwtDuration))},
				IssuedAt:  &jwt.NumericDate{Time: time.Now()},
				NotBefore: nil,
			},
			Id:       *user.ExternalId,
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
		s.Logger.Info(msg)
		return ctx.JSON(http.StatusOK, echo.Map{
			"token": token.String(),
		})
	*/
}

func (s Service) GetResetPasswordEmail(ctx echo.Context) error {
	s.Logger.Debug("trace: entering GetResetPasswordEmail()")
	//TODO implement me
	panic("implement me")
}

func (s Service) SendResetPassword(ctx echo.Context) error {
	s.Logger.Debug("trace: entering SendResetPassword()")
	//TODO implement me
	panic("implement me")
}

func (s Service) GetResetPasswordToken(ctx echo.Context) error {
	s.Logger.Debug("trace: entering GetResetPasswordToken()")
	//TODO implement me
	panic("implement me")
}

func (s Service) ResetPassword(ctx echo.Context) error {
	s.Logger.Debug("trace: entering ResetPassword()")
	//TODO implement me
	panic("implement me")
}

func (s Service) GetStatus(ctx echo.Context) error {
	handlerName := "GetStatus"
	goHttpEcho.TraceRequest(handlerName, ctx.Request(), s.Logger)
	// get the current user from JWT TOKEN
	claims := s.Server.JwtCheck.GetJwtCustomClaimsFromContext(ctx)
	currentUserId := claims.User.UserId
	s.Logger.Info("in %s : currentUserId: %d", handlerName, currentUserId)

	username := claims.User.Name
	idUser := claims.User.UserId

	s.Logger.Info("%s(user:%s, id:%d)", handlerName, username, idUser)
	return ctx.JSON(http.StatusOK, claims)
}
