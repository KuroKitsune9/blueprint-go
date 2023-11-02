package controller

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"Users/helpers"
	"Users/model"
	"Users/service"
)

type Controller struct {
	service service.Service
}

func NewController(service service.Service) *Controller {
	return &Controller{service}
}

func (c *Controller) GetAlltaskController(ctx echo.Context) error {
	Claims := helpers.ClaimToken(ctx)
	id := Claims.ID

	task, err := c.service.GetAlltask(id)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "succes get",
		"data":    task,
	})
}

func (c *Controller) GetTaskById(ctx echo.Context) error {
	Claims := helpers.ClaimToken(ctx)
	id := Claims.ID
	taskId := ctx.Param("id")
	taskidtrue, err := strconv.Atoi(taskId)
	if err != nil {
		return err
	}

	task, err := c.service.GetTaskById(id, taskidtrue)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "succes get by id",
		"data":    task,
	})
}

func (c *Controller) CreateTasksController(ctx echo.Context) error {
	var req model.TaskReq

	err := ctx.Bind(&req)
	if err != nil {
		return err
	}
	// Menerima file gambar dari form dengan nama "image"
	image, err := ctx.FormFile("image")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Tidak dapat memproses file gambar"})
	}

	// Buka file yang diunggah
	src, err := image.Open()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal membuka file gambar"})
	}
	defer src.Close()

	// Lokasi penyimpanan file gambar lokal
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal membuat direktori penyimpanan"})
	}

	// Generate nama file unik
	dstPath := filepath.Join(uploadDir, image.Filename)

	// Membuka file tujuan untuk penyimpanan
	dst, err := os.Create(dstPath)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal membuat file gambar"})
	}
	defer dst.Close()

	// Salin isi file dari file asal ke file tujuan
	if _, err = io.Copy(dst, src); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal menyalin file gambar"})
	}

	// Membuat URL ke gambar yang diunggah
	imageURL := "http://localhost:8080/uploads/" + image.Filename

	Claims := helpers.ClaimToken(ctx)
	Id := Claims.ID

	// validate .......
	task, err := c.service.CreateTask(req, Id, imageURL)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"task":    task,
		"message": "successfully created task",
	})
}

func (c *Controller) DeleteTasksController(ctx echo.Context) error {
	Id := ctx.Param("id")
	IdAsli, err := strconv.Atoi(Id)
	if err != nil {
		return err
	}

	err = c.service.DeleteTask(IdAsli)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "Task deleted successfully",
	})
}

func (c *Controller) Login(ctx echo.Context) error {
	var req model.LoginRequest
	err := ctx.Bind(&req)
	if err != nil {
		return err
	}

	data, err := c.service.Login(req.Email, req.Password)
	if err != nil {
		return err
	}

	var (
		jwtToken  *jwt.Token
		secretKey = []byte("secret")
	)

	jwtClaims := &model.Claims{
		ID:    data.Id,
		Email: data.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	jwtToken = jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	token, err := jwtToken.SignedString(secretKey)
	if err != nil {
		return err
	}

	err = c.service.SaveToken(token, int(data.Id))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"token":   token,
		"message": "success Login",
		"data":    data,
	})
}
