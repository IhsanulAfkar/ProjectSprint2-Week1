package controllers

import (
	"fmt"
	"week1/db"
	"week1/forms"
	"week1/helper"
	"week1/helper/hash"
	"week1/helper/jwt"
	"week1/models"

	"github.com/gin-gonic/gin"
)
type UserController struct{}


func (h UserController)SignUp(c *gin.Context) {

	conn := db.CreateConn()
	var userForm forms.UserForms
	if err := c.ShouldBindJSON(&userForm); err != nil {
		c.JSON(400, gin.H{
			"message":err.Error()})
		return
    }
	// check email format
	if !helper.IsValidEmail(userForm.Email) {
		c.JSON(400,gin.H{
			"message":"incorrect email address"})
		return
	}
	// check user exist with email
	var user models.User
	result := conn.Raw("SELECT * FROM public.user WHERE email = ? LIMIT 1", userForm.Email).Scan(&user)
	if result.RowsAffected > 0{
		fmt.Println(user)
		c.JSON(409, gin.H{"message": "user with email already exist"})
		return
	}
	// check username
	if len(userForm.Name) < 5 || len(userForm.Name) > 50 {
	
		c.JSON(400, gin.H{
			"message":"name cannot below 5 nor exceed 15 characters"})
		return
	}
	// check password
	if len(userForm.Password) < 5 || len(userForm.Password) >15 {
		
		c.JSON(400, gin.H{
			"message":"password cannot below 5 nor exceed 15 characters"})
		return
	}
	hashedPass, err:= hash.HashPassword(userForm.Password)
	if err != nil{
		c.JSON(500, gin.H{"Message":err.Error()})
		return
	}
	
	res := conn.Exec("INSERT INTO \"user\" (email, name, password) VALUES (?,?,?)",userForm.Email, userForm.Name, hashedPass)
	if res.Error != nil {
		c.JSON(500, gin.H{"message":"failed to create user", "error": res.Error.Error()})
		return
	}
	res2 := conn.Raw("SELECT * FROM \"user\" WHERE email = ? LIMIT 1;", userForm.Email).Scan(&user)
	if res2.RowsAffected == 0 {
		c.JSON(500, gin.H{
			"message":"server error"});
		return
	}
	// create access token
	accessToken := jwt.SignJWT(user)
	// res3 := conn.Exec("UPDATE public.user SET \"accessToken\" = ? WHERE email = ?", accessToken, user.Email)
	// if res3.RowsAffected == 0 {
	// 	fmt.Println(res3.Error.Error())
	// 	c.JSON(500, gin.H{
	// 		"message":"server error"})
	// 		return
	// }
	data := map[string]string{
		"email":userForm.Email,
		"name":userForm.Name,
		"accessToken":accessToken}

	c.JSON(201, gin.H{
		"message":"user created successfully",
		"data": data});
}

func (h UserController)SignIn(c *gin.Context){
	conn:=db.CreateConn()
	var loginForm forms.LoginForms
	if err:= c.ShouldBindJSON(&loginForm); err != nil{
		c.JSON(400, gin.H{
			"message":err.Error()})
		return
	}
	if loginForm.Email == "" || loginForm.Password == ""{
		c.JSON(400, gin.H{"message":"bad request"})
		return
	}
	if !helper.IsValidEmail(loginForm.Email) {
		c.JSON(400, gin.H{"message":"invalid email format"})
		return
	}
	if len(loginForm.Password) < 5 || len(loginForm.Password) >15 {
		c.JSON(400, gin.H{
			"message":"password cannot below 5 nor exceed 15 characters"})
		return
	}
	var user models.User
	result := conn.Debug().Raw("SELECT * FROM public.user WHERE email = ? LIMIT 1", loginForm.Email).Scan(&user)
	// workaround
	// var pkId int
	// result2 := conn.Debug().Raw("SELECT \"pkId\" FROM public.user WHERE email = ? LIMIT 1", loginForm.Email).Scan(&pkId)
	// user.PkId = pkId
	//
	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"message": "user not found"})
		return
	}
	fmt.Println(loginForm.Password, user.Password)
	if !hash.CheckPasswordHash(loginForm.Password, user.Password) {
		c.JSON(400, gin.H{"message": "invalid password"})
		return
	}
	// generate accessToken
	// fmt.Println("user ", user)
	accessToken := jwt.SignJWT(user)
	// res := conn.Exec("UPDATE public.user SET \"accessToken\" = ? WHERE email = ?", accessToken, user.Email)
	// if res.RowsAffected == 0 {
	// 	fmt.Println(res.Error.Error())
	// 	c.JSON(500, gin.H{
	// 		"message":"server error"})
	// 		return
	// }
	data := map[string]string{
		"email":user.Email,
		"name":user.Name,
		"accessToken":accessToken}

	c.JSON(200, gin.H{
		"message":"User logged successfully",
		"data": data});
}