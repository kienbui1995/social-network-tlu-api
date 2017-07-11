package controllers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/configs"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// ClassController struct
type ClassController struct {
	Service services.ClassServiceInterface
}

// Admin

// GetAll func
func (controller ClassController) GetAll(c *gin.Context) {
	// code, errParseInt := strconv.ParseInt(c.Query("student_code"), 10, 64)
	// if errParseInt != nil {
	// 	helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid params")
	// 	return
	// }

	//check permisson
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, 200, "Permissions error")
		return
	}
	role, errGetRoleFromUserID := GetRoleFromUserID(myUserID)
	if errGetRoleFromUserID != nil {
		helpers.ResponseAuthJSON(c, 200, "Permissions error")
		fmt.Printf("GetRoleFromUserID controller: %s\n", errGetRoleFromUserID.Error())
		return
	}
	if role != configs.IStudentRole && role != configs.IAdminRole && role != configs.ISupervisorRole && role != configs.ITeacherRole {
		helpers.ResponseAuthJSON(c, 200, "Permissions error")
		return
	}
	params := helpers.ParamsGetAll{}
	params.Skip, _ = strconv.Atoi(c.DefaultQuery("skip", configs.SSkip))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", configs.SLimit))
	// params.Type = c.DefaultQuery("type", configs.SPost)
	// if (params.Type != configs.SPostPhoto && params.Type != configs.SPostStatus) && params.Type != configs.SPost {
	// 	helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: type")
	// 	return
	// }
	params.Sort = c.DefaultQuery("sort", configs.SSort)
	params.Sort, _ = helpers.ConvertSort(params.Sort)
	classes, errGetAll := controller.Service.GetAll(params)
	if errGetAll != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetAll service: %s\n", errGetAll.Error())
		return
	}
	helpers.ResponseEntityListJSON(c, 1, "Get class list successful", classes, params, len(classes))
}

// // Get func
// func (controller SemesterController) Get(c *gin.Context) {
// 	postID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
// 	if errParseInt != nil {
// 		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid post id")
// 		return
// 	}
// 	// check exists
// 	exist, errCheckExistPost := controller.Service.CheckExistPost(postID)
// 	if errCheckExistPost != nil {
// 		helpers.ResponseServerErrorJSON(c)
// 		fmt.Printf("CheckExistPost service: %s\n", errCheckExistPost.Error())
// 	}
// 	if exist != true {
// 		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "No exist post")
// 		return
// 	}
//
// 	// get myuserID
// 	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
// 	if errGetUserIDFromToken != nil {
// 		helpers.ResponseServerErrorJSON(c)
// 		fmt.Printf("GetUserIDFromToken func: %s\n", errGetUserIDFromToken.Error())
// 	}
// 	post, errGet := controller.Service.Get(postID, myUserID)
// 	if errGet != nil {
// 		helpers.ResponseServerErrorJSON(c)
// 		fmt.Printf("Get service: %s\n", errGet.Error())
// 	}
//
// 	// Valida se existe a pessoa (404)
// 	if post.IsEmpty() {
// 		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not Found Post")
// 		return
// 	}
//
// 	helpers.ResponseSuccessJSON(c, 1, "Info of user", post)
// }
//
// // Delete func
// func (controller SemesterController) Delete(c *gin.Context) {
// 	postID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
// 	if errParseInt != nil {
// 		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid post_id")
// 		return
// 	}
// 	// check exists
// 	exist, errCheckExistPost := controller.Service.CheckExistPost(postID)
// 	if errCheckExistPost != nil {
// 		helpers.ResponseServerErrorJSON(c)
// 		fmt.Printf("CheckExistPost service: %s\n", errCheckExistPost.Error())
// 	}
// 	if exist != true {
// 		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "No exist post")
// 		return
// 	}
//
// 	userID, errGetUserIDByPostID := controller.Service.GetUserIDByPostID(postID)
// 	if errGetUserIDByPostID != nil {
// 		helpers.ResponseServerErrorJSON(c)
// 		fmt.Printf("GetUserIDByPostID service: %s\n", errGetUserIDByPostID.Error())
// 		return
// 	}
//
// 	// check permisson
// 	if myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token")); myUserID != userID || errGetUserIDFromToken != nil {
// 		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permissions error")
// 		return
// 	}
//
// 	// delete action
// 	deleted, errDelete := controller.Service.Delete(postID)
// 	if errDelete == nil && deleted == true {
// 		helpers.ResponseNoContentJSON(c)
//
// 		// auto noti
// 		go func() {
//
// 		}()
//
// 		return
// 	}
// 	if errDelete != nil {
// 		helpers.ResponseServerErrorJSON(c)
// 		fmt.Printf("Delete service: %s\n", errDelete)
// 	}
// }
//
// // Create func
// func (controller SemesterController) Create(c *gin.Context) {
// 	userID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
// 	if errParseInt != nil {
// 		helpers.ResponseBadRequestJSON(c, configs.EcParamUserID, "Invalid user id")
// 	}
//
// 	// Check permisson
// 	if myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token")); myUserID != userID || errGetUserIDFromToken != nil {
// 		helpers.ResponseAuthJSON(c, 200, "Permissions error")
// 		return
// 	}
//
// 	// json struct
// 	json := struct {
// 		Photo   string `json:"photo"`
// 		Message string `json:"message"`
// 		Privacy int    `json:"privacy"`
// 		Status  int    `json:"status"`
// 	}{}
//
// 	// BadRequest (400)
// 	if errBindJSON := c.BindJSON(&json); errBindJSON != nil {
// 		helpers.ResponseBadRequestJSON(c, configs.EcParam, "BindJSON: "+errBindJSON.Error())
// 		return
// 	}
//
// 	// validation
// 	if len(json.Message) == 0 {
// 		helpers.ResponseJSON(c, 400, 100, "Missing a few fields:  Message is NULL", nil)
// 		return
// 	}
// 	if json.Status == 0 {
// 		json.Status = 1
// 	}
// 	if json.Privacy == 0 {
// 		json.Privacy = 1
// 	}
//
// 	post := models.Post{}
// 	helpers.Replace(json, &post)
// 	postID, errCreate := controller.Service.Create(post, userID)
// 	if errCreate == nil && postID >= 0 {
// 		helpers.ResponseSuccessJSON(c, 1, "Create user post successful", map[string]interface{}{"id": postID})
//
// 		// // auto noti
// 		// go func() {
// 		//
// 		// }()
//
// 		// push noti
// 		go func() {
// 			Notification := NotificationController{Service: services.NewNotificationService()}
// 			err := Notification.UpdatePostNotification(userID)
// 			if err != nil {
// 				fmt.Printf("UpdateLikeNotification: %s\n", err.Error())
// 			}
// 		}()
// 		return
// 	}
// 	helpers.ResponseServerErrorJSON(c)
// 	if errCreate != nil {
// 		fmt.Printf("Create services: %s\n", errCreate.Error())
// 	} else {
// 		fmt.Printf("Create services: Don't create Post")
// 	}
// }

