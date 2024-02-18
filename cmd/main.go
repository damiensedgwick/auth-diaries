package main

import (
	"database/sql"
	"html/template"
	"io"

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
	db, err := sql.Open("sqlite3", "auth-diaries.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	e := echo.New()

	e.Renderer = newTemplate()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Static("/static", "static")

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", nil)
	})

	e.POST("/api/v1/auth/login", func(c echo.Context) error {
		email := c.FormValue("email")
		password := c.FormValue("password")

		// Prepare statement
		stmt, err := db.Prepare("SELECT * FROM users WHERE email = ?")
		if err != nil {
			panic(err)
		}
		defer stmt.Close()

		// Execute query with dynamic value
		rows, err := stmt.Query(email)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		var user User

		// Iterate over rows
		for rows.Next() {
			// Update user fields
			if err := rows.Scan(&user.Name, &user.Email, &user.Password); err != nil {
				panic(err)
			}
		}
		if err := rows.Err(); err != nil {
			panic(err)
		}

		if password != user.Password {
			return echo.ErrUnauthorized
		}

		return c.Render(200, "user-card", newPageData(user))
	})

	e.POST("/api/v1/auth/logout", func(c echo.Context) error {
		return c.Render(200, "auth-form", nil)
	})

	e.Logger.Fatal(e.Start(":8080"))
}

type PageData struct {
	User User
}

func newPageData(user User) PageData {
	return PageData{
		User: user,
	}
}

type User struct {
	Name     string
	Email    string
	Password string
}

func newUser() User {
	return User{
		Name:  "Damien Sedgwick",
		Email: "damienksedgwick@gmail.com",
	}
}
