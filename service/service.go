package service

import (
	"time"

	"Users/helpers"
	"Users/model"
	"Users/repository"
)

type Service interface {
	GetUser() ([]model.User, error)
	GetAlltask(id int) ([]model.TaskRes, error)
	GetTaskById(id int, taskId int) (model.TaskRes, error)
	CreateTask(req model.TaskReq, Id int, ImageURL string) (model.TaskRes, error)
	DeleteTask(Id int) error
	Login(email string, password string) (model.LoginRes, error)
	SaveToken(token string, userId int) error
}

type service struct {
	repository repository.Repository
}

func NewService(repository repository.Repository) *service {
	return &service{repository}
}

func (s *service) GetUser() ([]model.User, error) {
	user, err := s.repository.GetUser()

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) CreateTask(req model.TaskReq, Id int, ImageURL string) (model.TaskRes, error) {
	layout := "2006-01-02 15:04"
	parsedDate, err := time.Parse(layout, req.Date)
	if err != nil {
		return model.TaskRes{}, err
	}

	data, err := s.repository.CreateTask(req, parsedDate, ImageURL, Id)
	if err != nil {
		return model.TaskRes{}, err
	}

	return data, nil
}

func (s *service) DeleteTask(Id int) error {
	err := s.repository.DeleteTask(Id)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Login(email string, password string) (model.LoginRes, error) {
	data, err := s.repository.Login(email)
	if err != nil {
		return model.LoginRes{}, err
	}

	match, err := helpers.ComparePassword(data.Password, password)
	if err != nil {
		return model.LoginRes{}, err
	}
	if !match {
		return model.LoginRes{}, err
	}
	return data, nil
}

func (s *service) SaveToken(token string, userId int) error {
	err := s.repository.SaveToken(token, userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetAlltask(id int) ([]model.TaskRes, error) {
	data, err := s.repository.GetAlltask(id)
	if err != nil {
		return []model.TaskRes{}, err
	}
	return data, nil
}

func (s *service) GetTaskById(id int, taskId int) (model.TaskRes, error) {
	data, err := s.repository.GetTaskById(id, taskId)
	if err != nil {
		return model.TaskRes{}, err
	}

	return data, nil
}