// UpdateFromTLU func
func (controller ClassController) UpdateFromTLU(c *gin.Context) {
	fromSemester, _ := strconv.ParseInt(c.Query("from"), 10, 64)
	toSemester, _ := strconv.ParseInt(c.Query("to"), 10, 64)

	// get myUserID from token
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil || myUserID < 0 {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permissions error")
		return
	}
	// check role
	role, errGetRoleFromUserID := GetRoleFromUserID(myUserID)
	if errGetRoleFromUserID != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetRoleFromUserID controller: %s\n", errGetRoleFromUserID.Error())
		return
	}
	if role != configs.IAdminRole {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permissions error")
		return
	}
	// var newPost models.Post
	// if errBindJSON := c.BindJSON(&newPost); errBindJSON != nil {
	// 	helpers.ResponseBadRequestJSON(c, configs.EcParam, "BindJSON: "+errBindJSON.Error())
	// 	return
	// }
	for index := fromSemester; index <= toSemester; index++ {

		_, errUpdate := controller.Service.UpdateFromTLU(strconv.FormatInt(index, 10))

		if errUpdate != nil {
			fmt.Printf("UpdateFromTLU service: %s\n", errUpdate.Error())
			// return
		}

	}

	helpers.ResponseSuccessJSON(c, 1, "Update class successful", nil)
}

// Student

