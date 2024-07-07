package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kleankare/user-service/internal/core/domain"
	"github.com/kleankare/user-service/internal/core/ports"
)

type HttpUserHandler interface {
	CreateUser(ctx *gin.Context)
	ReadUsers(ctx *gin.Context)
	ReadUser(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
}

type httpUserHandler struct {
	service ports.UserService
}

func NewHttpUserHandler(service ports.UserService) HttpUserHandler {
	return &httpUserHandler{
		service: service,
	}
}

func (handler *httpUserHandler) CreateUser(ctx *gin.Context) {
	var user domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	newUser, err := handler.service.CreateUser(user.Username, user.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, newUser)
}

func (handler *httpUserHandler) ReadUsers(ctx *gin.Context) {
	users, err := handler.service.ReadUsers()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (handler *httpUserHandler) ReadUser(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := handler.service.ReadUser(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (handler *httpUserHandler) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")

	var user domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	updatedUser, err := handler.service.UpdateUser(id, user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, updatedUser)
}

func (handler *httpUserHandler) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	err := handler.service.DeleteUser(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
