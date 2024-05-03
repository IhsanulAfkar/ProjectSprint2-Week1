package controllers

import (
	"errors"
	"fmt"
	"strconv"
	"week1/db"
	"week1/forms"
	"week1/helper"
	"week1/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MatcherController struct{}

func (h MatcherController) CreateMatch(c *gin.Context){
	userId := c.MustGet("userId")
	conn := db.CreateConn()
	var createMatch forms.CreteMatch
	if err := c.ShouldBindJSON(&createMatch); err != nil {
		c.JSON(400, gin.H{
			"message":err.Error()})
		return
	}
	// check if inputs are uuid
	if !helper.IsUUID(createMatch.MatchCatId) || !helper.IsUUID(createMatch.UserCatId){
		c.JSON(400, gin.H{
			"message":"cat id must be uuid"})
		return
	}
	// check if exists
	listError := make(map[string]interface{})
	var matchCat models.Cat
	var userCat models.Cat
	result := conn.Raw("SELECT * FROM cat WHERE id = ?", createMatch.MatchCatId).Scan(&matchCat)
	result2 := conn.Raw("SELECT * FROM cat WHERE id = ?", createMatch.UserCatId).Scan(&userCat)
	fmt.Println(userCat.HasMatched)
	fmt.Println(matchCat.HasMatched)
	if result.RowsAffected == 0 {
		listError["matchCatId"] = "cat not found"
	}
	if result2.RowsAffected == 0 {
		listError["userCatId"] = "cat not found"
	}
	if len(createMatch.Message) <5 || len(createMatch.Message) > 120 {
		listError["message"] = "message length must be between 5-120 characters"
	}
	if len(listError) > 0 {
		c.JSON(400, gin.H{
			"message":"cat not found",
			"errors":listError})
		return
	}
	// fmt.Println(userCat.UserId, userId)
	userIdStr := fmt.Sprintf("%s", userId)
	val, _ := strconv.Atoi(userIdStr)
	if  userCat.UserId != val{
		c.JSON(400, gin.H{"message":"its not your cat :)"})		
		return
	}
	// if cats owner is same
	if userCat.UserId == matchCat.UserId {
		c.JSON(400, gin.H{"message":"both cat's owned by same user"})		
		return
	}
	// one or those already matched
	
	if userCat.HasMatched {
		c.JSON(400, gin.H{"message":"user cat already matched"})
		return
	}
	if matchCat.HasMatched {
		c.JSON(400, gin.H{"message":"match cat already matched"})
		return
	}
	// same sex
	if matchCat.Sex == userCat.Sex {
		c.JSON(400, gin.H{"message":"same sex not allowed"})
		return	
	}
	var matchResult models.Matcher
	result3 := conn.Raw("INSERT INTO matcher (\"userId\", \"matchUserId\", \"userCatId\", \"matchCatId\", message) VALUES (?,?,?,?,?)", userId, matchCat.UserId, userCat.PkId, matchCat.PkId, createMatch.Message).Scan(&matchResult)
	if result3.Error != nil {
		fmt.Println(result3.Error.Error())
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	c.JSON(201, gin.H{
		"message":"match created successfully"})
}

func (h MatcherController)GetAllMatches(c *gin.Context){
	userIdStr := fmt.Sprintf("%s", c.MustGet("userId"))
	userId,_ := strconv.Atoi(userIdStr) 
	conn := db.CreateConn()
	query := "SELECT * FROM matcher WHERE \"userId\" = ? OR \"matchUserId\" = ? ORDER BY created_at DESC"
	rows, err := conn.Raw(query, userId, userId).Rows()
	if err != nil{
		fmt.Println(err.Error())
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	arrMatcher := make([]any, 0)
	var matcher models.Matcher
	var user models.User
	var userCat models.Cat
	var matchCat models.Cat
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&matcher.PkId, &matcher.Id, &matcher.UserId, &matcher.MatchUserId, &matcher.UserCatId, &matcher.MatchCatId, &matcher.Message, &matcher.IsApproved, &matcher.IsValid, &matcher.CreatedAt, &matcher.UpdatedAt)
		if err != nil {
			fmt.Println(err.Error())
		}
		if !matcher.IsValid {
			continue
		}
		userResult := conn.Raw("SELECT * FROM public.user WHERE \"pkId\" = ?", matcher.UserId).Scan(&user)
		if userResult.Error != nil {
			fmt.Println(userResult.Error.Error())
		}
		userCatResult := conn.Raw("SELECT * FROM cat WHERE \"pkId\" = ? ", matcher.UserCatId).Scan(&userCat)
		if userCatResult.Error != nil {
			fmt.Println(userCatResult.Error.Error())
		}
		matchCatResult := conn.Raw("SELECT * FROM cat WHERE \"pkId\" = ? ", matcher.MatchCatId).Scan(&matchCat)
		if matchCatResult.Error != nil {
			fmt.Println(matchCatResult.Error.Error())
		}
		data := map[string]any{
			"id":matcher.Id,
			"issuedBy": map[string]string{
				"name":user.Name,
				"email":user.Email,
				"createdAt":helper.FormatToIso860(user.CreatedAt)},
			"matchCatDetail": map[string]any {
				"id": matchCat.Id,
				"name":matchCat.Name,
				"race":matchCat.Race,
				"sex":matchCat.Sex,
				"description":matchCat.Description,
				"ageInMonth":matchCat.AgeInMonth,
				"imageUrls":matchCat.ImageUrls,
				"hasMatched": matchCat.HasMatched,
				"createdAt": helper.FormatToIso860(matchCat.CreatedAt)},
			"userCatDetail": map[string]any {
				"id": userCat.Id,
				"name":userCat.Name,
				"race":userCat.Race,
				"sex":userCat.Sex,
				"description":userCat.Description,
				"ageInMonth":userCat.AgeInMonth,
				"imageUrls":userCat.ImageUrls,
				"hasMatched": userCat.HasMatched,
				"createdAt": helper.FormatToIso860(userCat.CreatedAt)},
			"message":matcher.Message,
			"createdAt":helper.FormatToIso860(matcher.CreatedAt)}
		arrMatcher = append(arrMatcher, data)
	}
	c.JSON(200, gin.H{
		"message":"success",
		"data":arrMatcher,
	})
}

func (h MatcherController) ApproveMatch(c *gin.Context){
	userIdStr := fmt.Sprintf("%s", c.MustGet("userId"))
	userId,_ := strconv.Atoi(userIdStr) 
	var approveMatch forms.RequestMatch
	if err:=c.ShouldBindJSON(&approveMatch); err != nil {
		c.JSON(400, gin.H{"message":err.Error()})
		return
	}

	conn := db.CreateConn()
	if !helper.IsUUID(approveMatch.MatchId) {
		c.JSON(400, gin.H{"message":"match id must be uuid"})
		return
	}
	// find match and is still valid?
	var findMatch models.Matcher
	result := conn.Raw("SELECT * FROM matcher WHERE id = ? AND \"matchUserId\" = ?", approveMatch.MatchId, userId).Scan(&findMatch)
	if result.RowsAffected == 0 {
		c.JSON(404,gin.H{"message":"no match found"})
		return
	}
	fmt.Println(findMatch)
	if !findMatch.IsValid {
		c.JSON(400, gin.H{"message":"match no longer valid"})
		return
	}
	// check if both cat are virgin
	var userCat models.Cat
	var matchCat models.Cat
	err := conn.Raw("SELECT * FROM cat WHERE \"pkId\" = ?", findMatch.UserCatId).Scan(&userCat).Error
	err2 := conn.Raw("SELECT * FROM cat WHERE \"pkId\" = ?", findMatch.MatchCatId).Scan(&matchCat).Error
	if err != nil || err2!= nil {
		fmt.Println(err.Error())
		fmt.Println(err2.Error())
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	if userCat.HasMatched {
		c.JSON(400, gin.H{"message":"user cat already matched"})
		return
	}
	if matchCat.HasMatched {
		c.JSON(400, gin.H{"message":"matched cat already matched"})
		return
	}
	finalErr := conn.Transaction(func(tx *gorm.DB) error {
		// update match status
		tx = tx.Debug()
		
		fmt.Println(approveMatch.MatchId)
		txRes := tx.Exec("UPDATE matcher SET \"isApproved\" = true WHERE id = ?", approveMatch.MatchId)
		if txRes.RowsAffected == 0{
			return errors.New("no rows affected")
		}
		
		// update userCat
		txRes = tx.Exec("UPDATE cat SET \"hasMatched\" = true WHERE \"pkId\" = ?",userCat.PkId)
		if txRes.RowsAffected == 0{
			return errors.New("no rows affected")
		}
		
		// update matchCat
		txRes = tx.Exec("UPDATE cat SET \"hasMatched\" = true WHERE \"pkId\" = ?",matchCat.PkId)
		if txRes.RowsAffected == 0{
			return errors.New("no rows affected")
		}
		// delete other records as well
		deleteErr := tx.Exec("DELETE FROM matcher WHERE \"matchCatId\" = ? AND \"userCatId\" = ? AND id != ?", matchCat.PkId , userCat.PkId, approveMatch.MatchId).Error
		if deleteErr != nil {
			return deleteErr
		}
		return nil
	})
	if finalErr != nil {
		fmt.Println(finalErr.Error())
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	c.JSON(200, gin.H{"message":"congratulations! match approved"})
}

func (h MatcherController) RejectMatch(c *gin.Context){
	userIdStr := fmt.Sprintf("%s", c.MustGet("userId"))
	userId,_ := strconv.Atoi(userIdStr) 
	var rejectMatch forms.RequestMatch
	if err:=c.ShouldBindJSON(&rejectMatch); err != nil {
		c.JSON(400, gin.H{"message":err.Error()})
		return
	}
	conn := db.CreateConn()
	if !helper.IsUUID(rejectMatch.MatchId){
		c.JSON(400, gin.H{"message":"match id must be uuid"})
	}
	var findMatch models.Matcher
	result := conn.Raw("SELECT * FROM matcher WHERE id = ? AND \"matchUserId\" = ?", rejectMatch.MatchId, userId).Scan(&findMatch)
	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"message":"no match found"})
	}
	fmt.Println(findMatch)
	// get cats
	var userCat models.Cat
	var matchCat models.Cat
	result = conn.Raw("SELECT * FROM cat WHERE id = ?",findMatch.UserCatId).Scan(&userCat)
	if result.RowsAffected == 0 {
		c.JSON(400, gin.H{"message": "user cat already deleted"})
		return
	}
	result = conn.Raw("SELECT * FROM cat WHERE id = ?",findMatch.MatchCatId).Scan(&matchCat)
	if result.RowsAffected == 0 {
		c.JSON(400, gin.H{"message": "match cat already deleted"})
		return
	}
	if matchCat.HasMatched {
		c.JSON(400, gin.H{"message": "match cat already matched"})
		return
	}
	if userCat.HasMatched {
		c.JSON(400, gin.H{"message": "user cat already matched"})
		return
	}
	if !findMatch.IsValid {
		c.JSON(400, gin.H{"message":"match no longer valid"})
		return
	}
	if findMatch.IsApproved {
		c.JSON(400, gin.H{"message":"match already approved"})
		return
	}
	result = conn.Exec("UPDATE matcher SET \"isApproved\" = true WHERE id = ?", findMatch.Id)
	if result.RowsAffected == 0{
		c.JSON(500,gin.H{"message":"server error"})
		return
	}
	c.JSON(200, gin.H{"message":"match rejected successfully"})	
}

func (h MatcherController)DeleteMatch(c *gin.Context){
	userId := c.MustGet("userId")
	matchId := c.Param("matchId")
	if !helper.IsUUID(matchId) {
		c.JSON(400, gin.H{"message":"match id is not uuid"})
		return
	}
	conn := db.CreateConn()
	var findMatch models.Matcher
	result := conn.Raw("SELECT * FROM matcher WHERE id = ? AND \"userId\" = ?", matchId, userId).Scan(&findMatch)
	if result.Error != nil {
		c.JSON(500, gin.H{"message":"server error"})
		return	
	}

	if result.RowsAffected == 0{
		c.JSON(404, gin.H{"message":"match not found"})
		return
	}
	if findMatch.IsApproved {
		c.JSON(400, gin.H{"message":"match already approved"})
		return
	}
	if !findMatch.IsValid {
		c.JSON(400, gin.H{"message":"match already rejected"})
	}
	result = conn.Exec("DELETE FROM matcher WHERE id = ?", matchId)
	if result.RowsAffected ==0 {
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	c.JSON(200, gin.H{"message":"match deleted successfully"})
}