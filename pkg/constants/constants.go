package constants

const (
	EmptyString   = ""
	Authorization = "authorization"
)

type ErrorDefinition struct {
	Code    int
	Message string
}

var (
	NotFound = ErrorDefinition{
		Code:    40005,
		Message: "Not Found",
	}
	WrongPassword = ErrorDefinition{
		Code:    40006,
		Message: "Wrong Password",
	}
	InvalidEmail = ErrorDefinition{
		Code:    40008,
		Message: "Invalid Email",
	}
	PasswordTooShort = ErrorDefinition{
		Code:    40009,
		Message: "Password too short",
	}
	NewPasswordMismatch = ErrorDefinition{
		Code:    40010,
		Message: "New password and confirm password mismatch",
	}
	UserAlreadyExists = ErrorDefinition{
		Code:    40011,
		Message: "User already exists",
	}
	InvalidUUID = ErrorDefinition{
		Code:    40012,
		Message: "Invalid UUID",
	}
	OldPasswordMismatch = ErrorDefinition{
		Code:    40013,
		Message: "Old password mismatch",
	}
	RoleEmpty = ErrorDefinition{
		Code:    40014,
		Message: "Role cannot be empty",
	}
)
