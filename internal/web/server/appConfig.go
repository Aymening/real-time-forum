package server

import (
	"forum/internal/api"
	"forum/internal/db"
	"forum/internal/web/handlers"
)

type App struct {
	Handlers handlers.Handler
	Api      api.Api
	DB       *db.Database
}

func NewApp() (*App, error) {
	database, err := db.NewDatabase()
	if err != nil {
		return nil, err
	}

	if err := database.ExecuteAllTableInDataBase(); err != nil {
		return nil, err
	}

	return &App{
		Handlers: handlers.Handler{
			DB: database,
		},
		Api: api.Api{
			Users: make(map[string]*db.User),
			Comments : make(map[int][]*db.Comment),
		},
		DB: database,
	}, nil
}

func InitApp() (*App, error) {
	app, err := NewApp()
	if err != nil {
		return nil, err
	}

	if err := app.DB.GetPostsFromDB(app.Api.Users, &app.Api.Posts); err != nil {
		return nil, err
	}
	
	if err := app.DB.GetAllCommentsFromDataBase(app.Api.Comments); err != nil {
		return nil, err
	}

	app.Handlers.Api = &app.Api

	return app, nil
}
