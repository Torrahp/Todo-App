package middleware

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret"

func makeToken(t *testing.T, secret string, method jwt.SigningMethod, claims jwt.MapClaims) string {
	t.Helper() // ถ้า error จะชี้บรรทัดที่เรียก makeToken แทนบรรทัดในนี้
	s, err := jwt.NewWithClaims(method, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return s
}

func TestAuthRequired(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	gin.SetMode(gin.TestMode)

	future := time.Now().Add(time.Hour).Unix()
	past := time.Now().Add(-time.Hour).Unix()

	valid := makeToken(t, testSecret, jwt.SigningMethodHS256, jwt.MapClaims{"user_id": 42, "exp": future})
	expired := makeToken(t, testSecret, jwt.SigningMethodHS256, jwt.MapClaims{"user_id": 42, "exp": past})
	wrongSecret := makeToken(t, "other-secret", jwt.SigningMethodHS256, jwt.MapClaims{"user_id": 42, "exp": future})
	wrongAlg := makeToken(t, testSecret, jwt.SigningMethodHS512, jwt.MapClaims{"user_id": 42, "exp": future})
	noUserID := makeToken(t, testSecret, jwt.SigningMethodHS256, jwt.MapClaims{"exp": future})

	tests := []struct {
		name       string
		header     string
		wantStatus int
		wantUserID int
	}{
		{"no header", "", 401, 0},
		{"not a bearer token", "Token " + valid, 401, 0},
		{"garbage token", "Bearer abc.def.ghi", 401, 0},
		{"wrong secret", "Bearer " + wrongSecret, 401, 0},
		{"expired", "Bearer " + expired, 401, 0},
		{"wrong algorithm", "Bearer " + wrongAlg, 401, 0},
		{"missing user_id", "Bearer " + noUserID, 401, 0},
		{"valid", "Bearer " + valid, 200, 42},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var gotUserID int
			r := gin.New()
			r.GET("/", AuthRequired(), func(c *gin.Context) {
				gotUserID = c.GetInt("userID")
				c.Status(200)
			})

			req := httptest.NewRequest("GET", "/", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tc.wantStatus)
			}
			if gotUserID != tc.wantUserID {
				t.Errorf("userID = %d, want %d", gotUserID, tc.wantUserID)
			}
		})
	}
}
