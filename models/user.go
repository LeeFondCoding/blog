package models

import (
	"encoding/json"
	"errors"
)

type User struct {
	UserID       uint64 `db:"user_id, string" json:"user_id"`
	UserName     string `db:"user_name" json:"user_name"`
	PassWord     string `db:"pass_word" json:"pass_word"`
	Email        string `db:"email" json:"email"`
	Gender       int    `db:"gender" json:"gender"`
	AccessToken  string
	RefreshToken string
}

// 将json数据解码到User结构体中
// 检查是否缺少字段user name, pass word
func (u *User) UnMarshalJSON(data []byte) (err error) {
	required := struct {
		UserID   uint64 `db:"user_id, string" json:"user_id"`
		UserName string `db:"user_name" json:"user_name"`
		PassWord string `db:"pass_word" json:"pass_word"`
		Email    string `db:"email" json:"email"`
		Gender   int    `db:"gender" json:"gender"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return
	} else if len(required.UserName) == 0 {
		err = errors.New("缺少必填字段user name")
	} else if len(required.PassWord) == 0 {
		err = errors.New("缺少必填字段pass word")
	} else {
		u.UserName = required.UserName
		u.PassWord = required.PassWord
		u.Email = required.Email
		u.Gender = required.Gender
	}
	return
}

type RegisterForm struct {
	UserName        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Gender          int    `json:"gender" binding:"oneof=0 1 2"` // 性别 0:未知 1:男 2:女
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type LoginForm struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// json数据解码到RegisterForm结构体中
// 检查是否缺少字段username, password, email
func (r *RegisterForm) UnMarshalJSON(data []byte) (err error) {
	required := struct {
		UserName        string `json:"username"`
		Email           string `json:"email"`
		Gender          int    `json:"gender"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}{}

	err = json.Unmarshal(data, &required)

	if err != nil {
		return
	} else if len(required.UserName) == 0 {
		err = errors.New("缺少必填字段username")
	} else if len(required.Password) == 0 {
		err = errors.New("缺少必填字段password")
	} else if len(required.Email) == 0 {
		err = errors.New("缺少必填字段email")
	} else if required.Password != required.ConfirmPassword {
		err = errors.New("两次密码不一致")
	} else {
		r.UserName = required.UserName
		r.Email = required.Email
		r.Gender = required.Gender
		r.Password = required.Password
		r.ConfirmPassword = required.ConfirmPassword
	}
	return
}

type VoteDataForm struct {
	PostID string `json:"post_id" binding:"required"`
	Direction int8 `json:"direction,string" binding:"oneof=1 0 -1"`
}

// 将json解码到VoteDataForm结构体中
//检查是否缺少字段post_id, direction
func (v *VoteDataForm) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		PostID    string `json:"post_id"`
		Direction int8   `json:"direction"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return
	} else if len(required.PostID) == 0 {
		err = errors.New("缺少必填字段post_id")
	} else if required.Direction == 0 {
		err = errors.New("缺少必填字段direction")
	} else {
		v.PostID = required.PostID
		v.Direction = required.Direction
	}
	return
}
