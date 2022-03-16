package model

import (
	"sync"

	"github.com/HideInBush7/go_im/pkg/mysql"
	"github.com/jmoiron/sqlx"
)

// 表user结构
type User struct {
	Uid      int64  `db:"uid"`
	Username string `db:"username"`
	Password string `db:"password"`
}

func NewUserModel() *UserModel {
	return &UserModel{
		db:    mysql.GetInstance(),
		table: "user",
	}
}

type UserModel struct {
	db    *sqlx.DB
	table string
	sync.Pool
}

func (u *UserModel) Insert(user *User) (err error) {
	query := `INSERT INTO user (username,password) VALUES (?, ?)`
	_, err = u.db.Exec(query, user.Username, user.Password)
	return
}

func (u *UserModel) GetByUid(uid int64) (*User, error) {
	var result = &User{}
	query := `SELECT uid, username, password FROM user WHERE uid=?`
	err := u.db.Get(result, query, uid)
	return result, err
}

func (u *UserModel) GetByUsername(username string) (*User, error) {
	var result = &User{}
	query := `SELECT uid, username, password FROM user WHERE username=?`
	err := u.db.Get(result, query, username)
	return result, err
}

func (u UserModel) Table() string {
	return u.table
}
