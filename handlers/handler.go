package handlers

import (
	"light-backend/auth"
	"light-backend/config"
	"light-backend/middleware"
	"light-backend/model"
	"light-backend/service"
	"light-backend/validation"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(c *fiber.Ctx) error {
	config := auth.ConfigGoogle()
	url := config.AuthCodeURL("state")
	return c.Redirect(url)
}

func Callback(c *fiber.Ctx) error {
	token, err := auth.ConfigGoogle().Exchange(c.Context(), c.FormValue("code"))
	if err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}
	user, err := auth.GetGoogleResponse(token.AccessToken)
	if err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}

	tokens, err := service.OAuthConnect(c, model.UserSchema{Email: user.Email,
		UserName: user.UserName, Fullname: user.Fullname, IsActivated: user.Verified})
	if err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}

	c.Cookie(&fiber.Cookie{Name: middleware.CookieJWT, Value: tokens.Refresh,
		Expires: time.Now().Add(service.RefreshokenExpires), SessionOnly: false})

	return c.SendStatus(fiber.StatusCreated)
}

func Registration(c *fiber.Ctx) error {

	myValidator := validation.XValidator{Validator: validator.New()}
	type RegistrationInput struct {
		Email    string `json:"email" validate:"required,email,min=3"`
		UserName string `json:"username" validate:"required,min=3,max=50"`
		Password string `json:"password" validate:"required,min=8,max=72"`
		Fullname string `json:"fullname" validate:"required,min=3,max=50"`
	}

	user := new(RegistrationInput)
	if err := c.BodyParser(user); err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}

	if errs := myValidator.Validate(user); len(errs) > 0 && errs[0].Error {

		return validation.GenerateErrorResp(&errs)
	}

	tokens, err := service.Register(c, model.UserSchema{Email: user.Email, UserName: user.UserName,
		Password: []byte(user.Password), IsActivated: false, Fullname: user.Fullname})

	if err != nil {
		return &fiber.Error{Code: fiber.ErrInternalServerError.Code, Message: err.Error()}
	}

	c.Cookie(&fiber.Cookie{Name: middleware.CookieJWT, Value: tokens.Refresh,
		HTTPOnly: true, Expires: time.Now().Add(service.RefreshokenExpires), SessionOnly: false})
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"token": tokens.Access, "email": user.Email})
}

func Login(c *fiber.Ctx) error {
	myValidator := validation.XValidator{Validator: validator.New()}
	type LoginInput struct {
		Email    string `json:"email" validate:"required,email,min=3"`
		Password string `json:"password" validate:"required,min=8,max=72"`
	}

	user := new(LoginInput)
	if err := c.BodyParser(user); err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}

	if errs := myValidator.Validate(user); len(errs) > 0 && errs[0].Error {

		return validation.GenerateErrorResp(&errs)
	}

	tokens, err := service.Login(c, model.UserSchema{Email: user.Email, Password: []byte(user.Password)})
	if err != nil {
		return &fiber.Error{Code: fiber.ErrInternalServerError.Code, Message: err.Error()}
	}

	c.Cookie(&fiber.Cookie{Name: middleware.CookieJWT, Value: tokens.Refresh,
		HTTPOnly: true, Expires: time.Now().Add(service.RefreshokenExpires), SessionOnly: false})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": tokens.Access, "email": user.Email})
}

func Logout(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	claims, err := service.ClaimModel(&userToken.Raw, []byte(config.Config("JWT_REFRESH_SECRET")))
	if err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}
	err = service.Logout(c, claims)
	if err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}

	return c.SendStatus(fiber.StatusOK)
}

func Activate(c *fiber.Ctx) error {

	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"status": "error", "message": "Not Implemented"})
}

func Refresh(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)

	tokens, err := service.Refresh(c, userToken)
	if err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}

	c.Cookie(&fiber.Cookie{Name: middleware.CookieJWT, Value: tokens.Refresh,
		HTTPOnly: true, Expires: time.Now().Add(service.RefreshokenExpires)})
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"token": tokens.Access})
}

func GetBasics(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	user, err := service.GetBasics(c, userToken)
	if err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}
	return c.Status(fiber.StatusOK).JSON(user)
}
func UploadImage(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	user, err := service.GetBasics(c, userToken)
	if err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}

	file, err := c.FormFile("document")
	if err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}

	imgId, err := service.UploadPicture(c, file, &model.ImageMetadata{UserId: user.ID, Header: file.Header})
	if err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}

	err = service.PushImage(c, user, imgId)
	if err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"token": imgId.Hex()})
}

// TODO: allow to download only for owner or allower users
func DownloadImage(c *fiber.Ctx) error {
	myValidator := validation.XValidator{Validator: validator.New()}
	type ImageInput struct {
		ImageId string `json:"imageid" validate:"required,len=24"`
	}
	body := new(ImageInput)
	if err := c.BodyParser(body); err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}

	if errs := myValidator.Validate(body); len(errs) > 0 && errs[0].Error {

		return validation.GenerateErrorResp(&errs)
	}

	fstream, file, err := service.DownloadPictureSt(c, &body.ImageId)
	if err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}
	// TODO: sets the same mime info that users passes, potential vulnerability
	for key := range file.Metadata.Header {
		c.Set(key, file.Metadata.Header.Get(key))

	}
	return c.SendStream(fstream)
}
