package core

import "fmt"

const (
	ErrorTeamExists  string = "TEAM_EXISTS"
	ErrorPRExists    string = "PR_EXISTS"
	ErrorPRMerged    string = "PR_MERGED"
	ErrorNotAssigned string = "NOT_ASSIGNED"
	ErrorNoCandidate string = "NO_CANDIDATE"
	ErrorNotFound    string = "NOT_FOUND"
	ErrorUserExists  string = "USER_EXISTS"
)

func Throw(code string, msg string) error {
	return fmt.Errorf("%s: %s", code, msg)
}
