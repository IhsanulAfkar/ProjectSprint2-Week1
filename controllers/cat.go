package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"week1/db"
	"week1/forms"
	"week1/helper"
	"week1/models"

	"github.com/gin-gonic/gin"
)

type CatController struct{}

func (h CatController) CreateCat(c *gin.Context){
	userId := c.MustGet("userId")
	var catForm forms.CreateCatForms
	listError := make(map[string]interface{})
	conn := db.CreateConn()
	var createdCat models.Cat
	if err := c.ShouldBindJSON(&catForm); err != nil {
		c.JSON(400, gin.H{
			"message":err.Error()})
			return
		}
		// validations
	if len(catForm.Name) < 1 || len(catForm.Name) > 30 {
		listError["name"] = "name cannot below 5 nor exceed 30"
		// c.JSON(400, gin.H{
		// 	"message":"name cannot below 5 nor exceed 30"})
		// 	return
	}
	
	if !helper.Includes(catForm.Race, models.CatRaces[:]){
		listError["race"] = "invalid cat race"
	}
	if !helper.Includes(catForm.Sex, models.Sex[:]){
		listError["sex"] = "invalid sex"
	}
	if catForm.AgeInMonth < 1 || catForm.AgeInMonth > 120082 {
		listError["age"] = "age in month must between 1-120082"
	}
	if len(catForm.Description) < 1 || len(catForm.Description) > 200 {
		listError["description"] = "description length must between 1-200 characters"
	}
	if len(catForm.ImageUrls) == 0 {
		listError["imageUrls"] = "imageUrls cannot be empty"
	}
	for idx, value := range catForm.ImageUrls {
		index := strconv.Itoa(idx)
		if value == "" {
			listError["imageUrl["+index+"]"] = "item cannot be empty"
		} else {
			// if not an url
			if !helper.IsURL(value){
				listError["imageUrl["+index+"]"] = "item is not an url"
			}
			// if it is, pass
		}
	}
	
	// have error
	if len(listError) > 0 {
		c.JSON(400, gin.H{
			"message":"invalid input",
			"errors":listError})
		return
	}

	
	// createdAt := time.Now().Format("2006-01-02T15:04:05Z")
	query:= "INSERT INTO cat (\"userId\", name, race, sex, \"ageInMonth\", description, \"imageUrls\") VALUES (?,?,?,?,?,?,?) RETURNING *;"
	res := conn.Raw(query, userId, catForm.Name, catForm.Race, catForm.Sex, catForm.AgeInMonth, catForm.Description, helper.ConvertToPGList(catForm.ImageUrls)).Scan(&createdCat)
	fmt.Println(res.Statement.SQL.String())
	if res.Error != nil {	
		c.JSON(500, gin.H{
			"message":"failed to create cat",
			"error":res.Error.Error()})
		return
	}
	
	fmt.Println("cat controller: ", createdCat)
	// yay
	data := map[string]string{
		"id":createdCat.Id,
		"createdAt":helper.FormatToIso860(createdCat.CreatedAt)}
	c.JSON(201, gin.H{
		"message":"success",
		"data": data})
}

