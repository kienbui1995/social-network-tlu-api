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

// UserController controller
type UserController struct {
	Service services.UserServiceInterface
}

// GetAll func
func (controller UserController) GetAll(c *gin.Context) {

	params := helpers.ParamsGetAll{}
	params.Skip, _ = strconv.Atoi(c.DefaultQuery("skip", "0"))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", "25"))
	users, errGetAll := controller.Service.GetAll(params)
	if errGetAll != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("GetAll service: %s\n", errGetAll.Error())
		return
	}
	helpers.ResponseEntityListJSON(c, 1, "Get user list successful", users, nil, len(users))
}

// Get func
func (controller UserController) Get(c *gin.Context) {
	userID := c.Param("id")
	user, _ := controller.Service.Get(userID)
	// Valida se existe a pessoa (404)
	if user.IsEmpty() {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not Found User")
		return
	}

	helpers.ResponseSuccessJSON(c, 1, "Info of user", user)
}

// Delete func
func (controller UserController) Delete(c *gin.Context) {
	userID := c.Param("id")
	user, _ := controller.Service.Get(userID)
	// Valida se existe a pessoa que será excluida (404)
	if user.IsEmpty() {
		helpers.ResponseNotFoundJSON(c, configs.EcNoExistObject, "Not Found User")
		return
	}
	// Valida se deu erro ao tentar excluir (500)
	if _, err := controller.Service.Delete(userID); err != nil {
		helpers.ResponseServerErrorJSON(c)
		fmt.Printf("Delete service: %s\n", err)
		return
	}
	helpers.ResponseNoContentJSON(c)
}

// Create func
func (controller UserController) Create(c *gin.Context) {
	var user models.User
	// BadRequest (400)
	if err := c.BindJSON(&user); err != nil {
		helpers.ResponseBadRequestJSON(c, configs.EcParam, "message")
		return
	}
	// Valida Invalid Entity (422)
	if _, err := user.Validate(); err != nil {
		c.JSON(422, gin.H{"errors": err})
		return
	}
	// Valida se deu erro ao inserir (500)
	// if err := controller.Repository.Create(&user); err != nil {
	// 	c.JSON(500, gin.H{"errors": "Houve um erro no servidor"})
	// 	return
	// }

	// c.JSON(201, gin.H{"person": person})
}

// func (controller PersonController) Update(c *gin.Context) {
// 	personId := c.Param("id")
// 	person := controller.Repository.Get(personId)
// 	// Valida se existe a pessoa que será editada (404)
// 	if person.IsEmpty() {
// 		c.JSON(404, gin.H{"errors": "Registros não encontrado."})
// 		return
// 	}
// 	// Valida BadRequest (400)
// 	if err := c.BindJSON(&person); err != nil {
// 		c.JSON(400, gin.H{"errors: ": err.Error()})
// 		return
// 	}
// 	// Valida Invalid Entity (422)
// 	if err := person.Validate(); err == nil {
// 		c.JSON(422, gin.H{"errors": err})
// 		return
// 	}
// 	// Valida se deu erro ao inserir (500)
// 	if err := controller.Repository.Update(&person); err != nil {
// 		c.JSON(500, gin.H{"errors": "Houve um erro no servidor."})
// 		return
// 	}
//
// 	c.JSON(201, gin.H{"person": person})
// }
