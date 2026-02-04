package service

import "backend/internal/repository"

type GraphService struct {
	repo repository.Repository
}

func NewGraphService(repo repository.Repository) *GraphService {
	return &GraphService{repo: repo}
}
