package handlers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"vk_tarantool_project/internal/domain"
)

func (h *Handler) Login(c echo.Context) error {

	var userInfo domain.UserInfo
	if err := c.Bind(&userInfo); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if userInfo.Username == "" || userInfo.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, domain.ErrUsernamePasswordRequired)
	}

	token, err := h.service.Login(c.Request().Context(), userInfo)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) || errors.Is(err, domain.ErrPasswordMismatch) {
			return echo.NewHTTPError(http.StatusUnauthorized, domain.ErrInvalidCredentials)
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})

}
