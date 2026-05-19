package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// Register のバリデーションエラーを httptest で検証する例。
// authService は呼ばれないため nil でよい。
func TestAuthController_Register_ValidationError(t *testing.T) {
	ctrl := NewAuthController(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"email":"not-an-email","password":"short"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	ctrl.Register(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body = %s", w.Code, http.StatusBadRequest, w.Body.String())
	}
}

// ActivateUser で token 未指定のとき 400 になることを検証する例。
func TestAuthController_ActivateUser_MissingToken(t *testing.T) {
	ctrl := NewAuthController(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/activate", nil)

	ctrl.ActivateUser(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body = %s", w.Code, http.StatusBadRequest, w.Body.String())
	}
}
