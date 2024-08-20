package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"vk_tarantool_project/internal/domain"
)

// WriteData handle for get and write data to storage
func (h *Handler) WriteData(c echo.Context) error {

	var data domain.Data
	statusResponseBody := map[string]string{"status": "success"}

	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if len(data.Data) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, domain.ErrDataRequired)
	}

	if err := h.service.WriteData(c.Request().Context(), data); err != nil {
		statusResponseBody["status"] = "error"
		c.Logger().Errorf("Error saving data: %v", err)
		return c.JSON(http.StatusInternalServerError, statusResponseBody)
	}

	return c.JSON(http.StatusCreated, statusResponseBody)
}
