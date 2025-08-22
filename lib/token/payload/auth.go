package payload

type Auth struct {
	UserId  uint `json:"UserId"`
	IsAdmin bool `json:"IsAdmin"`
}

func NewAuth(userId uint, isAdmin bool) Auth {
	return Auth{
		UserId:  userId,
		IsAdmin: isAdmin}
}
