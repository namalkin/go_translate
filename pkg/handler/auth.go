package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/namalkin/go_translate/pkg/tables"
)

func (h *Handler) signUp(c *gin.Context) {
	var input tables.SignUpInput

	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user := tables.User{
		Name:     input.Name,
		Username: input.Username,
		Password: input.Password,
	}

	id, err := h.services.Authorisation.CreateUser(user)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) signIn(c *gin.Context) {
	var input tables.SignInInput

	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorisation.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
