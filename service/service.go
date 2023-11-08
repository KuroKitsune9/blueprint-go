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
	UpdateTask(req model.TaskReq, ImageURL string, Id int) (model.TaskRes, error)
	DeleteTask(Id int) error
	BulkDeleteTask(taskIds []int, Id int) error
	Login(email string, password string) (model.LoginRes, error)
	Regis(email string, password string) (model.UserRegisRespon, error)
	Logout(reqToken string) error
	SaveToken(token string, userId int) error
	CountTasks(id int, keywoard string, parsedDate time.Time) (int, error)
	SearchTasks(id int, keywoard string, parsedDate time.Time, limit, offset int) ([]model.TaskRes, error)
	CountTask(Id int) (model.Count, error)
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

func (s *service) GetAlltask(id int) ([]model.TaskRes, error) {
	data, err := s.repository.GetAlltask(id)
	if err != nil {
		return []model.TaskRes{}, err
	}
	return data, nil
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

func (s *service) UpdateTask(req model.TaskReq, ImageURL string, Id int) (model.TaskRes, error) {
	layout := "2006-01-02 15:04"
	parsedDate, err := time.Parse(layout, req.Date)
	if err != nil {
		return model.TaskRes{}, err
	}

	data, err := s.repository.UpdateTasks(req, parsedDate, ImageURL, Id)
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

func (s *service) BulkDeleteTask(taskIds []int, Id int) error {
	err := s.repository.BulkDeleteTask(taskIds, Id)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Regis(email string, password string) (model.UserRegisRespon, error) {

	HasPassword, err := helpers.HashPassword(password)
	if err != nil {
		return model.UserRegisRespon{}, err
	}

	data, err := s.repository.Regis(email, HasPassword)
	if err != nil {
		return model.UserRegisRespon{}, err
	}
	return data, nil
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

func (s *service) GetTaskById(id int, taskId int) (model.TaskRes, error) {
	data, err := s.repository.GetTaskById(id, taskId)
	if err != nil {
		return model.TaskRes{}, err
	}

	return data, nil
}

func (s *service) Logout(reqToken string) error {
	err := s.repository.Logout(reqToken)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) CountTasks(id int, keywoard string, parsedDate time.Time) (int, error) {
	data, err := s.repository.CountTasks(id, keywoard, parsedDate)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (s *service) SearchTasks(id int, keywoard string, parsedDate time.Time, limit, offset int) ([]model.TaskRes, error) {
	data, err := s.repository.SearchTasks(id, keywoard, parsedDate, limit, offset)
	if err != nil {
		return []model.TaskRes{}, err
	}

	return data, nil
}

func (s *service) CountTask(Id int) (model.Count, error) {
	data, err := s.repository.CountTask(Id)
	if err != nil {
		return model.Count{}, err
	}

	return data, nil
}
