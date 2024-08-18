package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"vk_tarantool_project/internal/domain"
)

// ReadData handler for get data from storage by Keys
func (h *Handler) ReadData(c echo.Context) error {
	var keys domain.DataKeys

	if err := c.Bind(&keys); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	data, err := h.service.ReadData(
		c.Request().Context(),
		keys,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, data)
}
