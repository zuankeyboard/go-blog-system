package utils

import (
	"errors"
	"go-blog-system/config"
	"go-blog-system/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Claims 自定义JWT声明，包含用户ID、用户名和过期时间
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken 生成JWT Token
func GenerateToken(user *models.User) (string, error) {
	// 设置Token过期时间
	expireTime := time.Now().Add(time.Hour * time.Duration(config.TokenExpire))

	// 构造自定义声明
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), // 过期时间（秒级）
			IssuedAt:  time.Now().Unix(), // 签发时间
			Issuer:    "blog-backend",    // 签发者
		},
	}

	// 创建Token（使用HS256算法签名）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 用密钥签名并生成字符串
	tokenString, err := token.SignedString([]byte(config.JWTSecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析并验证JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	// 解析Token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名算法")
		}
		return []byte(config.JWTSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证Token有效性并提取声明
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("无效的Token")
}
