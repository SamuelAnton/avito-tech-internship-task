package handlers

import (
	"PR_reviewer_assign_service/internal/errors"
	"PR_reviewer_assign_service/internal/models"
	"regexp"
)

// Regular Expressions to check input parameters
var (
	TeamNameRegExp *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9_\.\-\s]*$`)
	UserIDRegExp   *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9\-]*$`)
	UserNameRegExp *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z]*$`)
	PRIDRegExp     *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9\-]*$`)
	PRNameRegExp   *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9_\.\-\s]*$`)
)

// validation for string field
func validateStringField(field string, reg *regexp.Regexp, min_len, max_len int) (errors.ErrorCode, string) {
	if len(field) < min_len {
		return errors.ErrorCodeInvalidInput, " is too small"
	}
	if len(field) > max_len {
		return errors.ErrorCodeInvalidInput, " is too long"
	}

	if !reg.MatchString(field) {
		return errors.ErrorCodeInvalidInput, " is invalid"
	}
	return "", ""
}

/*
	 Team {
		TeamName : string
		members : []TeamMember - {
				UserID : string
				Username : string
				IsActive : boolean
			}
	}
*/
func ValidateTeam(team models.Team) (errors.ErrorCode, string) {
	// Check teamName
	if err, msg := validateStringField(team.TeamName, TeamNameRegExp, 1, 100); err != "" {
		return err, "Team name " + team.TeamName + msg
	}
	// Check members
	if len(team.Members) == 0 {
		return errors.ErrorCodeInvalidInput, "team must have at least one member"
	}
	userIds := make(map[string]bool)
	for _, member := range team.Members {
		// Check member id
		if err, msg := validateStringField(member.UserID, UserIDRegExp, 1, 50); err != "" {
			return err, "Member ID " + member.UserID + msg
		}
		// Check member name
		if err, msg := validateStringField(member.Username, UserNameRegExp, 2, 50); err != "" {
			return err, "Member name " + member.Username + msg
		}
		// Check duplicates
		if userIds[member.UserID] {
			return errors.ErrorCodeInvalidInput, "Duplicate user ID: " + member.UserID
		}
	}
	return "", ""
}

/*
	 TeamNameQuery {
		TeamName : string
	}
*/
func ValidateTeamNameQuery(name models.TeamNameQuery) (errors.ErrorCode, string) {
	// Check teamName
	if err, msg := validateStringField(name.TeamName, TeamNameRegExp, 1, 100); err != "" {
		return err, "Team name " + name.TeamName + msg
	}
	return "", ""
}

/*
	UserActiveQuery {
		UserID : string
		IsActive : boolean
	}
*/
func ValidateUserActiveQuery(user models.UserActiveQuery) (errors.ErrorCode, string) {
	// Check user id
	if err, msg := validateStringField(user.UserID, UserIDRegExp, 1, 50); err != "" {
		return err, "User ID " + user.UserID + msg
	}
	return "", ""
}

/*
	PullRequestCreateQuery {
		PullRequestID : string
		PullRequestName : string
		AuthorID : string
	}
*/
func ValidatePullRequestCreateQuery(pr models.PullRequestCreateQuery) (errors.ErrorCode, string) {
	// Check pr id
	if err, msg := validateStringField(pr.PullRequestID, PRIDRegExp, 1, 50); err != "" {
		return err, "Pull Request ID " + pr.PullRequestID + msg
	}
	// Check pr name
	if err, msg := validateStringField(pr.PullRequestName, PRNameRegExp, 1, 150); err != "" {
		return err, "Pull Request name " + pr.PullRequestName + msg
	}
	// Check author id
	if err, msg := validateStringField(pr.AuthorID, UserIDRegExp, 1, 50); err != "" {
		return err, "Author ID " + pr.AuthorID + msg
	}
	return "", ""
}

/*
	PullRequestMergeQuery {
		PullRequestID : string
	}
*/
func ValidatePullRequestMergeQuery(pr models.PullRequestMergeQuery) (errors.ErrorCode, string) {
	// Check pr id
	if err, msg := validateStringField(pr.PullRequestID, PRIDRegExp, 1, 50); err != "" {
		return err, "Pull Request ID" + pr.PullRequestID + msg
	}
	return "", ""
}

/*
	PullRequestReassignQuery {
		PullRequestID : string
		OldUserID : string
	}
*/
func ValidatePullRequestReassignQuery(pr models.PullRequestReassignQuery) (errors.ErrorCode, string) {
	// Check pr id
	if err, msg := validateStringField(pr.PullRequestID, PRIDRegExp, 1, 50); err != "" {
		return err, "Pull Request ID " + pr.PullRequestID + msg
	}
	// Check user id
	if err, msg := validateStringField(pr.OldUserID, UserIDRegExp, 1, 50); err != "" {
		return err, "Old user's ID" + pr.OldUserID + msg
	}

	return "", ""
}

/*
	UserIDQuery {
		UserID : string
	}
*/
func ValidateUserIDQuery(user models.UserIDQuery) (errors.ErrorCode, string) {
	// Check user id
	if err, msg := validateStringField(user.UserID, UserIDRegExp, 1, 50); err != "" {
		return err, "User ID " + user.UserID + msg
	}

	return "", ""
}
