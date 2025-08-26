package web

import (
	"net/http"

	"github.com/RomanGhost/buratino_bot.git/internal/vpn/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

type getNewUserData struct {
	TelegramID int64 `json:"telegram_id"`
	AuthID     uint  `json:"auth_id"`
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var newUserData getNewUserData
	// Привязываем JSON к структуре Data
	if err := c.ShouldBindJSON(&newUserData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createNewUserError := h.userService.AddNewUser(newUserData.TelegramID, newUserData.AuthID)
	if createNewUserError != nil {
		c.JSON(http.StatusFailedDependency, gin.H{"error": createNewUserError.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})

}
