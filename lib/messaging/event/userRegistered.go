package event

const UserRegisteredName = "user.event.registered"

type UserRegistered struct {
	IdUser uint   `json:"id_user"`
	Nama   string `json:"nama"`
}

func NewUserRegistered(IdUser uint, nama string) UserRegistered {
	return UserRegistered{IdUser: IdUser, Nama: nama}
}