func (h CatController) GetAllCats(c *gin.Context){
	userId := c.MustGet("userId")
	paramId := c.Query("id")
	if !helper.IsUUID(paramId){
		paramId =""
	}
	limit, errLimit := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if errLimit != nil || limit < 0 {
		limit = 5
	}
	offset, errOffset := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if errOffset != nil || offset < 0 {
		offset = 0
	}
	race := c.Query("race")
	sex := c.Query("sex")
	hasMatched := c.Query("hasMatched")
	if !helper.Includes(hasMatched, []string{"true", "false"}) {
		hasMatched = ""
	}
	ageInMonth := c.Query("ageInMonth")
	owned := c.Query("owned")
	search := c.Query("search")
	
	fmt.Println(paramId, limit,offset,race,sex,hasMatched, ageInMonth,owned, search)

	// validate query
	if !helper.Includes(race, models.CatRaces[:]){
		race = ""
	}
	if !helper.Includes(sex, models.Sex[:]){
		sex = ""
	}
	// Build the SQL query string
	sqlQuery := "SELECT \"pkId\", id, name, race, sex, \"ageInMonth\",\"imageUrls\", description, \"hasMatched\", \"created_at\" FROM cat"

	// Prepare parameters slice
	var args []interface{}
	// queryParams := make(map[string]interface{})
	var queryParams []string
	where := false
	// Add conditions based on parameters
	// WHERE
	if paramId != "" {
		where = true
		queryParams = append(queryParams, " id = ?") 
		args = append(args, paramId)
	}
	if race != ""{
		where =true
		queryParams = append(queryParams, " race = ?")
		args = append(args, race)
	}
	if sex != "" {
		where = true 
		queryParams = append(queryParams, " sex = ?")
		args = append(args, sex)
	}
	if hasMatched != "" {
		where = true
		queryParams = append(queryParams, " \"hasMatched\" = " + hasMatched) 
	}
	if ageInMonth != "" {
		// if no lt nor gt
		if !(ageInMonth[0] == '>' || ageInMonth[0] == '<') {
			if ageInMonth[0] == '='{
				ageInMonth = ageInMonth[1:]
			}
			_, err := strconv.Atoi(ageInMonth)
			if err == nil {
				// not an int
				where = true
				queryParams = append(queryParams, " \"ageInMonth\" = " + ageInMonth)
			}
		} else {
			operand := ageInMonth[0]
			rest := ageInMonth[1:]
			_, err := strconv.Atoi(rest)
			if err == nil && (operand =='>'|| operand=='<'){
				where = true
				queryParams = append(queryParams, " \"ageInMonth\" "+ string(operand) + " " + rest)
			}
		}
	}
	if owned != ""{
		where = true
		if owned == "true" {
			queryParams = append(queryParams, " \"userId\" = ?")
			args = append(args, userId)
		} else {
			queryParams = append(queryParams, " \"userId\" != ?")
			args = append(args, userId)
		}
	}

	if search != ""{
		where=true
		queryParams = append(queryParams, " name LIKE '%"+search+"%'")
	}
	
	if where {
		allQuery := strings.Join(queryParams, " AND")
		sqlQuery += " WHERE" + allQuery 
		fmt.Println(sqlQuery)
		fmt.Println(args...)
	}
	sqlQuery += " LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset)
	conn := db.CreateConn()
	// userId := c.MustGet("userId")
	
	rows,err := conn.Raw(sqlQuery, args...).Rows()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500, gin.H{
			"message":"server error"})
		return
	}
	// var arrCat []models.Cat
	arrCat := make([]models.Cat,0)
	var cat models.Cat
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&cat.PkId,&cat.Id, &cat.Name, &cat.Race, &cat.Sex, &cat.AgeInMonth, &cat.ImageUrls, &cat.Description, &cat.HasMatched, &cat.CreatedAt)
		if err != nil {
			fmt.Println(err.Error())
		}
		arrCat = append(arrCat, cat)
	  
		// do something
	  }
	cleanArr := make([]models.GetCat, 0)
	for _, oldCat := range arrCat {
		newCat :=models.GetCat{
			Id:oldCat.Id,
			Name: oldCat.Name,
			Race: oldCat.Race,
			Sex: oldCat.Sex,
			AgeInMonth: oldCat.AgeInMonth,
			Description: oldCat.Description,
			HasMatched: oldCat.HasMatched,
			ImageUrls: oldCat.ImageUrls,
			CreatedAt: helper.FormatToIso860(oldCat.CreatedAt)}
		cleanArr = append(cleanArr, newCat)
	}
	c.JSON(200, gin.H{
		"message":"success",
		"data":cleanArr})
}

