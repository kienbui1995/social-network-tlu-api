package controllers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/configs"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/models"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// PostController controller
type PostController struct {
	Service services.PostServiceInterface
}

// GetAll func
func (controller PostController) GetAll(c *gin.Context) {
	userID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParamUserID, "Invalid user id")
		return
	}

	//check permisson
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, 200, "Permissions error")
		return
	}
	params := helpers.ParamsGetAll{}
	params.Skip, _ = strconv.Atoi(c.DefaultQuery("skip", configs.SSkip))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", configs.SLimit))
	params.Type = c.DefaultQuery("type", configs.SPost)
	if (params.Type != configs.SPostPhoto && params.Type != configs.SPostStatus) && params.Type != configs.SPost {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: type")
		return
	}
	params.Sort = c.DefaultQuery("sort", configs.SSort)
	params.Sort, _ = helpers.ConvertSort(params.Sort)
	posts, errGetAll := controller.Service.GetAll(params, userID, myUserID)
	if errGetAll != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetAll service: %s\n", errGetAll.Error())
		return
	}
	helpers.ResponseEntityListJSON(c, 1, "Get user list successful", posts, params, len(posts))
}

// Get func
func (controller PostController) Get(c *gin.Context) {
	postID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid post id")
		return
	}
	// check exists
	exist, errCheckExistPost := controller.Service.CheckExistPost(postID)
	if errCheckExistPost != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistPost service: %s\n", errCheckExistPost.Error())
	}
	if exist != true {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "No exist post")
		return
	}

	// get myuserID
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDFromToken func: %s\n", errGetUserIDFromToken.Error())
	}
	post, errGet := controller.Service.Get(postID, myUserID)
	if errGet != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Get service: %s\n", errGet.Error())
	}

	// Valida se existe a pessoa (404)
	if post.IsEmpty() {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not Found Post")
		return
	}

	helpers.ResponseSuccessJSON(c, 1, "Info of user", post)
}

// Delete func
func (controller PostController) Delete(c *gin.Context) {
	postID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid post_id")
		return
	}
	// check exists
	exist, errCheckExistPost := controller.Service.CheckExistPost(postID)
	if errCheckExistPost != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistPost service: %s\n", errCheckExistPost.Error())
	}
	if exist != true {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "No exist post")
		return
	}

	userID, errGetUserIDByPostID := controller.Service.GetUserIDByPostID(postID)
	if errGetUserIDByPostID != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDByPostID service: %s\n", errGetUserIDByPostID.Error())
		return
	}

	// check permisson
	if myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token")); myUserID != userID || errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permissions error")
		return
	}

	// delete action
	deleted, errDelete := controller.Service.Delete(postID)
	if errDelete == nil && deleted == true {
		helpers.ResponseNoContentJSON(c)

		// auto Decrease Posts
		go func() {
			ok, errDecreasePosts := controller.Service.DecreasePosts(userID)
			if errDecreasePosts != nil {
				fmt.Printf("DecreasePosts service: %s\n", errDecreasePosts.Error())
			}
			if ok != true {
				fmt.Printf("DecreasePosts service: Don't decrease\n")
			}
		}()

		return
	}
	if errDelete != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Delete service: %s\n", errDelete)
	}
}

