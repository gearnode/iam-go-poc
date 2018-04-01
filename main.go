package main

import (
	"fmt"
	"regexp"
	"strings"
	//  "database/sql"
	_ "github.com/lib/pq"
)

type Organization struct {
	ID     string
	Prefix string
	Email  string
	Name   string
}

func (organization *Organization) ORN() *ORN {
	return &ORN{
		Namespace:    "jws",
		Application:  "organizations",
		OwnerID:      fmt.Sprintf("%s", organization.ID),
		ResourceType: "account",
		ResourcePath: fmt.Sprintf("%s/%s", organization.Prefix, organization.ID),
	}
}

func test() {
	org := Organization{
		ID:     "101977606264",
		Prefix: "o-jw41v5abna",
		Email:  "sample@domain.com",
	}
	fmt.Println("Object Resource Name =>", org.ORN().String())
}

type ORN struct {
	Namespace    string
	Application  string
	ResourceType string
	ResourcePath string
	OwnerID      string
}

func (orn *ORN) String() string {
	return fmt.Sprintf(
		"orn:%s:%s::%s:%s/%s",
		orn.Namespace,
		orn.Application,
		orn.OwnerID,
		orn.ResourceType,
		orn.ResourcePath,
	)
}

type Statement struct {
	ID       ORN
	Resource string
	Slug     string
	Actions  []string
}

func Authorize(identity ORN, resource string, action string) bool {
	statement := Statement{
		ID: ORN{
			Namespace:    "coretech",
			Application:  "iam",
			ResourceType: "statement",
			ResourcePath: "policy/someAcl",
		},
		Resource: "orn:campus-management:cockpit:*:action/*",
		Actions:  []string{"cockpit:showAction"},
		Slug:     "policy/someAcl",
	}
	statements := []Statement{statement}
	for i := range statements {
		if isPermitedResource(statements[i], resource) &&
			isPermitedAction(statements[i], action) {
			return true
		}
	}
	return false
}

func isPermitedResource(statement Statement, resource string) bool {
	expr := strings.Replace(statement.Resource, "*", ".*", -1)
	matched, err := regexp.MatchString(expr, resource)
	if err == nil && matched {
		return true
	}
	return false
}

func isPermitedAction(statement Statement, action string) bool {
	// Use HashSet to have O(1) func
	for i := range statement.Actions {
		if statement.Actions[i] == action {
			return true
		}
	}
	return false
}

func main() {
	test()
	// db, err := sql.Open("postgres", "postgres://postgres@localhost/iam")
	// if err != nil { panic(err) }

	// rows, err := db.Query("SELECT * FROM organizations")
	// if err != nil { panic(err) }
	// defer db.Close()

	// x := Authorize(
	// 	ORN{},
	// 	"orn:campus-management:cockpit:93724:action/55",
	// 	"cockpit:showAction",
	// )
	// fmt.Println("DEBUG –– %t", x)
}
