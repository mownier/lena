package errors

import (
	"fmt"
	"strings"
)

const (
	ErrCodeUserAlreadyExists = iota + 1
	ErrCodeUserDoesNotExist
	ErrCodeSessionAlreadyExists
	ErrCodeSessionDoesNotExist
	ErrCodeGettingAccessToken
	ErrCodeSigningOut
	ErrCodeVerifyingAccessToken
	ErrCodeMetadataNotOkay
	ErrCodeAuthorizationNotSet
	ErrCodeEmptyAuthorization
	ErrCodeRegistering
	ErrCodeSigningIn
	ErrCodeRefreshingAccessToken
	ErrCodeOpeningSqliteDB
	ErrCodeMigratingSqliteDB
	ErrCodeCreatingMigrationsTable
	ErrCodeQueryingMigrations
	ErrCodeCannotBeginDBTx
	ErrCodeExecutingMigrationCmd
	ErrCodeMigrationModify
	ErrCodeInsertingMigration
	ErrCodeDBTxCommitHasFailed
	ErrCodeRowScanHasFailed
	ErrCodeUserCreationTimeCannotBeDetermined
	ErrCodeCheckingIfUserExists
	ErrCodeQueryingUsers
	ErrCodeSessionCreationTimeCannotBeDetermined
	ErrCodeSessionArchivedTimeCannotBeDetermined
	ErrCodeSessionAcccessTokenExpiryCannotBeDetermined
	ErrCodeSessionRefreshTokenExpiryCannotBeDetermined
	ErrCodeGettingSession
	ErrCodeExecutingSessionUpdateStmt
	ErrCodePreparingSessionUpdateStmt
	ErrCodeNoSessionToUpdate
	ErrCodeInsertingSession
	ErrCodeSessionAlreadyInvalidated
	ErrCodeInvalidRefreshToken
	ErrCodeInvalidPassword
	ErrCodeHashingPassword
	ErrCodeSessionCanNoLongerBeExtended
	ErrCodeSessionExpired
	ErrCodeHTTPMethodNotAllowed
	ErrCodeHTTPBodyCannotBeRead
	ErrCodeHTTPBodyMalformed
	ErrCodeGeneratingResponse
)

var errMessages = map[int]string{
	ErrCodeUserAlreadyExists:                           "user already exists",
	ErrCodeUserDoesNotExist:                            "user does not exist",
	ErrCodeSessionAlreadyExists:                        "session already exists",
	ErrCodeSessionDoesNotExist:                         "session does not exist",
	ErrCodeGettingAccessToken:                          "problem on getting access token",
	ErrCodeSigningOut:                                  "problem on signing out",
	ErrCodeVerifyingAccessToken:                        "problem on verifying access token",
	ErrCodeMetadataNotOkay:                             "metadata not okay",
	ErrCodeAuthorizationNotSet:                         "authorization is not set",
	ErrCodeEmptyAuthorization:                          "authorization is empty",
	ErrCodeRegistering:                                 "problem on registering",
	ErrCodeSigningIn:                                   "problem on signing in",
	ErrCodeRefreshingAccessToken:                       "problem on refreshing access token",
	ErrCodeOpeningSqliteDB:                             "problem opening sqlite db",
	ErrCodeMigratingSqliteDB:                           "problem on migrating sqlite db",
	ErrCodeCreatingMigrationsTable:                     "problem on creating migrations table",
	ErrCodeQueryingMigrations:                          "problem on querying migrations",
	ErrCodeCannotBeginDBTx:                             "problem on starting db transaction",
	ErrCodeExecutingMigrationCmd:                       "problem on executing migration command",
	ErrCodeMigrationModify:                             "problem on modify action of migration",
	ErrCodeInsertingMigration:                          "problem on inserting migration version",
	ErrCodeDBTxCommitHasFailed:                         "db transaction commit has failed",
	ErrCodeRowScanHasFailed:                            "row scan has failed",
	ErrCodeUserCreationTimeCannotBeDetermined:          "can not determine user creation time",
	ErrCodeCheckingIfUserExists:                        "problem on checking if user exists",
	ErrCodeQueryingUsers:                               "problem on querying users",
	ErrCodeSessionCreationTimeCannotBeDetermined:       "can not determine session creation time",
	ErrCodeSessionArchivedTimeCannotBeDetermined:       "can not determine session archived time",
	ErrCodeSessionAcccessTokenExpiryCannotBeDetermined: "can not determine session access token expiry time",
	ErrCodeSessionRefreshTokenExpiryCannotBeDetermined: "can not determine session refresh token expiry time",
	ErrCodeGettingSession:                              "problem on getting session",
	ErrCodeExecutingSessionUpdateStmt:                  "problem on executing update statement of session",
	ErrCodePreparingSessionUpdateStmt:                  "problem on prepating update statement of session",
	ErrCodeNoSessionToUpdate:                           "no session to update",
	ErrCodeInsertingSession:                            "problem on inserting session",
	ErrCodeSessionAlreadyInvalidated:                   "session is already invalidated",
	ErrCodeInvalidRefreshToken:                         "refresh token is invalid",
	ErrCodeInvalidPassword:                             "password is invalid",
	ErrCodeHashingPassword:                             "problem on hashing password",
	ErrCodeSessionCanNoLongerBeExtended:                "session can no longer be extended",
	ErrCodeSessionExpired:                              "session has expired",
	ErrCodeHTTPMethodNotAllowed:                        "method not allowed",
	ErrCodeHTTPBodyCannotBeRead:                        "body of http request cannot be read",
	ErrCodeHTTPBodyMalformed:                           "body of http request is malformed",
	ErrCodeGeneratingResponse:                          "problem on generating response",
}

