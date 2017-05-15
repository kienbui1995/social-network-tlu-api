package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//const
const (
	GraphFacebookAPI = "https://graph.facebook.com"
)

//VerifyFacebookID func to check logged in via Facebook
func VerifyFacebookID(id string, accessToken string) bool {
	url := fmt.Sprintf("%s/me?fields=id&access_token=%s", GraphFacebookAPI, accessToken)
	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}

	// read json http response
	jsonDataFromHTTP, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		panic(err)
	}

	var jsonData struct {
		ID string `json:"id"`
	}

	err = json.Unmarshal([]byte(jsonDataFromHTTP), &jsonData) // here!

	if err != nil {
		panic(err)
	}
	if jsonData.ID == id {
		return true
	}
	return false
}
