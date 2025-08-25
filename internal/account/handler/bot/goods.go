package bot

import "github.com/RomanGhost/buratino_bot.git/internal/account/service"

type GoodsHandler struct {
	goodsService *service.GoodsService
}

func NewGoodsHadler() *GoodsHandler {
	return &GoodsHandler{}
}

func (h *GoodsHandler) GetAll() {

}
