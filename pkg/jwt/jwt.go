package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

type MyClaims struct {
	UserID   uint64 `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.StandardClaims
}

const (
	TokenExpiredDuration = time.Hour * 24

	AccessTokenExpireDuration  = time.Hour * 24
	RefreshTokenExpireDuration = time.Hour * 24 * 7
)

var mySecret = []byte("夏天夏天悄悄过去")

func keyFunc(_ *jwt.Token) (i any, err error) {
	return mySecret, nil
}

func GenToken(userID uint64, userName string) (accessToken, refreshToken string, err error) {
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, MyClaims{
		UserID:   userID,
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AccessTokenExpireDuration).Unix(),
			Issuer:    "lichun",
		},
	}).SignedString(mySecret)
	if err != nil {
		zap.L().Error("jwt.NewWithClaims failed", zap.Error(err))
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, MyClaims{
		UserID:   userID,
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(RefreshTokenExpireDuration).Unix(),
			Issuer:    "lichun",
		},
	}).SignedString(mySecret)
	if err != nil {
		zap.L().Error("jwt.NewWithClaims failed", zap.Error(err))
	}

	return
}

func ParseToken(tokenString string) (claims *MyClaims, err error) {
	var token *jwt.Token
	claims = new(MyClaims)
	token, err = jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		zap.L().Error("ParseToken failed", zap.Error(err))
		return
	}
	if !token.Valid {
		err = errors.New("invalid token")
	}
	return
}

func RefreshToken(accessToken, refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	// check token
	if _, err = jwt.Parse(refreshToken, keyFunc); err != nil {
		return
	}

	var claim MyClaims
	_, err = jwt.ParseWithClaims(accessToken, &claim, keyFunc)
	v, _ := err.(*jwt.ValidationError)

	if v.Errors == jwt.ValidationErrorExpired {
		return GenToken(claim.UserID, claim.UserName)
	}
	return
}
