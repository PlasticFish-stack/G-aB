package controll

import (
	"fmt"
	"project/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserName string
	UserId   uint
	jwt.RegisteredClaims
}

// Valid implements jwt.Claims.
func (c *Claims) Valid() error {
	panic("unimplemented")
}

var AccessSecret = []byte(getSecretKey("access"))
var RefreshSecret = []byte(getSecretKey("refresh"))

func getSecretKey(key string) string {
	sk, _ := utils.GetValue(key)
	if sk == "" {
		sk = utils.RandString(15, utils.ALL)
		utils.SetValue(key, sk)
	}
	return sk
}

func GenerateJwt(username string, userid uint) (accessToken string, refreshToken string, expTime time.Time, err error) {
	accessExpiration := time.Now().Add(30 * time.Minute)
	expTime = time.Now().Add(30 * time.Minute)
	refreshExpiration := time.Now().Add(24 * time.Hour)
	accessClaims := &Claims{
		UserName: username,
		UserId:   userid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiration),
		},
	}
	refreshClaims := &Claims{
		UserName: username,
		UserId:   userid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiration),
		},
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(AccessSecret)
	if err != nil {
		return
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(RefreshSecret)
	return
}

func Refresh(resultToken string) (tokens JwtToken, err error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(resultToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return RefreshSecret, nil
	})
	if err != nil || !token.Valid {
		return tokens, fmt.Errorf("解析token出错,请检查")
	}
	tokens.AccessToken, tokens.RefreshToken, tokens.ExpTime, err = GenerateJwt(claims.UserName, claims.UserId)
	return
}
