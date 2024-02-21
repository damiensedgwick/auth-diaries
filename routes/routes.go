package routes

import (
	"fmt"

	"github.com/damiensedgwick/auth-diaries/database"
	"github.com/damiensedgwick/auth-diaries/model/page"
	"github.com/damiensedgwick/auth-diaries/model/user"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func HomeRoute(c echo.Context) error {
	sess, _ := session.Get("session", c)

	if sess.Values["user"] != nil {
		fmt.Println(sess.Values["user"])
	}

	return c.Render(200, "index", nil)
}

func LoginRoute(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	// Prepare statement
	stmt, err := database.DBConn.Prepare("SELECT * FROM users WHERE email = ?")
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

	var u user.User

	// Iterate over rows
	for rows.Next() {
		// Update user fields
		if err := rows.Scan(&u.Name, &u.Email, &u.Password); err != nil {
			panic(err)
		}
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return echo.ErrUnauthorized
	}

	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["user"] = u.Email
	sess.Save(c.Request(), c.Response())

	return c.Render(200, "user-card", page.NewPage(u))
}

func LogoutRoute(c echo.Context) error {
	return c.Render(200, "auth-form", nil)
}
