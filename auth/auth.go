package auth

import (
	"fmt"
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hsinnnlu/software_engineering_project/db"
	"github.com/hsinnnlu/software_engineering_project/models"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			c.Redirect(http.StatusFound, "/login")
		} else {
			c.Set("user_id", userID)
			c.Next()
		}
	}
}

// GetCurrentUser retrieves the current user based on the session data
func GetCurrentUser(c *gin.Context) (*models.User, error) {
    DB := db.DB
	session := sessions.Default(c)

    // Get user_id from session
    user_id, ok := session.Get("user_id").(string)
    if !ok || user_id == "" {
        return nil, errors.New("user not logged in or invalid session")
    }

    // Fetch user from the database by user_id
    user, err := db.GetUserById(DB, user_id)
    if err != nil {
        return nil, fmt.Errorf("error fetching user: %v", err)
    }

    return user, nil
}

