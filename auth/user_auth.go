package auth

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/hsinnnlu/software_engineering_project/db"
)

var jwtKey = []byte("your_secret_key")

// 2024.11.26 做使用者帳號的驗證
func AuthenticateUser(username, password string) (string, bool, error) {
	var hashedPassword string
	var userPermission string

	query := `SELECT password_hash, user_permission FROM Users WHERE user_id = ?`
	err := db.DB.QueryRow(query, username).Scan(&hashedPassword, &userPermission)
	fmt.Println("dberr: ", err)
	if err == sql.ErrNoRows {
		fmt.Printf("err: %s\n", err)
		return "", false, nil // 用戶不存在
	} else if err != nil {
		fmt.Printf("err: %s\n", err)
		return "", false, err // 一些奇怪的錯誤報錯
	}

	newPassword := HashPassword(password)
	// 比較哈希值
	if hashedPassword != newPassword {
		return "", false, nil // 密碼不匹配
	}
	return userPermission, true, nil
}

func AuthenticateEmail(email string) (string, bool, error) {
    var hashedPassword string

    query := `SELECT password_hash FROM Users WHERE user_email = ?`
    err := db.DB.QueryRow(query, email).Scan(&hashedPassword)
    if err == sql.ErrNoRows {
        // email 不存在
        return "", false, nil
    } else if err != nil {
        // 其他錯誤
        return "", false, fmt.Errorf("查詢錯誤: %w", err)
    }

    // 返回結果
    return hashedPassword, true, nil
}

func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil)) // 以十六進制字串形式返回哈希值
}

func GenerateToken(username, permission string) (string, error) {
	claims := jwt.MapClaims{
		"username":   username,
		"permission": permission,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateToken(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["username"].(string), claims["permission"].(string), nil
	}
	return "", "", err
}
