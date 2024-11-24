package logic

import (
	"blog/dao/mysql"
	"blog/models"
	"blog/pkg/snowflake"
	"blog/pkg/jwt"
)

// 注册逻辑
func SignUp(register *models.RegisterForm) (err error) {
	if err := mysql.CheckUserExist(register.UserName); err != nil {
		return err
	}

	userID, err := snowflake.GetID()
	if err != nil {
		return mysql.ErrorGenIDFailed
	}

	return mysql.InsertUser(models.User{
		UserID:   userID,
		UserName: register.UserName,
		PassWord: register.Password,
		Email: register.Email,
		Gender: register.Gender,
	})
}

// 登录逻辑
func Login(login *models.LoginForm) (*models.User, error) {
	user := &models.User{
		UserName: login.UserName,
		PassWord: login.Password,
	}

	if err := mysql.Login(user); err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := jwt.GenToken(user.UserID, user.UserName)
	if err != nil {
		return nil, err
	}

	user.AccessToken = accessToken
	user.RefreshToken = refreshToken
	return user, nil
}