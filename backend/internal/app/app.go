// Package app wires the application's dependencies together (composition root).
// Keeping construction here keeps main.go thin and the layers decoupled.
package app

import (
	authhandler "github.com/example/ems/internal/auth/handler"
	authrepo "github.com/example/ems/internal/auth/repository"
	authservice "github.com/example/ems/internal/auth/service"
	"github.com/example/ems/internal/config"
	dashhandler "github.com/example/ems/internal/dashboard/handler"
	dashrepo "github.com/example/ems/internal/dashboard/repository"
	dashservice "github.com/example/ems/internal/dashboard/service"
	emphandler "github.com/example/ems/internal/employee/handler"
	emprepo "github.com/example/ems/internal/employee/repository"
	empservice "github.com/example/ems/internal/employee/service"
	"github.com/example/ems/pkg/jwt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Container holds the fully constructed handlers and shared infrastructure used
// by the router.
type Container struct {
	Config           *config.Config
	Logger           *zap.Logger
	JWT              *jwt.Manager
	AuthHandler      *authhandler.Handler
	EmployeeHandler  *emphandler.Handler
	DashboardHandler *dashhandler.Handler
}

// New builds the dependency graph: repositories -> services -> handlers.
func New(cfg *config.Config, db *gorm.DB, log *zap.Logger) *Container {
	jwtMgr := jwt.NewManager(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessTTL,
		cfg.JWT.RefreshTTL,
		cfg.JWT.Issuer,
	)

	// Auth
	aRepo := authrepo.New(db)
	aSvc := authservice.New(aRepo, jwtMgr)
	aHandler := authhandler.New(aSvc)

	// Employee
	eRepo := emprepo.New(db)
	eSvc := empservice.New(eRepo)
	eHandler := emphandler.New(eSvc)

	// Dashboard
	dRepo := dashrepo.New(db)
	dSvc := dashservice.New(dRepo)
	dHandler := dashhandler.New(dSvc)

	return &Container{
		Config:           cfg,
		Logger:           log,
		JWT:              jwtMgr,
		AuthHandler:      aHandler,
		EmployeeHandler:  eHandler,
		DashboardHandler: dHandler,
	}
}
