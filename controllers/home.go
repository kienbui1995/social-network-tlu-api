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

// HomeController controller
type HomeController struct {
	Service services.HomeServiceInterface
}

// FindUser func
func (controller HomeController) FindUser(c *gin.Context) {
	name := c.Query("name")
	if len(name) == 0 {
		helpers.ResponseBadRequestJSON(c, configs.EcParamMissingField, "Missing a few fields: name")
		return
	}

	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permission error")
		return
	}

	users, errFindUserByUsernameAndFullName := controller.Service.FindUserByUsernameAndFullName(name, myUserID)
	if errFindUserByUsernameAndFullName != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("ERROR FindUserByUsernameAndFullName: %s", errFindUserByUsernameAndFullName.Error())
	} else {
		helpers.ResponseEntityListJSON(c, 1, "User list found", users, nil, len(users))
	}
}

// NewsFeed func
func (controller HomeController) NewsFeed(c *gin.Context) {
	//check permisson and get user id
	myUserID, errGetUserIDFromToken := GetUserIDFromToken(c.Request.Header.Get("token"))
	if myUserID < 0 || errGetUserIDFromToken != nil {
		helpers.ResponseAuthJSON(c, configs.EcPermission, "Permissions error")
		return
	}
	var params helpers.ParamsGetAll
	sort := c.DefaultQuery("sort", "pagerank")
	if sort != "pagerank" {
		var errConvertSort error
		sort, errConvertSort = helpers.ConvertSort(sort)
		if errConvertSort != nil {
			helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: sort "+errConvertSort.Error())
			return
		}
	}
	skip, errSkip := strconv.Atoi(c.DefaultQuery("skip", configs.SSkip))
	if errSkip != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: skip "+errSkip.Error())
		return
	}
	limit, errLimit := strconv.Atoi(c.DefaultQuery("limit", configs.SLimit))
	if errLimit != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: limit "+errLimit.Error())
		return
	}
	params.Sort = sort
	params.Skip = skip
	params.Limit = limit
	var posts []models.Post
	var errList error
	if params.Sort == "pagerank" {
		posts, errList = controller.Service.GetNewsFeedWithPageRank(params, myUserID)
	} else {
		posts, errList = controller.Service.GetNewsFeed(params, myUserID)
	}
	if errList == nil {
		helpers.ResponseEntityListJSON(c, 1, "Get news feed successful", posts, nil, len(posts))
		return
	}

	helpers.ResponseServerErrorJSON(c)
	if errList != nil {
		fmt.Printf("GetNewsFeed services: %s", errList.Error())
	} else {
		fmt.Printf("GetNewsFeed services: Don't get news feed")
	}

}
