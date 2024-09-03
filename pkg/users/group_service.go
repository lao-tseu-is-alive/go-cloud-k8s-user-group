package users

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/goHttpEcho"
	"net/http"
)

// GroupCreate will store a new Group in the store
func (s Service) GroupCreate(ctx echo.Context) error {
	handlerName := "GroupCreate"
	goHttpEcho.TraceRequest(handlerName, ctx.Request(), s.Logger)
	// get the current user from JWT TOKEN
	claims := s.Server.JwtCheck.GetJwtCustomClaimsFromContext(ctx)
	currentUserId := int32(claims.User.UserId)
	s.Logger.Info("in %s : currentUserId: %d", handlerName, currentUserId)
	// IF USER IS NOT ADMIN RETURN 401 Unauthorized
	if !s.Store.IsUserAdmin(currentUserId) {
		return echo.NewHTTPError(http.StatusUnauthorized, "current user has no admin privilege")
	}
	newGroup := &Group{
		Id:      0,
		Creator: int32(currentUserId),
	}
	if err := ctx.Bind(newGroup); err != nil {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("GroupCreate has invalid format [%v]", err))
	}
	if len(newGroup.Name) < 1 {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprint("GroupCreate name cannot be empty"))
	}
	if len(newGroup.Name) < 5 {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprint("GroupCreate name minLength is 5"))
	}
	s.Logger.Info("# GroupCreate() newGroup : %#v\n", newGroup)
	groupCreated, err := s.Store.CreateGroup(*newGroup)
	if err != nil {
		msg := fmt.Sprintf("GroupCreate had an error saving group:%#v, err:%#v", *newGroup, err)
		s.Logger.Info(msg)
		return ctx.JSON(http.StatusBadRequest, msg)
	}
	s.Logger.Info("# GroupCreate() Group %#v\n", groupCreated)
	return ctx.JSON(http.StatusCreated, groupCreated)

}

// GroupGet will retrieve the Group in the store and return then
// to test it with curl you can try :
// curl -s -H "Content-Type: application/json" -H "Authorization: Bearer $token" 'http://localhost:8888/api/users' |jq
func (s Service) GroupGet(ctx echo.Context, id int32) error {
	handlerName := "GroupGet"
	goHttpEcho.TraceRequest(handlerName, ctx.Request(), s.Logger)
	claims := s.Server.JwtCheck.GetJwtCustomClaimsFromContext(ctx)
	currentUserId := int32(claims.User.UserId)
	s.Logger.Info("in %s : currentUserId: %d", handlerName, currentUserId)
	// IF USER IS NOT ADMIN RETURN 401 Unauthorized
	if !s.Store.IsUserAdmin(currentUserId) {
		return echo.NewHTTPError(http.StatusUnauthorized, "current user has no admin privilege")
	}
	group, err := s.Store.GetGroup(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("problem retrieving group :%v", err))
	}
	return ctx.JSON(http.StatusOK, group)
}

// GroupList will retrieve all Groups in the store and return then
func (s Service) GroupList(ctx echo.Context, params GroupListParams) error {
	handlerName := "GroupList"
	goHttpEcho.TraceRequest(handlerName, ctx.Request(), s.Logger)
	// get the current user from JWT TOKEN
	claims := s.Server.JwtCheck.GetJwtCustomClaimsFromContext(ctx)
	currentUserId := int32(claims.User.UserId)
	s.Logger.Info("in %s : currentUserId: %d", handlerName, currentUserId)
	list, err := s.Store.ListGroup(0, 100)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("problem calling store.ListGroup :%v", err))
	}
	return ctx.JSON(http.StatusOK, list)
}

// GroupDelete will remove the given id entry from the store, and if not present will return 400 Bad Request
// curl -v -XDELETE -H "Content-Type: application/json" -H "Authorization: Bearer $token" 'http://localhost:8888/api/groups/3' ->  204 No Content if present and delete it
func (s Service) GroupDelete(ctx echo.Context, id int32) error {
	handlerName := "GroupDelete"
	goHttpEcho.TraceRequest(handlerName, ctx.Request(), s.Logger)
	claims := s.Server.JwtCheck.GetJwtCustomClaimsFromContext(ctx)
	currentUserId := int32(claims.User.UserId)
	s.Logger.Info("in %s : currentUserId: %d", handlerName, currentUserId)
	if !s.Store.IsUserAdmin(currentUserId) {
		return echo.NewHTTPError(http.StatusUnauthorized, "current user has no admin privilege")
	}
	err := s.Store.DeleteGroup(id)
	if err != nil {
		msg := fmt.Sprintf("GroupDelete(%d) got an error: %#v ", id, err)
		s.Logger.Info(msg)
		return echo.NewHTTPError(http.StatusInternalServerError, msg)
	}
	return ctx.NoContent(http.StatusNoContent)

}

// GroupUpdate will store the modified information in the store for the given id
func (s Service) GroupUpdate(ctx echo.Context, id int32) error {
	handlerName := "GroupUpdate"
	goHttpEcho.TraceRequest(handlerName, ctx.Request(), s.Logger)
	claims := s.Server.JwtCheck.GetJwtCustomClaimsFromContext(ctx)
	currentUserId := int32(claims.User.UserId)
	s.Logger.Info("in %s : currentUserId: %d", handlerName, currentUserId)
	if !s.Store.IsUserAdmin(currentUserId) {
		return echo.NewHTTPError(http.StatusUnauthorized, "current user has no admin privilege")
	}
	g := new(Group)
	if err := ctx.Bind(g); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("GroupUpdate has invalid format [%v]", err))
	}
	g.LastModificationUser = &currentUserId
	if len(g.Name) < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprint("GroupUpdate name cannot be empty"))
	}
	if len(g.Name) < 5 {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprint("GroupUpdate name minLength is 5"))
	}

	updatedGroup, err := s.Store.UpdateGroup(id, *g)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("GroupUpdate got problem updating group :%v", err))
	}
	return ctx.JSON(http.StatusOK, updatedGroup)
}
