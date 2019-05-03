package main

import (
	elastic "gopkg.in/olivere/elastic.v3"

	"fmt"
	"reflect"
	"regexp"
)

const (
	TYPER_USER = "user"
)

var (
	usernamePattern = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
)

type User struct {
	Username string `json:"username"`
	Password string `jsont:"password"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
}

func checkUser(username, password string) bool {
	esClient, err := elastic.NewClient(elastic.SetURL(ES_URL), elastic.SetSniff(false))
	if err != nil {
		fmt.Printf("ES is not setup %v\n", err)
		return false
	}

	//search in a term query
	termQuery := elastic.NewTermQuery("username", username)
	queryResult, err := esClient.Search().
		Index(INDEX).
		Query(termQuery).
		Pretty(true).
		Do()
	if err != nil {
		fmt.Printf("ES query failed %v\n", err)
		return false
	}

	var tyu User
	for _, item := range queryResult.Each(reflect.TypeOf(tyu)) {
		u := item.(User)
		return u.Password == password && u.Username == username
	}
	return false

}

func addUser(user User) bool {
	esClient, err := elastic.NewClient(elastic.SetURL(ES_URL), elastic.SetSniff(false))
	if err != nil {
		fmt.Printf("ES is not setup %v\n", err)
		return false
	}
	termQeury := elastic.NewTermQuery("username", user.Username)
	queryResult, err := esClient.Search().
		Index(INDEX).
		Query(termQeury).
		Pretty(true).
		Do()
	if err != nil {
		fmt.Printf("ES query failed %v\n", err)
		return false
	}

	if queryResult.TotalHits() > 0 {
		fmt.Printf("User %s already existed\n", user.Username)
		return false
	}

	_, errr := esClient.Index().
		Index(INDEX).
		Type(TYPER_USER).
		Id(user.Username).
		BodyJson(user).
		Refresh(true).
		Do()
	if errr != nil {
		fmt.Printf("ES save user failed %v\n", errr)
		return false
	}
	return true
}
