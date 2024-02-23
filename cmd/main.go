package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	tmpl *template.Template
}

func newTemplate() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("templates/*.html")),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	e.Static("/static", "static")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	store := sessions.NewCookieStore([]byte("secret"))
	e.Use(session.Middleware(store))

	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		sess, _ := session.Get("session", c)

		fmt.Println(sess)

		if sess.Values["user"] != nil {
			var user User

			err := json.Unmarshal(sess.Values["user"].([]byte), &user)
			if err != nil {
				fmt.Println("error unmarshalling user value")
				return err
			}

			return c.Render(200, "index", newPageData(user))
		}

		return c.Render(200, "index", nil)
	})

	e.GET("/auth/sign-in", func(c echo.Context) error {
		return c.Render(200, "auth-form", nil)
	})

	e.POST("/auth/sign-in", func(c echo.Context) error {
		user := newUser()

		sess, _ := session.Get("session", c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
		}
		userBytes, err := json.Marshal(user)
		if err != nil {
			fmt.Println("error marshalling user value")
			return err
		}
		sess.Values["user"] = userBytes
		err = sess.Save(c.Request(), c.Response())
		if err != nil {
			fmt.Println("error saving session")
			return err
		}

		return c.Render(200, "index", newPageData(user))
	})

	e.POST("/auth/sign-out", func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		sess.Options.MaxAge = -1
		err := sess.Save(c.Request(), c.Response())
		if err != nil {
			fmt.Println("error saving session")
			return err
		}

		return c.Render(200, "index", nil)
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