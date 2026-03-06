package service

import (
	"backend/internal/repository"
	"backend/internal/repository/storage"
)

type FileService struct {
	repo  repository.Repository
	files storage.FileStore
}

func NewFileService(repo repository.Repository, files storage.FileStore) *FileService {
	return &FileService{repo: repo, files: files}
}
