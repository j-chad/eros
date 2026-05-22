package admin

import (
	"backend/internal/service"
)

type Handler struct {
	authService   *service.AuthService
	favourService *service.FavourService
	fileService   *service.FileService
	graphService  *service.GraphService
	pushService   *service.PushService
}

func NewHandler(
	authService *service.AuthService,
	favourService *service.FavourService,
	graphService *service.GraphService,
	fileService *service.FileService,
	pushService *service.PushService,
) *Handler {
	return &Handler{
		authService:   authService,
		favourService: favourService,
		fileService:   fileService,
		graphService:  graphService,
		pushService:   pushService,
	}
}
