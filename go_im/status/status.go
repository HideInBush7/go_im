package status

const (
	SUCCESS             = 0
	ERROR               = 500
	USER_NOT_EXIST      = 50001
	USER_ALREADY_EXIST  = 50002
	USER_PASSWORD_WRONG = 50003
	INVALID_TOKEN       = 50004
)

var codeMsg = map[int32]string{
	SUCCESS:             "success.",
	ERROR:               "service error.",
	USER_NOT_EXIST:      "user not exist.",
	USER_ALREADY_EXIST:  "user already exist.",
	USER_PASSWORD_WRONG: "user password wrong.",
	INVALID_TOKEN:       "invalid token.",
}

func Message(code int32) string {
	return codeMsg[code]
}
