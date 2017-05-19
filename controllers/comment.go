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

// CommentController controller
type CommentController struct {
	Service services.CommentServiceInterface
}

// GetAll func
func (controller CommentController) GetAll(c *gin.Context) {
	postID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: post_id")
		return
	}

	//check permisson
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, 200, "Permissions error")
		return
	}
	if allowed, _ := controller.Service.CheckPostInteractivePermission(postID, myUserID); allowed == false {
		helpers.ResponseForbiddenJSON(c, configs.EcPermissionPost, "Post not visible")
		return
	}
	params := helpers.ParamsGetAll{}
	var err error
	params.Skip, err = strconv.Atoi(c.DefaultQuery("skip", configs.SSkip))
	if err != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: skip")
		return
	}
	params.Limit, err = strconv.Atoi(c.DefaultQuery("limit", configs.SLimit))
	if err != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: limit")
		return
	}
	params.Sort = c.DefaultQuery("sort", configs.SSort)
	params.Sort, _ = helpers.ConvertSort(params.Sort)
	comments, errGetAll := controller.Service.GetAll(postID, params, myUserID)
	if errGetAll != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetAll service: %s\n", errGetAll.Error())
		return
	}
	helpers.ResponseEntityListJSON(c, 1, "Get comment list successful", comments, params, len(comments))
}

// Get func
func (controller CommentController) Get(c *gin.Context) {
	commentID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid comment id")
		return
	}

	//check exist
	exist, errCheckExistComment := controller.Service.CheckExistComment(commentID)
	if errCheckExistComment != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistComment service: %s\n", errCheckExistComment.Error())
	}

	if exist != true {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "No exist comment")
		return
	}

	comment, errGet := controller.Service.Get(commentID)
	if errGet == nil && comment.ID == commentID {
		helpers.ResponseSuccessJSON(c, 1, "Get comment successful", comment)
		return
	}

	helpers.ResponseServerErrorJSON(c)
	if errGet != nil {
		fmt.Printf("ERROR in GetComment services: %s", errGet.Error())
	} else {
		fmt.Printf("ERROR in GetComment services: Don't GetComment")
	}
}

// Create func
func (controller CommentController) Create(c *gin.Context) {
	postID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: post_id")
	}
	//check permisson
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permissions error")
		return
	}

	json := models.Comment{}
	if errBind := c.BindJSON(&json); errBind != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "BindJSON: %s\n"+errBind.Error())
		return
	}

	// validation
	if len(json.Message) == 0 {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Missing a few fields:  Message is NULL")
		return
	}
	// validation

	if json.Status == 0 {
		json.Status = 1
	}
	json.UserID = myUserID
	commentID, errCreate := controller.Service.Create(json, postID)
	if errCreate == nil && commentID >= 0 {
		helpers.ResponseSuccessJSON(c, 1, "Create comment successful", map[string]interface{}{"id": commentID})

		// auto Increase Posts
		go func() {
			ok, errIncreasePostComments := controller.Service.IncreasePostComments(postID)
			if errIncreasePostComments != nil {
				fmt.Printf("IncreasePostComments service: %s\n", errIncreasePostComments.Error())
			}
			if ok != true {
				fmt.Printf("IncreasePostComments service")
			}
		}()

		// push noti
		go func() {
			user, _ := services.NewUserService().Get(myUserID)
			writerID, _ := services.NewPostService().GetUserIDByPostID(postID)

			notify := models.Notification{
				UserID:     writerID,
				ObjectID:   postID,
				ObjectType: "post",
				Title:      "@" + user.Username + " bình luận bài đăng của bạn",
				Message:    json.Message,
			}
			PushTest(notify)

		}()
		return
	}
	helpers.ResponseServerErrorJSON(c)
	if errCreate != nil {
		fmt.Printf("Create services: %s\n", errCreate.Error())
	} else {
		fmt.Printf("Create services: Don't Create")
	}
}

// Update func
func (controller CommentController) Update(c *gin.Context) {

}

// Delete func
func (controller CommentController) Delete(c *gin.Context) {
	commentID, errParseInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errParseInt != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid comment id")
		return
	}

	//check exist
	exist, errCheckExistComment := controller.Service.CheckExistComment(commentID)
	if errCheckExistComment != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("CheckExistComment service: %s\n", errCheckExistComment.Error())
		return
	}
	if exist != true {
		helpers.ResponseNotFoundJSON(c, 2, "No exist this object")
		return
	}

	writerID, errGetUserIDByComment := controller.Service.GetUserIDByComment(commentID)
	if errGetUserIDByComment != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetUserIDByComment service: %s\n", errGetUserIDByComment.Error())
		return
	}

	postID, errGetPostIDbyComment := controller.Service.GetPostIDbyComment(commentID)
	if errGetPostIDbyComment != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetPostIDbyComment service: %s\n", errGetPostIDbyComment.Error())
		return
	}

	//check permisson
	if id, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token")); writerID != id || errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permissions error")
		return
	}

	deleted, errDelete := controller.Service.Delete(commentID)
	if errDelete == nil && deleted == true {
		helpers.ResponseNoContentJSON(c)

		// auto Decrease Status Comments
		go func() {
			ok, errDecreasePostComments := controller.Service.DecreasePostComments(postID)
			if errDecreasePostComments != nil {
				fmt.Printf("DecreasePostComments service: %s\n", errDecreasePostComments.Error())
			}
			if ok != true {
				fmt.Printf("DecreasePostComments service")
			}
		}()
		return
	}

	helpers.ResponseServerErrorJSON(c)
	if errDelete != nil {
		fmt.Printf("Delete services: %s", errDelete.Error())
	} else {
		fmt.Printf("Delete services: Don't delete comment")
	}
}
