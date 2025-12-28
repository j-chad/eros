package service

import "backend/internal/repository"

type AdminService struct {
	repo repository.Repository
}

func NewAdminService(repo repository.Repository) *AdminService {
	return &AdminService{repo: repo}
}
