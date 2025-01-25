// utils/jwt.go
package utils

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type Claims struct {
	UserID string `json:"userId"`
	// jwt.StandardClaimsを埋め込むことで、有効期限や発行日時などの標準的なクレームを使うことができる。
	// クレームとは、トークンに含まれる情報のこと。
	jwt.StandardClaims
}

func GenerateJWT(userID string, secret string) (string, error) {
	// トークンの有効期限を24時間に設定
	expirationTime := time.Now().Add(24 * time.Hour)

	// トークンに埋め込むクレームを作成
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			// Unix()メソッドは、time.Time型の値をUnix時間に変換する。
			// Unix時間は、1970年1月1日からの経過秒数を表す。
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	// NewWithClaims：指定したアルゴリズムとクレームを使って、新しいトークンを作成する。
	// SigningMethodHS256：HS256アルゴリズムを使って署名を行う。
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// SignedString：指定した文字列を使ってトークンを署名し、文字列に変換する。
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWT(tokenStr, secret string) (*Claims, error) {
	claims := &Claims{}
	// ParseWithClaims：指定したトークン文字列をパースし、指定したクレームに格納する。
	// パースとは、トークン文字列を構造化して、トークンの内容を取り出すこと。
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}