// Create func
func (controller PostController) Create(c *gin.Context) {
	userID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParamUserID, "Invalid user id")
	}

	// Check permisson
	if myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token")); myUserID != userID || errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, 200, "Permissions error")
		return
	}

	// json struct
	json := struct {
		Photo   string `json:"photo"`
		Message string `json:"message"`
		Privacy int    `json:"privacy"`
		Status  int    `json:"status"`
	}{}

	// BadRequest (400)
	if errBindJSON := c.BindJSON(&json); errBindJSON != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "BindJSON: "+errBindJSON.Error())
		return
	}

	// validation
	if len(json.Message) == 0 {
		helpers.ResponseJSON(c, 400, 100, "Missing a few fields:  Message is NULL", nil)
		return
	}
	if json.Status == 0 {
		json.Status = 1
	}
	if json.Privacy == 0 {
		json.Privacy = 1
	}

	action := " cập nhật trạng thái"
	if len(json.Photo) > 0 {
		action = " đăng ảnh"
	}
	post := models.Post{}
	helpers.Replace(json, &post)
	postID, errCreate := controller.Service.Create(post, userID)
	if errCreate == nil && postID >= 0 {
		helpers.ResponseSuccessJSON(c, 1, "Create user post successful", map[string]interface{}{"id": postID})

		// auto Increase Posts
		go func() {
			ok, errIncreasePosts := controller.Service.IncreasePosts(userID)
			if errIncreasePosts != nil {
				fmt.Printf("IncreasePosts service: %s\n", errIncreasePosts.Error())
			}
			if ok != true {
				fmt.Printf("IncreasePosts service: Don't increase\n")
			}
		}()

		// push noti
		go func() {
			user, _ := services.NewUserService().Get(userID)
			notify := models.Notification{
				ObjectID:   postID,
				ObjectType: "post",
				Title:      "@" + user.Username + action,
				Message:    json.Message,
			}
			ids, errGetIDs := services.NewSubscriberService().GetFollowerIDs(userID)
			if len(ids) > 0 && errGetIDs == nil {
				for index := 0; index < len(ids); index++ {
					notify.UserID = ids[index]
					PushTest(notify)
				}
			}
		}()
		return
	}
	helpers.ResponseServerErrorJSON(c)
	if errCreate != nil {
		fmt.Printf("Create services: %s\n", errCreate.Error())
	} else {
		fmt.Printf("Create services: Don't create Post")
	}
}

// Update func
func (controller PostController) Update(c *gin.Context) {
	postID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	exist, errCheckExistPost := controller.Service.CheckExistPost(postID)
	if errCheckExistPost != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistPost service: %s\n", errCheckExistPost.Error())
		return
	}
	if exist == false {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Don't exist post")
		return
	}

	// get myUserID from token
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil || myUserID < 0 {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permissions error")
		return
	}
	olderPost, errGet := controller.Service.Get(postID, myUserID)
	if errGet != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Replace helpers: %s\n", errGet.Error())
		return
	}
	if myUserID != olderPost.Owner.ID {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permissions error")
		return
	}

	var newPost models.Post
	if errBindJSON := c.BindJSON(&newPost); errBindJSON != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "BindJSON: "+errBindJSON.Error())
		return
	}

	errReplace := helpers.Replace(olderPost, &newPost)
	if errReplace != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Replace helpers: %s\n", errReplace.Error())
	}

	userUpdate, errUpdate := controller.Service.Update(newPost)
	if errUpdate != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Update service: %s\n", errUpdate.Error())
		return
	}

	helpers.ResponseSuccessJSON(c, 1, "Update user successful", userUpdate)
}

// CreateLike func
func (controller PostController) CreateLike(c *gin.Context) {
	postID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: post_id")
		return
	}

	exist, errCheckExistPost := controller.Service.CheckExistPost(postID)
	if errCheckExistPost != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistPost service: %s\n", errCheckExistPost.Error())
		return
	}
	if exist != true {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not exist object")
		return
	}
	//check permisson
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDFromToken controller: %s\n", errGetUserIDFromToken.Error())
		return
	}
	// check liked
	if liked, _ := controller.Service.CheckExistLike(postID, myUserID); liked == true {
		helpers.ResponseBadRequestJSON(c, configs.EcExisObject, "Exist this object: Likes")
		return
	}

	likes, errCreateLike := controller.Service.CreateLike(postID, myUserID)
	if errCreateLike == nil && likes >= 0 {
		helpers.ResponseSuccessJSON(c, 1, "Like post successful", map[string]int{"likes": likes})

		// auto Increase post Likes
		go func() {
			ok, errIncreaseLikes := controller.Service.IncreaseLikes(postID)
			if errIncreaseLikes != nil {
				fmt.Printf("IncreaseLikes service: %s\n", errIncreaseLikes.Error())
			}
			if ok != true {
				fmt.Printf("IncreaseLikes service: don't increase like")
			}
		}()

		// push noti
		go func() {
			post, _ := controller.Service.Get(postID, myUserID)
			userLiked, _ := services.NewUserService().Get(myUserID)
			notify := models.Notification{
				UserID:     post.Owner.ID,
				ObjectID:   post.ID,
				ObjectType: "post",
				Title:      "@" + userLiked.Username + " thích trạng thái của bạn",
				Message:    post.Message,
			}
			PushTest(notify)
		}()
		return
	}

	helpers.ResponseServerErrorJSON(c)
	if errCreateLike != nil {
		fmt.Printf("CreateLike services: %s\n", errCreateLike.Error())
	} else {
		fmt.Printf("CreateLike services: Don't Create Like\n")
	}
}