// GetAllByStudent func
func (controller ClassController) GetAllByStudent(c *gin.Context) {
	studentCode := c.Query("student_code")
	semesterCode := c.Query("semester_code")

	//check permisson
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, 200, "Permissions error")
		return
	}
	role, errGetRoleFromUserID := GetRoleFromUserID(myUserID)
	if errGetRoleFromUserID != nil {
		helpers.ResponseAuthJSON(c, 200, "Permissions error")
		fmt.Printf("GetRoleFromUserID controller: %s\n", errGetRoleFromUserID.Error())
		return
	}
	if role != configs.IStudentRole && role != configs.IAdminRole && role != configs.ISupervisorRole && role != configs.ITeacherRole {
		fmt.Printf("role: %v\n", role)
		helpers.ResponseAuthJSON(c, 200, "Permissions error")
		return
	}
	params := helpers.ParamsGetAll{}
	params.Skip, _ = strconv.Atoi(c.DefaultQuery("skip", configs.SSkip))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", configs.SLimit))
	// params.Type = c.DefaultQuery("type", configs.SPost)
	// if (params.Type != configs.SPostPhoto && params.Type != configs.SPostStatus) && params.Type != configs.SPost {
	// 	helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: type")
	// 	return
	// }
	params.Sort = c.DefaultQuery("sort", configs.SSort)
	params.Sort, _ = helpers.ConvertSort(params.Sort)
	classes, errGetAll := controller.Service.GetAllByStudent(params, semesterCode, studentCode)
	if errGetAll != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetAll service: %s\n", errGetAll.Error())
		return
	}
	helpers.ResponseEntityListJSON(c, 1, "Get class list successful", classes, params, len(classes))
}

// GetAllByRoom func
func (controller ClassController) GetAllByRoom(c *gin.Context) {
	roomCode := c.Param("id")
	day := c.Query("day")

	//check permisson
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, 200, "Permissions error")
		return
	}
	role, errGetRoleFromUserID := GetRoleFromUserID(myUserID)
	if errGetRoleFromUserID != nil {
		helpers.ResponseAuthJSON(c, 200, "Permissions error")
		fmt.Printf("GetRoleFromUserID controller: %s\n", errGetRoleFromUserID.Error())
		return
	}
	if role != configs.IAdminRole && role != configs.ISupervisorRole && role != configs.ITeacherRole {
		//fmt.Printf("role: %v\n", role)
		helpers.ResponseAuthJSON(c, 200, "Permissions error")
		return
	}
	params := helpers.ParamsGetAll{}
	params.Skip, _ = strconv.Atoi(c.DefaultQuery("skip", configs.SSkip))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", configs.SLimit))
	// params.Type = c.DefaultQuery("type", configs.SPost)
	// if (params.Type != configs.SPostPhoto && params.Type != configs.SPostStatus) && params.Type != configs.SPost {
	// 	helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: type")
	// 	return
	// }
	params.Sort = c.DefaultQuery("sort", configs.SSort)
	params.Sort, _ = helpers.ConvertSort(params.Sort)
	classes, errGetAll := controller.Service.GetAllByRoom(params, day, roomCode)
	if errGetAll != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetAll service: %s\n", errGetAll.Error())
		return
	}
	helpers.ResponseEntityListJSON(c, 1, "Get class list successful", classes, params, len(classes))
}

// // GetAllByStudent func
// func (controller ClassController) GetAllByStudent(c *gin.Context) {
// 	studentCode := c.Query("student_code")
// 	semesterCode := c.Query("semester_code")
//
// 	//check permisson
// 	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
// 	if errGetUserIDFromToken != nil {
// 		helpers.ResponseAuthJSON(c, 200, "Permissions error")
// 		return
// 	}
// 	role, errGetRoleFromUserID := GetRoleFromUserID(myUserID)
// 	if errGetRoleFromUserID != nil {
// 		helpers.ResponseAuthJSON(c, 200, "Permissions error")
// 		fmt.Printf("GetRoleFromUserID controller: %s\n", errGetRoleFromUserID.Error())
// 		return
// 	}
// 	if role != configs.IStudentRole && role != configs.IAdminRole && role != configs.ISupervisorRole && role != configs.ITeacherRole {
// 		fmt.Printf("role: %v\n", role)
// 		helpers.ResponseAuthJSON(c, 200, "Permissions error")
// 		return
// 	}
// 	params := helpers.ParamsGetAll{}
// 	params.Skip, _ = strconv.Atoi(c.DefaultQuery("skip", configs.SSkip))
// 	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", configs.SLimit))
// 	// params.Type = c.DefaultQuery("type", configs.SPost)
// 	// if (params.Type != configs.SPostPhoto && params.Type != configs.SPostStatus) && params.Type != configs.SPost {
// 	// 	helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: type")
// 	// 	return
// 	// }
// 	params.Sort = c.DefaultQuery("sort", configs.SSort)
// 	params.Sort, _ = helpers.ConvertSort(params.Sort)
// 	classes, errGetAll := controller.Service.GetAllByStudent(params, semesterCode, studentCode)
// 	if errGetAll != nil {
// 		helpers.ResponseServerErrorJSON(c)
// 		fmt.Printf("GetAll service: %s\n", errGetAll.Error())
// 		return
// 	}
// 	helpers.ResponseEntityListJSON(c, 1, "Get class list successful", classes, params, len(classes))
// }
