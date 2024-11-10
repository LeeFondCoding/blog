package mysql

import (
	"blog/models"
	"blog/pkg/snowflake"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

const (
	secret = "liiy"
)

// 对密码进行加密
func encryptPassword(data []byte) (result string) {
	hash_algo := md5.New()
	hash_algo.Write([]byte(secret))
	return hex.EncodeToString(hash_algo.Sum(data))
}

// 检测用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return errors.New(ErrorUserExit)
	}
	return
}

// 向数据库中插入用户数据
func InsertUser(user models.User) (err error) {
	user.PassWord = encryptPassword([]byte(user.PassWord))
	sqlstr := `insert into user(user_id,username,password,email,gender) values(?,?,?,?,?)`
	_, err = db.Exec(sqlstr, user.UserID, user.UserName, user.PassWord, user.Email, user.Gender)
	return
}

// 用户注册
func Register(user *models.User) (err error) {
	sqlStr := "select count(user_id) from user where username = ?"
	var count int64
	err = db.Get(&count, sqlStr, user.UserName)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	if count > 0 {
		return errors.New(ErrorUserExit)
	}

	userID, err := snowflake.GetID()
	if err != nil {
		return ErrorGenIDFailed
	}

	password := encryptPassword([]byte(user.PassWord))
	sqlStr = "insert into user(user_id, username, password) values (?,?,?)"
	_, err = db.Exec(sqlStr, userID, user.UserName, password)
	return
}

// 用户登录
func Login(user *models.User) (err error) {
	originPassword := user.PassWord
	sqlStr := "select user_id, username, password from user where username = ?"
	err = db.Get(user, sqlStr, user.UserName)

	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		return errors.New(ErrorUserNotExit)
	}

	password := encryptPassword([]byte(originPassword))
	if user.PassWord != password {
		return errors.New(ErrorPasswordWrong)
	}
	return nil
}

// 根据用户id查询用户信息
func GetUserByID(id uint64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id, username from user where user_id = ?`
	err = db.Get(user, sqlStr, id)
	return
}
