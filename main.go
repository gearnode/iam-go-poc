package main

import (
	"fmt"
	"regexp"
	"strings"
)

type ORN struct {
	Namespace    string
	Application  string
	ResourceType string
	ResourcePath string
	OwnerID      string
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
	for i := range statement.Actions {
		if statement.Actions[i] == action {
			return true
		}
	}
	return false
}

func main() {
	x := Authorize(
		ORN{},
		"orn:campus-management:cockpit:93724:action/55",
		"cockpit:showAction",
	)
	fmt.Println("DEBUG –– %t", x)
}
