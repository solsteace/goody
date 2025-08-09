package internal

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type loginPayload struct {
	Password string `json:"kata_sandi"`
	Phone    string `json:"no_telp"`
}

type registerPayload struct {
	Name       string `json:"nama"`
	Password   string `json:"kata_sandi"`
	Phone      string `json:"no_telp"`
	Birthdate  string `json:"tanggal_lahir"`
	Occupation string `json:"pekerjaan"`
	Email      string `json:"email"`
	ProvinceId string `json:"id_provinsi"`
	CityId     string `json:"id_kota"`
}

type AuthController struct {
	service AuthService
}

func NewAuthController(as AuthService) AuthController {
	return AuthController{
		service: as,
	}
}

func (ac AuthController) login(c *fiber.Ctx) error {
	payload := new(loginPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	err := ac.service.login(payload.Phone, payload.Password)
	if err != nil {
		c.SendString(err.Error())
		return err
	}

	c.SendString("Login OK!")
	return nil
}

func (ac AuthController) register(c *fiber.Ctx) error {
	payload := new(registerPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	birthdate, err := time.Parse("02/01/2006", payload.Birthdate)
	if err != nil {
		return err
	}

	newUserId, err := ac.service.register(
		payload.Name,
		payload.Password,
		payload.Phone,
		birthdate,
		payload.Occupation,
		payload.Email,
		payload.ProvinceId,
		payload.CityId)
	if err != nil {
		return err
	}

	c.SendString(fmt.Sprintf("Your user id %d", newUserId))
	return nil
}
