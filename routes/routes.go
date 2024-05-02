package routes

import (
	// "Final-Project/controllers"
	"fmt"
	"net/http"
	"week1/controllers"
	"week1/db"
	"week1/middleware"
	"week1/models"

	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	userController := new(controllers.UserController)
	catController := new(controllers.CatController)
	matchController := new(controllers.MatcherController)
	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})
	v1 := router.Group("/v1")
	{
		v1.GET("/test", func(c *gin.Context) {
			conn := db.CreateConn()
			var cat models.Cat
			result := conn.Raw("SELECT * FROM cat LIMIT 1;").Scan(&cat)
			// var pkId int
			// result := conn.Debug().Raw("SELECT \"pkId\" FROM public.user WHERE email = 'ihsanul2001@gmail.com'").Scan(&pkId)
			if result.RowsAffected == 0 {
				// 
			}
			fmt.Println(cat)
			c.JSON(200, cat)
		})
		user := v1.Group("/user")
		{
			user.POST("/register",userController.SignUp)
			user.POST("/login",userController.SignIn)
		}
		// protected middleware
		v1.Use(middleware.AuthMiddleware)
		cat := v1.Group("/cat")
		{
			cat.POST("/", catController.CreateCat)
			cat.GET("/", catController.GetAllCats)
			cat.PUT("/:catId",catController.UpdateCat)
			cat.DELETE("/:catId",catController.DeleteCat)
			match := cat.Group("/match")
			{
				match.POST("/", matchController.CreateMatch)
				match.GET("/", matchController.GetAllMatches)
				match.POST("/approve", matchController.ApproveMatch)
				match.POST("/reject", matchController.RejectMatch)
				match.DELETE("/:matchId", matchController.DeleteMatch)
			}
		}
		

	}

	return router
}
