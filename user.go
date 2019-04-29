package main

type user struct {
	ID       int    `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u user) String() string {
	return "{" + "ID=" + string(u.ID) + ";UserName=" + u.UserName + ";Email=" +
		u.Email + ";Password=" + u.Password + "}"
}
