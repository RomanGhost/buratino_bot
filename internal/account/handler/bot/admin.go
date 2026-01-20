package bot

import "github.com/RomanGhost/buratino_bot.git/internal/account/service"

type AdminHandler struct {
	userService *service.UserService
}

func (h *AdminHandler) BanUser() {
	// TODO
}
