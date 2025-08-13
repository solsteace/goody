package payload

type Auth interface {
	GetUserId() uint
}

type AuthPayload struct {
	UserId  uint `json:"UserId"`
	IsAdmin bool `json:"IsAdmin"`
}

func NewAuth(userId uint, isAdmin bool) AuthPayload {
	return AuthPayload{
		UserId:  userId,
		IsAdmin: isAdmin}
}

func (ap AuthPayload) GetUserId() uint {
	return ap.UserId
}
