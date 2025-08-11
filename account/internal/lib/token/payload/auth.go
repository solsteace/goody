package payload

type Auth interface {
	GetUserId() uint
}

type AuthPayload struct {
	UserId uint `json:"UserId"`
}

func NewAuth(userId uint) AuthPayload {
	return AuthPayload{
		UserId: userId,
	}
}

func (ap AuthPayload) GetUserId() uint {
	return ap.UserId
}