// DeleteLike func
func (controller PostController) DeleteLike(c *gin.Context) {
	postID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid post id")
		return
	}
	//check exist
	exist, errCheckExistPost := controller.Service.CheckExistPost(postID)
	if errCheckExistPost != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistPost service: %s\n", errCheckExistPost.Error())
		return
	}
	if exist != true {
		helpers.ResponseBadRequestJSON(c, configs.EcNoExistObject, "No exist this object: post")
		return
	}

	//check permisson
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, 200, "Permissions error")
		return
	}

	// check liked
	if liked, _ := controller.Service.CheckExistLike(postID, myUserID); liked != true {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "No Exist this object: Likes")
		return
	}

	likes, errDeleteLike := controller.Service.DeleteLike(postID, myUserID)
	if errDeleteLike == nil && likes >= 0 {
		helpers.ResponseSuccessJSON(c, 1, "Unlike successful", map[string]int{"likes": likes})

		// auto Decrease post Likes
		go func() {
			ok, errDecreaseLikes := controller.Service.DecreaseLikes(postID)
			if errDecreaseLikes != nil {
				fmt.Printf("DecreaseLikes service: %s\n", errDecreaseLikes.Error())
			}
			if ok != true {
				fmt.Printf("DecreaseLikes service: don't decrease\n")
			}
		}()
		return
	}

	helpers.ResponseServerErrorJSON(c)
	if errDeleteLike != nil {
		fmt.Printf("DeletePostLike services: %s\n", errDeleteLike.Error())
	} else {
		fmt.Printf("DeletePostLike services: Don't Delete Like\n")
	}

}

// GetLikes func
func (controller PostController) GetLikes(c *gin.Context) {
	postID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, 100, "Invalid parameter: "+errParseInt.Error())
		return
	}

	//check permisson
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDFromToken controller: %s\n", errGetUserIDFromToken.Error())
		return
	}

	var params helpers.ParamsGetAll
	params.Sort = c.DefaultQuery("sort", "-liked_at")
	var errConvertSort error
	params.Sort, errConvertSort = helpers.ConvertSort(params.Sort)
	if errConvertSort != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: "+errConvertSort.Error())
		return
	}
	var errSkip error
	params.Skip, errSkip = strconv.Atoi(c.DefaultQuery("skip", configs.SSkip))
	if errSkip != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: "+errSkip.Error())
		return
	}
	var errLimit error
	params.Limit, errLimit = strconv.Atoi(c.DefaultQuery("limit", configs.SLimit))
	if errLimit != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: "+errLimit.Error())
		return
	}

	likeList, errGetLikes := controller.Service.GetLikes(postID, myUserID, params)
	if errGetLikes == nil {
		helpers.ResponseEntityListJSON(c, 1, " Posts Likes List", likeList, params, len(likeList))
		return
	}

	helpers.ResponseServerErrorJSON(c)
	if errGetLikes != nil {
		fmt.Printf("GetLikes services: %s\n", errGetLikes.Error())
	} else {
		fmt.Printf("GetLikes services: Don't GetLikes\n")
	}
}

