package model

import (
	"database/sql"

	"github.com/HideInBush7/go_im/pkg/mysql"
)

var ErrNoRows = sql.ErrNoRows

// 表user结构
type User struct {
	Uid      int64  `db:"uid"`
	Username string `db:"username"`
	Password string `db:"password"`
}

func (u User) Table() string {
	return "user"
}

// 新增用户
func InsertUser(user *User) (int64, error) {
	query := `INSERT INTO user (username,password) VALUES (?, ?)`
	res, err := mysql.GetInstance().Exec(query, user.Username, user.Password)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// 根据uid获取用户
func GetUserByUid(uid int64) (*User, error) {
	var result = &User{}
	query := `SELECT uid, username, password FROM user WHERE uid=?`
	err := mysql.GetInstance().Get(result, query, uid)
	return result, err
}

// 根据用户名返回用户
func GetUserByUsername(username string) (*User, error) {
	var result = &User{}
	query := `SELECT uid, username, password FROM user WHERE username=?`
	err := mysql.GetInstance().Get(result, query, username)
	return result, err
}
