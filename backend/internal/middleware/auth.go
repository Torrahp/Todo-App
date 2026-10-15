package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")     //ดึงส่วนAuthorization
		if !strings.HasPrefix(header, "Bearer ") { //เช็คว่ารูปแบบเป็น Bearer หรือไม่
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ") // ตัดคำว่า Bearer ออกเหลือแค่ token

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil //จะใช้Key อะไรตรวจสอบJWT
		}, jwt.WithValidMethods([]string{"HS256"})) //บังคับว่าจะใช้Algorนี่
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		//Type assertion เพื่อยืนยันว่า claims เป็น jwt.MapClaims ซึ่งเป็นชนิดข้อมูลที่ใช้เก็บข้อมูลใน JWT
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}
		userID, ok := claims["user_id"].(float64) //Go ถอดรหัสJson number ลง interface{} เป็น float64
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}

		c.Set("userID", int(userID)) // set userID in context for later use
		c.Next()                     // ส่งต่อไปยัง handler ถัดไป
	}
}
