package main

import (
	"html/template"
	"io"

	"github.com/damiensedgwick/auth-diaries/database"
	"github.com/damiensedgwick/auth-diaries/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

type Template struct {
	tmpl *template.Template
}

func newTemplate() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("views/*.html")),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

func main() {
	database.NewDBConn()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Static("/static", "static")

	e.Renderer = newTemplate()

	e.GET("/", routes.HomeRoute)
	e.POST("/api/v1/auth/login", routes.LoginRoute)
	e.POST("/api/v1/auth/logout", routes.LogoutRoute)

	e.Logger.Fatal(e.Start(":8080"))
}
