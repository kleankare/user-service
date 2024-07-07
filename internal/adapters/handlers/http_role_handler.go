package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kleankare/user-service/internal/core/domain"
	"github.com/kleankare/user-service/internal/core/ports"
)

type HttpRoleHandler interface {
	CreateRole(ctx *gin.Context)
	ReadRoles(ctx *gin.Context)
	ReadRole(ctx *gin.Context)
	UpdateRole(ctx *gin.Context)
	DeleteRole(ctx *gin.Context)
}

type httpRoleHandler struct {
	service ports.RoleService
}

func NewHttpRoleHandler(service ports.RoleService) HttpRoleHandler {
	return &httpRoleHandler{
		service: service,
	}
}

func (handler *httpRoleHandler) CreateRole(ctx *gin.Context) {
	var role domain.Role
	if err := ctx.ShouldBindJSON(&role); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	newRole, err := handler.service.CreateRole(role.Name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, newRole)
}

func (handler *httpRoleHandler) ReadRoles(ctx *gin.Context) {
	roles, err := handler.service.ReadRoles()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, roles)
}

func (handler *httpRoleHandler) ReadRole(ctx *gin.Context) {
	id := ctx.Param("id")

	role, err := handler.service.ReadRole(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, role)
}

func (handler *httpRoleHandler) UpdateRole(ctx *gin.Context) {
	id := ctx.Param("id")

	var role domain.Role
	if err := ctx.ShouldBindJSON(&role); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	updatedRole, err := handler.service.UpdateRole(id, role)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, updatedRole)
}

func (handler *httpRoleHandler) DeleteRole(ctx *gin.Context) {
	id := ctx.Param("id")

	err := handler.service.DeleteRole(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Role deleted"})
}