var userFriendlyErrMessages = map[int]string{
	ErrCodeRegistering:           "Something went wrong while registering the user",
	ErrCodeUserAlreadyExists:     "The name of the user to be registered is already used",
	ErrCodeSigningIn:             "Something went wrong while signing in the user",
	ErrCodeInvalidPassword:       "Password does not match with the name of the user",
	ErrCodeUserDoesNotExist:      "The user is not yet registered",
	ErrCodeGettingAccessToken:    "Something went wrong while validating the authentication of the user",
	ErrCodeSigningOut:            "Something went wrong while signing out the user",
	ErrCodeVerifyingAccessToken:  "Something went wrong while verifying the authentication of the user",
	ErrCodeRefreshingAccessToken: "Something went wrong while refreshing the authentication of the user",
	ErrCodeHTTPMethodNotAllowed:  "Method of the request is not allowed",
	ErrCodeHTTPBodyCannotBeRead:  "Something went wrong while processing body of the request",
	ErrCodeHTTPBodyMalformed:     "The content of the request's body is malformed",
	ErrCodeGeneratingResponse:    "Unable to generate response",
}

type UserFriendlyResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type AppError struct {
	Domain  string
	Code    int
	Message string
	Other   error
}

func (e *AppError) ContainsCode(code int) (*AppError, bool) {
	if code == e.Code {
		return e, true
	}
	other := e.Other
	for other != nil {
		if appError, ok := other.(*AppError); ok {
			if appError.Code == code {
				return appError, true
			}
			other = appError.Other
		} else {
			return nil, false
		}
	}
	return nil, false
}

func (e *AppError) AsUserFriendlyResponse() UserFriendlyResponse {
	message, exists := userFriendlyErrMessages[e.Code]
	if !exists {
		message = "An unknown error has occurred"
	}
	return UserFriendlyResponse{Code: e.Code, Message: message}
}

func (e *AppError) Error() string {
	stack := []string{}
	format := "error domain=%s code=%d message=%s"
	stack = append(stack, fmt.Sprintf(format, e.Domain, e.Code, e.Message))
	other := e.Other
	for other != nil {
		if appError, ok := other.(*AppError); ok {
			stack = append(stack, fmt.Sprintf("error domain=%s code=%d message=%s", appError.Domain, appError.Code, appError.Message))
			other = appError.Other
		} else {
			stack = append(stack, fmt.Sprintf("error unsupported type=%T, info=%v", other, other))
			break
		}
	}
	return strings.Join(stack, "\n")
}

func NewAppError(code int, domain string, other error) *AppError {
	message, exists := errMessages[code]
	if !exists {
		message = "unsupported error code"
	}
	return &AppError{
		Domain:  domain,
		Code:    code,
		Message: message,
		Other:   other,
	}
}
