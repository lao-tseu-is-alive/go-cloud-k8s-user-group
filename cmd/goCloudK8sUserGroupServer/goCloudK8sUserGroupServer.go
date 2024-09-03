package main

import (
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/config"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/database"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/goHttpEcho"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/golog"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/metadata"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/tools"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/users"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/version"
	"runtime"
	"strings"
)

const (
	defaultPort                = 8080
	defaultDBPort              = 5432
	defaultDBIp                = "127.0.0.1"
	defaultDBSslMode           = "prefer"
	defaultWebRootDir          = "goCloudK8sUserGroupFront/dist/"
	defaultSqlDbMigrationsPath = "db/migrations"
	defaultAdminUser           = "goadmin"
	defaultAdminEmail          = "goadmin@yourdomain.org"
	defaultAdminId             = 960901
	charsetUTF8                = "charset=UTF-8"
	MIMEAppJSON                = "application/json"
	MIMEHtml                   = "text/html"
	MIMEAppJSONCharsetUTF8     = MIMEAppJSON + "; " + charsetUTF8
	MIMEHtmlCharsetUTF8        = MIMEHtml + "; " + charsetUTF8
)

var db database.DB

// content holds our static web server content.
//
//go:embed goCloudK8sUserGroupFront/dist/*
var content embed.FS

// sqlMigrations holds our db migrations sql files using https://github.com/golang-migrate/migrate
// in the line above you SHOULD have the same path  as const defaultSqlDbMigrationsPath
//
//go:embed db/migrations/*.sql
var sqlMigrations embed.FS

func main() {
	// l := log.New(os.Stdout, fmt.Sprintf("%s ", version.APP), log.Ldate|log.Ltime|log.Lshortfile)
	prefix := fmt.Sprintf("%s ", version.APP)
	l, _ := golog.NewLogger("dev", golog.DebugLevel, prefix)
	l.Info("🚀🚀 Starting App:'%s', ver:%s, from: %s", version.APP, version.VERSION, version.REPOSITORY)

	dbDsn := config.GetPgDbDsnUrlFromEnvOrPanic(defaultDBIp, defaultDBPort, tools.ToSnakeCase(version.APP), version.AppSnake, defaultDBSslMode)
	db, err := database.GetInstance("pgx", dbDsn, runtime.NumCPU(), l)
	if err != nil {
		l.Fatal("💥💥 error doing database.GetInstance(pgx ...) error: %v", err)
	}
	defer db.Close()

	dbVersion, err := db.GetVersion()
	if err != nil {
		l.Fatal("💥💥 error doing dbConn.GetVersion() error: %v", err)
	}
	l.Info("connected to db version : %s", dbVersion)

	// checking metadata information
	metadataService := metadata.Service{Log: l, Db: db}
	metadataService.CreateMetadataTableOrFail()
	found, ver := metadataService.GetServiceVersionOrFail(version.APP)
	if found {
		l.Info("service %s was found in metadata with version: %s", version.APP, ver)
	} else {
		l.Info("service %s was not found in metadata", version.APP)
	}
	metadataService.SetServiceVersionOrFail(version.APP, version.VERSION)

	// begin section go-migrate db migration with embed files in go program
	// https://github.com/golang-migrate/migrate
	d, err := iofs.New(sqlMigrations, defaultSqlDbMigrationsPath)
	if err != nil {
		l.Fatal("💥💥 error doing iofs.New for db migrations  error: %v\n", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, strings.Replace(dbDsn, "postgres", "pgx5", 1))
	if err != nil {
		l.Fatal("💥💥 error doing migrate.NewWithSourceInstance(iofs, dbURL:%s)  error: %v\n", dbDsn, err)
	}

	err = m.Up()
	if err != nil {
		//if err == m.
		if !errors.Is(err, migrate.ErrNoChange) {
			l.Fatal("💥💥 error doing migrate.Up error: %v\n", err)
		}
	}
	// end section go-migrate db migration with embed files in go program

	myVersionReader := goHttpEcho.NewSimpleVersionReader(version.APP, version.VERSION, version.REPOSITORY)
	// Create a new JWT checker
	myJwt := goHttpEcho.NewJwtChecker(
		config.GetJwtSecretFromEnvOrPanic(),
		config.GetJwtIssuerFromEnvOrPanic(),
		version.APP,
		config.GetJwtDurationFromEnvOrPanic(60),
		l)
	// Create a new Authenticator with a simple admin user
	myAuthenticator := goHttpEcho.NewSimpleAdminAuthenticator(&goHttpEcho.UserInfo{
		UserId:     config.GetAdminIdFromEnvOrPanic(defaultAdminId),
		ExternalId: config.GetAdminExternalIdFromEnvOrPanic(9999999),
		Name:       "NewSimpleAdminAuthenticator_Admin",
		Email:      config.GetAdminEmailFromEnvOrPanic(defaultAdminEmail),
		Login:      config.GetAdminUserFromEnvOrPanic(defaultAdminUser),
		IsAdmin:    false,
	},
		config.GetAdminPasswordFromEnvOrPanic(),
		myJwt)

	server := goHttpEcho.CreateNewServerFromEnvOrFail(
		defaultPort,
		"0.0.0.0", // defaultServerIp,
		&goHttpEcho.Config{
			ListenAddress: "",
			Authenticator: myAuthenticator,
			JwtCheck:      myJwt,
			VersionReader: myVersionReader,
			Logger:        l,
			WebRootDir:    defaultWebRootDir,
			Content:       content,
			RestrictedUrl: "/api/v1",
		},
	)

	userStore := users.GetStorageInstanceOrPanic("pgx", db, l)

	userService := users.Service{
		Logger: l,
		DbConn: db,
		Store:  userStore,
		Server: server,
	}

	e := server.GetEcho()
	e.GET("/login", userService.GetLogin)
	e.POST("/login", userService.LoginUser)
	e.GET("/resetpassword", userService.GetResetPasswordEmail)
	e.POST("/resetpassword", userService.SendResetPassword)
	e.GET("/resetpassword/:resetPasswordToken", userService.GetResetPasswordToken)
	e.POST("/resetpassword/:resetPasswordToken", userService.ResetPassword)
	r := server.GetRestrictedGroup()
	users.RegisterHandlers(r, &userService) // register all openapi declared routes
	//some other restricted routes
	r.GET("/status", userService.GetStatus)
	r.GET("/users/maxid", userService.GetMaxId)

	err = server.StartServer()
	if err != nil {
		l.Fatal("💥💥 ERROR: 'calling echo.StartServer() got error: %v'\n", err)
	}

}
