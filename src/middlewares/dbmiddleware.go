package middlewares

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DatabaseMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("DB", db)
		c.Next()
	}
}

func GetDB(c *gin.Context) *gorm.DB {
	db, exists := c.Get("DB")
	if !exists {
		return nil
	}
	return db.(*gorm.DB)
}
