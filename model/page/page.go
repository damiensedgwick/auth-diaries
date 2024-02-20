package page

import "github.com/damiensedgwick/auth-diaries/model/user"

type PageData struct {
	User user.User
}

func NewPage(user user.User) PageData {
	return PageData{
		User: user,
	}
}