// CreateFollow func
func (controller PostController) CreateFollow(c *gin.Context) {
	postID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: post_id")
		return
	}

	exist, errCheckExistPost := controller.Service.CheckExistPost(postID)
	if errCheckExistPost != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistPost service: %s\n", errCheckExistPost.Error())
		return
	}
	if exist != true {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not exist object")
		return
	}
	//check permisson
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDFromToken controller: %s\n", errGetUserIDFromToken.Error())
		return
	}
	// check liked
	if liked, _ := controller.Service.CheckExistLike(postID, myUserID); liked == true {
		helpers.ResponseBadRequestJSON(c, configs.EcExisObject, "Exist this object: Likes")
		return
	}

	likes, errCreateLike := controller.Service.CreateLike(postID, myUserID)
	if errCreateLike == nil && likes >= 0 {
		helpers.ResponseSuccessJSON(c, 1, "Like post successful", map[string]int{"likes": likes})

		// auto Increase post Likes
		go func() {
			ok, errIncreaseLikes := controller.Service.IncreaseLikes(postID)
			if errIncreaseLikes != nil {
				fmt.Printf("IncreaseLikes service: %s\n", errIncreaseLikes.Error())
			}
			if ok != true {
				fmt.Printf("IncreaseLikes service: don't increase like")
			}
		}()

		// push noti
		go func() {
			post, _ := controller.Service.Get(postID, myUserID)
			userLiked, _ := services.NewUserService().Get(myUserID)
			notify := models.Notification{
				UserID:     post.Owner.ID,
				ObjectID:   post.ID,
				ObjectType: "post",
				Title:      "@" + userLiked.Username + " thích trạng thái của bạn",
				Message:    post.Message,
			}
			PushTest(notify)
		}()
		return
	}

	helpers.ResponseServerErrorJSON(c)
	if errCreateLike != nil {
		fmt.Printf("CreateLike services: %s\n", errCreateLike.Error())
	} else {
		fmt.Printf("CreateLike services: Don't Create Like\n")
	}
}

// DeleteFollow func
func (controller PostController) DeleteFollow(c *gin.Context) {

}

// CreateReport func
func (controller PostController) CreateReport(c *gin.Context) {

}

// DeleteReport func
func (controller PostController) DeleteReport(c *gin.Context) {

}

// GetUsers func to get mentioned users or can_mention users
func (controller PostController) GetUsers(c *gin.Context) {
	postID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: post_id")
		return
	}

	exist, errCheckExistPost := controller.Service.CheckExistPost(postID)
	if errCheckExistPost != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistPost service: %s\n", errCheckExistPost.Error())
		return
	}
	if exist != true {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not exist object")
		return
	}

	//check permisson
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDFromToken controller: %s\n", errGetUserIDFromToken.Error())
		return
	}
	params := helpers.ParamsGetAll{}
	params.Type = configs.SCanMention
	params.Type = c.DefaultQuery("type", params.Type)
	params.Skip, _ = strconv.Atoi(c.DefaultQuery("skip", configs.SSkip))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", configs.SLimit))
	params.Sort = c.DefaultQuery("sort", configs.SSort)
	params.Sort, _ = helpers.ConvertSort(params.Sort)

	// get can mentioned users
	if params.Type == configs.SCanMention {
		users, errGetCanMentionedUsers := controller.Service.GetCanMentionedUsers(postID, params, myUserID)
		if errGetCanMentionedUsers != nil {
			helpers.ResponseServerErrorJSON(c)
			fmt.Printf("GetCanMentionedUsers service: %s\n", errGetCanMentionedUsers.Error())
			return
		}
		helpers.ResponseEntityListJSON(c, 1, "Get user could mentioned successful", users, params, len(users))
		return
	}

	// get mentioned users
	if params.Type == configs.SMentioned {
		users, errGetMentionedUsers := controller.Service.GetMentionedUsers(postID, params, myUserID)
		if errGetMentionedUsers != nil {
			helpers.ResponseServerErrorJSON(c)
			fmt.Printf("GetMentionedUsers service: %s\n", errGetMentionedUsers.Error())
			return
		}
		helpers.ResponseEntityListJSON(c, 1, "Get mentioned users successful", users, params, len(users))
		return
	}

	// get liked users
	if params.Type == configs.SLikedPost {
		users, errGetLikedUsers := controller.Service.GetLikedUsers(postID, params, myUserID)
		if errGetLikedUsers != nil {
			helpers.ResponseServerErrorJSON(c)
			fmt.Printf("GetLikedUsers service: %s\n", errGetLikedUsers.Error())
			return
		}
		helpers.ResponseEntityListJSON(c, 1, "Get liked users successful", users, params, len(users))
		return
	}

	// get followed users
	if params.Type == configs.SFollowedPost {
		users, errGetFollowedUsers := controller.Service.GetFollowedUsers(postID, params, myUserID)
		if errGetFollowedUsers != nil {
			helpers.ResponseServerErrorJSON(c)
			fmt.Printf("GetFollowedUsers service: %s\n", errGetFollowedUsers.Error())
			return
		}
		helpers.ResponseEntityListJSON(c, 1, "Get followed users successful", users, params, len(users))
		return
	}

}

