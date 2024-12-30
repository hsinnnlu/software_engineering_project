package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/hsinnnlu/software_engineering_project/auth"
	"github.com/hsinnnlu/software_engineering_project/db"
	"gopkg.in/gomail.v2"
)

// 傳送修改密碼的連結
func GenerateResetpasswdToken(c *gin.Context, email, oldpassword string) (string, error) {
	jwtSecret := "f321a343233594d17697e0c9b83b6cb192a00e8562e4b1738e263c6ac90d3d1d"
	secretKey := jwtSecret + oldpassword
	payload := jwt.MapClaims{
		"email":       email,
		"expiredTime": time.Now().Add(24 * time.Hour).Unix(),
		"type":        "ResetPassword",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(secretKey))
}

func SendResetLink(email string, token string) error {
	resetLink := fmt.Sprintf("http://localhost:8080/webpage/login/reset_password?token=%s", token)
	body := fmt.Sprintf("請點擊以下連結來重設您的密碼: %s", resetLink)

	from := "nmg@cs.thu.edu.tw"
	password := "e04su3su;6"
	subject := "驗證碼"

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("mail.cs.thu.edu.tw", 587, from, password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func VerifylinkToken(URLtoken string) (jwt.MapClaims, error) {
	// 解析 Token 並取得 Email
	secretKeyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		// 需要從 token 中獲得 email，再組合 secretKey
		tokenClaims, _ := token.Claims.(jwt.MapClaims)
		email := tokenClaims["email"].(string)

		oldpassword, _, _ := auth.AuthenticateEmail(email)

		// 組合 secretKey
		jwtSecret := "f321a343233594d17697e0c9b83b6cb192a00e8562e4b1738e263c6ac90d3d1d"
		secretKey := jwtSecret + oldpassword
		return []byte(secretKey), nil
	}

	// 解析 Token
	token, err := jwt.Parse(URLtoken, secretKeyFunc)
	if err != nil {
		return nil, fmt.Errorf("token parsing failed: %v", err)
	}

	// 確認 Token 是否有效
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 驗證 Token 是否過期
		if exp, ok := claims["expiredTime"].(float64); ok {
			expTime := time.Unix(int64(exp), 0)
			if time.Now().After(expTime) {
				return nil, errors.New("token has expired")
			}
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func ResetPassword(claims jwt.MapClaims, password string) error{
	// 從 token 中提取 email
	emailFromClaims, ok := claims["email"].(string)
    if !ok || emailFromClaims == "" {
        return fmt.Errorf("token 驗證失敗或無效的 email 資訊")
    }

	newPassword := auth.HashPassword(password)
	

	// 更新數據庫中的密碼
	query := "UPDATE users SET password_hash = ? WHERE user_email = ?"
	_, err := db.DB.Exec(query, newPassword, emailFromClaims)
	if err != nil {
        return fmt.Errorf("密碼更新失敗: %v", err)
    }

	return nil
}

// func ResetPassword(c *gin.Context) {
// 	fmt.Println("pass ResetPassword")
// 	// 取得用戶提交的密碼重設表單數據
// 	requestBody, err := getResetPasswordRequestBody(c)
// 	log.Println("requestBody: ", requestBody)
// 	log.Println("requestBody_err: ", err)
// 	if err != nil {
// 		c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	newclaims, err := verifyToken(oldpassword, requestBody.Token)
// 	log.Println("newclaims: ", newclaims)
// 	log.Println("newclaims_err: ", err)
// 	if err != nil {
// 		c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
// 			"error": "Token 驗證失敗或已過期",
// 		})
// 		return
// 	}

// 	// 檢查新密碼與確認密碼是否一致
// 	err = validatePasswordMatch(requestBody.Password, requestBody.ConfirmPassword)
// 	log.Println("new_old_password_err: ", err)
// 	if err != nil {
// 		c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	// 從 token 中提取 email
// 	emailFromClaims, err := getEmailFromClaims(newclaims)
// 	log.Println("emailFromClaims: ", emailFromClaims)
// 	log.Println("emailFromClaims_err: ", err)
// 	if err != nil {
// 		c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
// 			"error": "Token 驗證失敗或已過期",
// 		})
// 		return
// 	}

// 	hashedPassword := GetHashedPassword(requestBody.Password)

// 	// 更新數據庫中的密碼
// 	err = updatePasswordInDB(emailFromClaims, hashedPassword)
// 	log.Println("updatedb_err: ", err)
// 	if err != nil {
// 		c.HTML(http.StatusInternalServerError, "reset_password.html", gin.H{
// 			"error": "更新密碼失敗",
// 		})
// 		return
// 	}

// 	c.HTML(http.StatusOK, "reset_password.html", gin.H{
// 		"message": "密碼重設成功，請重新登入",
// 	})
// }

// import (
// 	"errors"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"regexp"
// 	"time"

// 	"github.com/dgrijalva/jwt-go"

// 	"gopkg.in/gomail.v2"

// 	"github.com/gin-gonic/gin"
// 	"github.com/hsinnnlu/software_engineering_project/db"
// )

// var oldpassword string

// // 傳送連結主函式
// func SendLink(c *gin.Context) {
// 	fmt.Println("pass reset password") // 測試路由
// 	// 檢查email是否合法
// 	email, err := preProcessingEmail(c)
// 	log.Println("email: ", email) // 測試email資料

// 	if err != nil {
// 		return
// 	}

// 	// 檢查email是否存在資料庫
// 	err = checkEmail(c, email)
// 	log.Println("email_db_err: ", err)
// 	if err != nil {
// 		return
// 	}

// 	// 找舊密碼生成token
// 	oldpassword, _ = findPassword(c, email)
// 	log.Println("oldpassword: ", oldpassword)
// 	if oldpassword == "" {
// 		return
// 	}
// 	token, err := generateToken(email, oldpassword)
// 	log.Println("token: ", err)
// 	log.Println("token_err: ", err)
// 	if err != nil {
// 		c.HTML(http.StatusBadRequest, "login.html", gin.H{
// 			"error": "Token生成失敗",
// 		})
// 		return
// 	}

// 	// 送信件連結
// 	err = sendResetLink(email, token)
// 	log.Println("send_mail_err: ", err)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "發送郵件失敗"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"success": true, "message": "重設密碼的連結已發送"})
// }

// // 渲染更改密碼頁面主函式
// func ResetPasswordPage(c *gin.Context) {
// 	fmt.Println("pass ResetPasswordPage")
// 	// 取得URL中的連結
// 	URLtoken := GetURLtoken(c)
// 	log.Println("url_token: ", URLtoken)
// 	if URLtoken == "" {
// 		return
// 	}

// 	// 驗證token
// 	claims, err := verifyToken(oldpassword, URLtoken)
// 	log.Println("claims: ", claims)
// 	log.Println("claims_err: ", err)
// 	if err != nil {
// 		c.HTML(http.StatusBadRequest, "error.html", gin.H{
// 			"error": "Token 驗證失敗或已過期",
// 		})
// 		return
// 	}

// 	// 顯示重設密碼頁面
// 	Page(c, URLtoken, claims)
// }

// // 更改密碼主函式
// func ResetPassword(c *gin.Context) {
// 	fmt.Println("pass ResetPassword")
// 	// 取得用戶提交的密碼重設表單數據
// 	requestBody, err := getResetPasswordRequestBody(c)
// 	log.Println("requestBody: ", requestBody)
// 	log.Println("requestBody_err: ", err)
// 	if err != nil {
// 		c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	newclaims, err := verifyToken(oldpassword, requestBody.Token)
// 	log.Println("newclaims: ", newclaims)
// 	log.Println("newclaims_err: ", err)
// 	if err != nil {
// 		c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
// 			"error": "Token 驗證失敗或已過期",
// 		})
// 		return
// 	}

// 	// 檢查新密碼與確認密碼是否一致
// 	err = validatePasswordMatch(requestBody.Password, requestBody.ConfirmPassword)
// 	log.Println("new_old_password_err: ", err)
// 	if err != nil {
// 		c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	// 從 token 中提取 email
// 	emailFromClaims, err := getEmailFromClaims(newclaims)
// 	log.Println("emailFromClaims: ", emailFromClaims)
// 	log.Println("emailFromClaims_err: ", err)
// 	if err != nil {
// 		c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
// 			"error": "Token 驗證失敗或已過期",
// 		})
// 		return
// 	}

// 	hashedPassword := GetHashedPassword(requestBody.Password)

// 	// 更新數據庫中的密碼
// 	err = updatePasswordInDB(emailFromClaims, hashedPassword)
// 	log.Println("updatedb_err: ", err)
// 	if err != nil {
// 		c.HTML(http.StatusInternalServerError, "reset_password.html", gin.H{
// 			"error": "更新密碼失敗",
// 		})
// 		return
// 	}

// 	c.HTML(http.StatusOK, "reset_password.html", gin.H{
// 		"message": "密碼重設成功，請重新登入",
// 	})
// }

// func preProcessingEmail(c *gin.Context) (email string, err error) {
// 	var requestData struct {
// 		UserEmail string `json:"user_email"` // 使用與 JSON 中相同的 key
// 	}

// 	if err := c.BindJSON(&requestData); err != nil {
// 		c.HTML(http.StatusBadRequest, "login.html", gin.H{
// 			"error": "無效的請求格式",
// 		})
// 		return "", err
// 	}

// 	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
// 	if requestData.UserEmail != "" {
// 		log.Println("email_input: ", requestData.UserEmail)
// 		if emailRegex.MatchString(requestData.UserEmail) {
// 			email = requestData.UserEmail
// 		} else {
// 			c.HTML(http.StatusBadRequest, "login.html", gin.H{
// 				"error": "請輸入有效的email格式",
// 			})
// 			return "", err
// 		}
// 	} else {
// 		c.HTML(http.StatusBadRequest, "login.html", gin.H{
// 			"error": "請輸入email",
// 		})
// 		return "", err
// 	}

// 	return email, nil
// }

// func checkEmail(c *gin.Context, email string) error {
// 	_, err := db.GetUserByEmail(email)
// 	if err != nil {
// 		c.HTML(http.StatusBadRequest, "login.html", gin.H{
// 			"error": errors.New("email does not exist"),
// 		})
// 		return err
// 	}
// 	return nil
// }

// func findPassword(c *gin.Context, email string) (string, error) {
// 	user, err := db.GetpasswordByEmail(email) // 呼叫 GetPasswordByEmail 獲取用戶資料

// 	if err != nil {
// 		c.HTML(http.StatusBadRequest, "login.html", gin.H{
// 			"error": "資料庫中沒有找到密碼",
// 		})
// 		return "", err
// 	}

// 	// 返回用戶的密碼哈希
// 	return user.Password_hash, nil
// }

// func generateToken(email, password string) (string, error) {
// 	jwtSecret := "f321a343233594d17697e0c9b83b6cb192a00e8562e4b1738e263c6ac90d3d1d"
// 	secretKey := jwtSecret + password
// 	payload := jwt.MapClaims{
// 		"email":       email,
// 		"expiredTime": time.Now().Add(5 * time.Minute).Unix(),
// 		"type":        "ResetPassword",
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
// 	return token.SignedString([]byte(secretKey))
// }

// func sendResetLink(email string, token string) error {
// 	resetLink := fmt.Sprintf("http://localhost:8080/webpage/login/reset_password?token=%s", token)
// 	body := fmt.Sprintf("請點擊以下連結來重設您的密碼: %s", resetLink)

// 	from := "nmg@cs.thu.edu.tw"
// 	password := "e04su3su;6"
// 	subject := "驗證碼"

// 	m := gomail.NewMessage()
// 	m.SetHeader("From", from)
// 	m.SetHeader("To", email)
// 	m.SetHeader("Subject", subject)
// 	m.SetBody("text/plain", body)

// 	d := gomail.NewDialer("mail.cs.thu.edu.tw", 587, from, password)

// 	if err := d.DialAndSend(m); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func GetURLtoken(c *gin.Context) string {
// 	// user點link後，會傳到這做驗證
// 	token := c.Query("token")

// 	if token == "" {
// 		c.HTML(http.StatusBadRequest, "login.html", gin.H{
// 			"error": "無效的重設連結",
// 		})
// 		return ""
// 	} else {
// 		return token
// 	}
// }

// func verifyToken(password string, URLtoken string) (jwt.MapClaims, error) {
// 	secretKeyFunc := func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, errors.New("unexpected signing method")
// 		}
// 		jwtSecret := "f321a343233594d17697e0c9b83b6cb192a00e8562e4b1738e263c6ac90d3d1d"
// 		secretKey := jwtSecret + password
// 		return []byte(secretKey), nil
// 	}

// 	token, err := jwt.Parse(URLtoken, secretKeyFunc)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// token.Claims.(jwt.MapClaims) 將 token 的 claims 轉換為 jwt.MapClaims 型別
// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		// 如果 token 已經過期
// 		if exp, ok := claims["expiredTime"].(float64); ok {
// 			expTime := time.Unix(int64(exp), 0)
// 			if time.Now().After(expTime) {
// 				return nil, errors.New("token has expired")
// 			}
// 		}
// 		return claims, nil
// 	}

// 	return nil, errors.New("invalid token")
// }

// func Page(c *gin.Context, URLtoken string, claims jwt.MapClaims) {
// 	// Token有效，顯示更改密碼頁面
// 	c.HTML(http.StatusOK, "reset_password.html", gin.H{
// 		"token": URLtoken,        // 在HTML表單中以隱藏字段的形式保存token
// 		"email": claims["email"], // 顯示Email或用於後續驗證
// 	})
// }

// func getResetPasswordRequestBody(c *gin.Context) (struct {
// 	Token           string `form:"token"`
// 	email           string `form:"email"`
// 	Password        string `form:"password"`
// 	ConfirmPassword string `form:"confirm-password"`
// }, error) {
// 	var requestBody struct {
// 		Token           string `form:"token"`
// 		email           string `form:"email"`
// 		Password        string `form:"password"`
// 		ConfirmPassword string `form:"confirm-password"`
// 	}

// 	if err := c.Bind(&requestBody); err != nil {
// 		return requestBody, errors.New("無效的請求")
// 	}

// 	return requestBody, nil
// }

// func getEmailFromClaims(claims jwt.MapClaims) (string, error) {
// 	email, ok := claims["email"].(string)
// 	if !ok {
// 		return "", errors.New("從 Token 中提取 Email 失敗")
// 	}
// 	return email, nil
// }

// func validatePasswordMatch(password, confirmPassword string) error {
// 	if password != confirmPassword {
// 		return errors.New("密碼和確認密碼不一致")
// 	}
// 	return validatePassword(password)
// }

// func validatePassword(password string) error {
// 	// 密碼長度至少8個字元
// 	if len(password) < 8 {
// 		log.Println("Password:", password)
// 		return errors.New("密碼長度必須至少8個字元")
// 	}

// 	// 密碼必須包含至少一個大寫字母
// 	uppercasePattern := `[A-Z]`
// 	matched, _ := regexp.MatchString(uppercasePattern, password)
// 	if !matched {
// 		return errors.New("密碼必須包含至少一個大寫字母")
// 	}

// 	// 密碼不可包含空白字元或特殊符號
// 	illegalPattern := `[^\w!?]`
// 	matched, _ = regexp.MatchString(illegalPattern, password)
// 	if matched {
// 		return errors.New("密碼不可包含空白字元或特殊符號")
// 	}

// 	return nil
// }

// func updatePasswordInDB(email, hashedPassword string) error {
// 	DB = db.DB
// 	err := db.UpdatePasswordByEmail(DB, email, hashedPassword)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
