package internal

import (
	"net/http"
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
	reqPayload := new(loginPayload)
	if err := c.BodyParser(reqPayload); err != nil {
		return err
	}

	resData, err := ac.service.login(reqPayload.Phone, reqPayload.Password)
	if err != nil {
		return c.SendString(err.Error())
	}

	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"status":  true,
			"message": "Succeed to POST data",
			"errors":  nil,
			"data":    resData,
		})
}

func (ac AuthController) register(c *fiber.Ctx) error {
	reqPayload := new(registerPayload)
	if err := c.BodyParser(reqPayload); err != nil {
		return err
	}

	birthdate, err := time.Parse("02/01/2006", reqPayload.Birthdate)
	if err != nil {
		return err
	}

	err = ac.service.register(
		reqPayload.Name,
		reqPayload.Password,
		reqPayload.Phone,
		birthdate,
		reqPayload.Occupation,
		reqPayload.Email,
		reqPayload.ProvinceId,
		reqPayload.CityId)
	if err != nil {
		return err
	}

	return c.
		Status(http.StatusCreated).
		SendString("Register Succeed")
}
