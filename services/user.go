package services

import (
	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/models"
)

// UserServiceInterface include method list
type UserServiceInterface interface {
	GetAll(q helpers.ParamsGetAll) (models.PublicUsers, error)
	Get(id string) (models.User, error)
	Delete(id string) (bool, error)
	Create(p *models.User) (int, error)
	Update(p *models.User) (models.User, error)
}

// UserService struct
type userService struct{}

// NewUserService to constructor
func NewUserService() *userService {
	return new(userService)
}

// GetAll func
func (service userService) GetAll(params helpers.ParamsGetAll) (models.PublicUsers, error) {

	stmt := `
	MATCH (u:User)
	return u {id:ID(u), .*}  AS user

		SKIP {skip}
		LIMIT {limit}
		`
	p := map[string]interface{}{
		"skip":  params.Skip,
		"limit": params.Limit,
	}
	var res []struct {
		User models.PublicUser `json:"user"`
	}

	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: p,
		Result:     &res,
	}

	err := conn.Cypher(&cq)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}
	//fmt.Printf("res: %v", res)
	var list models.PublicUsers
	for _, val := range res {
		list = append(list, val.User)
	}

	return list, nil
}

// Get func
func (service userService) Get(id string) (models.User, error) {

	var user models.User

	return user, nil
}

// Delete func
func (service userService) Delete(id string) (bool, error) {

	return false, nil
}

// Create func
func (service userService) Create(p *models.User) (int, error) {
	//
	// record := models.User{
	// 	Type:             p.Type,
	// 	Name:             p.Name,
	// 	CityName:         p.CityName,
	// 	CompanyName:      p.CompanyName,
	// 	Address:          p.Address,
	// 	Number:           p.Number,
	// 	Complement:       p.Complement,
	// 	District:         p.District,
	// 	Zip:              p.Zip,
	// 	BirthDate:        p.BirthDate,
	// 	Cpf:              p.Cpf,
	// 	Rg:               p.Rg,
	// 	Gender:           p.Gender,
	// 	BusinessPhone:    p.BusinessPhone,
	// 	HomePhone:        p.HomePhone,
	// 	MobilePhone:      p.MobilePhone,
	// 	Cnpj:             p.Cnpj,
	// 	StateInscription: p.StateInscription,
	// 	Phone:            p.Phone,
	// 	Fax:              p.Fax,
	// 	Email:            p.Email,
	// 	Website:          p.Website,
	// 	Observations:     p.Observations,
	// 	RegisteredAt:     time.Now(),
	// 	RegisteredByUUID: p.RegisteredByUUID,
	// }
	//
	// err := db.Create(&record).Error
	// if err != nil {
	// 	log.Print(err.Error())
	// }
	//
	// *(p) = record
	return -1, nil
}

// Update func
func (service userService) Update(p *models.User) (models.User, error) {
	// record := models.User{
	// 	Name:             p.Name,
	// 	CityName:         p.CityName,
	// 	CompanyName:      p.CompanyName,
	// 	Address:          p.Address,
	// 	Number:           p.Number,
	// 	Complement:       p.Complement,
	// 	District:         p.District,
	// 	Zip:              p.Zip,
	// 	BirthDate:        p.BirthDate,
	// 	Cpf:              p.Cpf,
	// 	Rg:               p.Rg,
	// 	Gender:           p.Gender,
	// 	BusinessPhone:    p.BusinessPhone,
	// 	HomePhone:        p.HomePhone,
	// 	MobilePhone:      p.MobilePhone,
	// 	Cnpj:             p.Cnpj,
	// 	StateInscription: p.StateInscription,
	// 	Phone:            p.Phone,
	// 	Fax:              p.Fax,
	// 	Email:            p.Email,
	// 	Website:          p.Website,
	// 	Observations:     p.Observations,
	// }
	//
	// err := db.Model(&models.Person{}).
	// 	Where("uuid = ?", p.UUID).
	// 	Updates(&record).Error
	//
	// if err != nil {
	// 	log.Print(err.Error())
	// }
	return models.User{}, nil
}
