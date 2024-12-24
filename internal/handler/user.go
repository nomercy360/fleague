package handler

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/terrors"
	"net/http"
)

func getUserID(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*contract.JWTClaims)
	return claims.UID
}

func (h *Handler) ListMatches(c echo.Context) error {
	resp, err := h.service.GetActiveMatches(c.Request().Context(), getUserID(c))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetUserInfo(c echo.Context) error {
	username := c.Param("username")
	resp, err := h.service.GetUserInfo(c.Request().Context(), username)
	if err != nil {
		return terrors.InternalServer(err, "cannot get user info")
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) ListMyReferrals(c echo.Context) error {
	resp, err := h.service.GetUserReferrals(c.Request().Context(), getUserID(c))
	if err != nil {
		return terrors.InternalServer(err, "cannot get referrals")
	}

	return c.JSON(http.StatusOK, resp)
}