// CreateGroupPost func
func (controller PostController) CreateGroupPost(c *gin.Context) {
	groupID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParamUserID, "Invalid group id")
		return
	}

	// Check permisson
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, 200, "Permissions error")
		return
	}
	group, errGet := services.NewGroupService().Get(groupID, myUserID)
	if errGet != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Get Group service: %s\n", errGet.Error())
		return
	}
	if group.IsEmpty() {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "No exist group")
		return
	}

	role, errCheckUserRole := services.NewGroupService().CheckUserRole(groupID, myUserID)
	if errCheckUserRole != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckUserRole service: %s\n", errCheckUserRole.Error())
		return
	}

	if role != configs.IAdmin && role != configs.IMember {
		helpers.ResponseForbiddenJSON(c, 200, "Permissions error")
		return
	}
	// json struct
	json := struct {
		Photo   string `json:"photo"`
		Message string `json:"message"`
		Privacy int    `json:"privacy"`
		Status  int    `json:"status"`
	}{}

	// BadRequest (400)
	if errBindJSON := c.BindJSON(&json); errBindJSON != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "BindJSON: "+errBindJSON.Error())
		return
	}

	// validation
	if len(json.Message) == 0 {
		helpers.ResponseJSON(c, 400, 100, "Missing a few fields:  Message is NULL", nil)
		return
	}
	if json.Status == 0 {
		json.Status = 1
	}
	if json.Privacy == 0 {
		json.Privacy = 1
	}
	// ~doing ~needfix
	action := " viết bài"
	if len(json.Photo) > 0 {
		action = " đăng ảnh"
	}
	action += " trong " + group.Name
	post := models.Post{}

	helpers.Replace(json, &post)
	postID, errCreate := controller.Service.CreateGroupPost(post, groupID, myUserID)
	if errCreate == nil && postID >= 0 {
		helpers.ResponseSuccessJSON(c, 1, "Create user post successful", map[string]interface{}{"id": postID})

		// push noti
		go func() {
			user, _ := services.NewUserService().Get(myUserID)
			notify := models.Notification{
				ObjectID:   postID,
				ObjectType: "post",
				Title:      "@" + user.Username + action,
				Message:    json.Message,
			}
			ids, errGetIDs := services.NewSubscriberService().GetFollowerIDs(myUserID)
			if len(ids) > 0 && errGetIDs == nil {
				for index := 0; index < len(ids); index++ {
					notify.UserID = ids[index]
					PushTest(notify)
				}
			}
		}()
		return
	}
	helpers.ResponseServerErrorJSON(c)
	if errCreate != nil {
		fmt.Printf("Create services: %s\n", errCreate.Error())
	} else {
		fmt.Printf("Create services: Don't create group post")
	}
}

// GetAllGroupPost func
func (controller PostController) GetAllGroupPost(c *gin.Context) {
	groupID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParamUserID, "Invalid group id")
		return
	}

	//check permisson
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, 200, "Permissions error")
		return
	}
	params := helpers.ParamsGetAll{}
	params.Skip, _ = strconv.Atoi(c.DefaultQuery("skip", configs.SSkip))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", configs.SLimit))
	params.Type = c.DefaultQuery("type", configs.SPost)
	if (params.Type != configs.SPostPhoto && params.Type != configs.SPostStatus) && params.Type != configs.SPost {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: type")
		return
	}
	params.Sort = c.DefaultQuery("sort", configs.SSort)
	params.Sort, _ = helpers.ConvertSort(params.Sort)
	posts, errGetAll := controller.Service.GetAllGroupPosts(params, groupID, myUserID)
	if errGetAll != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetAll service: %s\n", errGetAll.Error())
		return
	}
	helpers.ResponseEntityListJSON(c, 1, "Get user list successful", posts, params, len(posts))
}