func (h CatController)UpdateCat(c *gin.Context){
	catId := c.Param("catId")
	userId := c.MustGet("userId")
	var catForm forms.CreateCatForms
	listError := make(map[string]interface{})
	if err := c.ShouldBindJSON(&catForm); err != nil {
		c.JSON(400, gin.H{
			"message":err.Error()})
		return
	}

	// validations
	if !helper.IsUUID(catId) {
		listError["id"] = "id is not uuid format"
	}
	if len(catForm.Name) < 1 || len(catForm.Name) > 30 {
		listError["name"] = "name cannot below 5 nor exceed 30"
		// c.JSON(400, gin.H{
		// 	"message":"name cannot below 5 nor exceed 30"})
		// 	return
	}
	
	if !helper.Includes(catForm.Race, models.CatRaces[:]){
		listError["race"] = "invalid cat race"
	}
	if !helper.Includes(catForm.Sex, models.Sex[:]){
		listError["sex"] = "invalid sex"
	}
	if catForm.AgeInMonth < 1 || catForm.AgeInMonth > 120082 {
		listError["age"] = "age in month must between 1-120082"
	}
	if len(catForm.Description) < 1 || len(catForm.Description) > 200 {
		listError["description"] = "description length must between 1-200 characters"
	}
	if len(catForm.ImageUrls) == 0 {
		listError["imageUrls"] = "imageUrls cannot be empty"
	}
	for idx, value := range catForm.ImageUrls {
		index := strconv.Itoa(idx)
		if value == "" {
			listError["imageUrl["+index+"]"] = "item cannot be empty"
		} else {
			// if not an url
			if !helper.IsURL(value){
				listError["imageUrl["+index+"]"] = "item is not an url"
			}
			// if it is, pass
		}
	}
	// have error
	if len(listError) > 0 {
		c.JSON(400, gin.H{
			"message":"invalid input",
			"errors":listError})
		return
	}
	conn := db.CreateConn()
	var cat models.Cat	
	result := conn.Raw("SELECT * FROM cat WHERE id = ? AND \"userId\" = ?", catId, userId).Scan(&cat)
	if result.RowsAffected == 0{
		c.JSON(404, gin.H{
			"message":"id not found"})
			return
		}
	result = conn.Exec("SELECT * FROM matcher WHERE \"userCatId\" = ? OR \"matchCatId\" = ?", cat.PkId, cat.PkId) 
	
	if cat.Sex != catForm.Sex && result.RowsAffected > 0 {
		c.JSON(400, gin.H{
			"message":"cannot edit cat's sex when already matched"})
			return
	}
	result2 := conn.Raw("UPDATE cat SET name = ?, race = ?, sex = ?, \"ageInMonth\" = ?, description = ?, \"imageUrls\" = ?, updated_at = ? WHERE id = ? RETURNING * ", catForm.Name, catForm.Race, catForm.Sex, catForm.AgeInMonth, catForm.Description, helper.ConvertToPGList(catForm.ImageUrls), time.Now(), catId).Scan(&cat)
	if result2.RowsAffected ==0 || result2.Error!= nil {
		fmt.Println(result2.Error.Error())
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	data:= map[string]string{
		"id":cat.Id,
		"updatedAt" : cat.UpdatedAt}
	c.JSON(200,gin.H{
		"message":"success update cat",
		"data": data})
}

func (h CatController)DeleteCat(c *gin.Context){
	catId := c.Param("catId")
	userId := c.MustGet("userId")
	conn := db.CreateConn()
	if !helper.IsUUID(catId) {
		c.JSON(400, gin.H{"message":"id is not uuid"})
		return
	}
	fmt.Println(catId, " ", userId)
	checkCat := conn.Exec("SELECT * FROM cat WHERE id = ? AND \"userId\" = ? LIMIT 1", catId, userId)
	if checkCat.RowsAffected == 0 {
		c.JSON(404, gin.H{"message":"cat not found"})
		return
	}
	var cat models.Cat
	deleteCat := conn.Raw("DELETE FROM cat WHERE id = ? ", catId).Scan(&cat)
	if deleteCat.Error != nil {
		c.JSON(500, gin.H{"message":"failed to delete cat"})
		return
	}
	c.JSON(200, gin.H{"message":"cat deleted successfully"})
}
