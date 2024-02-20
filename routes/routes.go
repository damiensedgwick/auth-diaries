package routes

import (
	"github.com/damiensedgwick/auth-diaries/database"
	"github.com/damiensedgwick/auth-diaries/handler/page"
	"github.com/damiensedgwick/auth-diaries/handler/user"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func HomeRoute(c echo.Context) error {
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

	return c.Render(200, "user-card", page.NewPage(u))
}

func LogoutRoute(c echo.Context) error {
	return c.Render(200, "auth-form", nil)
}
