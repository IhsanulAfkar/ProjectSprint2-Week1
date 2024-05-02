package middleware

import (
	"errors"
	"fmt"
	"strings"
	"week1/db"
	"week1/helper/jwt"
	"week1/models"

	"github.com/gin-gonic/gin"
)
func getBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}
func AuthMiddleware(c *gin.Context) {
	token, err := getBearerToken(c.GetHeader("Authorization"))
	if err!= nil {
		c.AbortWithStatusJSON(401, gin.H{
			"message": err.Error()})
		return
	}
	pkId, err := jwt.ParseToken(token)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(401, gin.H{
			"message":err.Error()})
		return
	}
	// find user
	conn := db.CreateConn()
	var user models.User
	result := conn.Raw("SELECT * FROM \"user\" WHERE \"pkId\" = ? LIMIT 1", pkId).Scan(&user)
	if result.RowsAffected ==0{
		c.AbortWithStatusJSON(404, gin.H{
			"message":"user not found"})
			return
	}
	c.Set("userId", pkId)
	c.Next()
}