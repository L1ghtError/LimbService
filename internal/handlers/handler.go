package handlers

import (
	"bufio"
	"light-backend/internal/auth"
	"light-backend/internal/middleware"
	"light-backend/internal/ports"
	"light-backend/internal/validation"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"light-backend/internal/domain/media"
	"light-backend/internal/domain/token"
	"light-backend/internal/domain/user"
)

type HttpServer struct {
	UserRepo  user.Repository
	TokenRepo token.Repository
	MediaRepo media.Repository

	MediaTransport media.Transport
}

func (h *HttpServer) GoogleOAuth(c *fiber.Ctx) error {
	config := auth.ConfigGoogle()
	url := config.AuthCodeURL("state")
	return c.Redirect(url)
}

func (h *HttpServer) GoogleOAuthCb(c *fiber.Ctx, params ports.GoogleOAuthCbParams) error {
	t, err := auth.ConfigGoogle().Exchange(c.Context(), params.Code)
	if err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}
	usr, err := auth.GetGoogleResponse(t.AccessToken)
	if err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}

	dbUser, err := h.UserRepo.GetUserByEmail(c.UserContext(), usr.Email)
	if err != nil {
		if err != user.ErrNotFound {
			return err
		}

		dbUser, err = h.UserRepo.Register(
			c.UserContext(),
			user.UserSchema{
				Email:    usr.Email,
				UserName: usr.UserName,
				Fullname: usr.Fullname,
			},
			true,
		)
		if err != nil {
			return err
		}
	}

	tokens, err := token.GenerateTokens(dbUser.ID, dbUser.Email)
	if err != nil {
		return &fiber.Error{Code: fiber.ErrInternalServerError.Code, Message: err.Error()}
	}

	err = h.TokenRepo.SaveToken(c.UserContext(), token.TokenSchema{UserId: dbUser.ID, RefreshToken: tokens.Refresh})
	if err != nil {
		return &fiber.Error{Code: fiber.ErrInternalServerError.Code, Message: err.Error()}
	}

	c.Cookie(&fiber.Cookie{Name: middleware.CookieJWT, Value: tokens.Refresh,
		Expires: time.Now().Add(token.RefreshokenExpires), SessionOnly: false})

	return c.SendStatus(fiber.StatusCreated)
}

func (h *HttpServer) Register(c *fiber.Ctx) error {

	myValidator := validation.XValidator{Validator: validator.New()}

	var userInput ports.RegisterInput
	if err := c.BodyParser(&userInput); err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}

	if errs := myValidator.Validate(userInput); len(errs) > 0 && errs[0].Error {

		return validation.GenerateErrorResp(&errs)
	}

	dbUser, err := h.UserRepo.Register(
		c.UserContext(),
		user.UserSchema{
			Email:    string(userInput.Email),
			UserName: userInput.Username,
			Password: []byte(userInput.Password),
			Fullname: userInput.Fullname,
		},
		false,
	)

	if err != nil {
		if err == user.ErrAlreadyExists {
			return &fiber.Error{Code: fiber.ErrConflict.Code, Message: err.Error()}
		}
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}

	tokens, err := token.GenerateTokens(dbUser.ID, dbUser.Email)
	if err != nil {
		return &fiber.Error{Code: fiber.ErrInternalServerError.Code, Message: err.Error()}
	}

	err = h.TokenRepo.SaveToken(c.UserContext(), token.TokenSchema{UserId: dbUser.ID, RefreshToken: tokens.Refresh})
	if err != nil {
		return &fiber.Error{Code: fiber.ErrInternalServerError.Code, Message: err.Error()}
	}

	c.Cookie(&fiber.Cookie{Name: middleware.CookieJWT, Value: tokens.Refresh,
		HTTPOnly: true, Expires: time.Now().Add(token.RefreshokenExpires), SessionOnly: false})
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"token": tokens.Access, "email": dbUser.Email})
}

func (h *HttpServer) Login(c *fiber.Ctx) error {
	myValidator := validation.XValidator{Validator: validator.New()}

	var userInput ports.LoginInput
	if err := c.BodyParser(&userInput); err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}

	if errs := myValidator.Validate(userInput); len(errs) > 0 && errs[0].Error {

		return validation.GenerateErrorResp(&errs)
	}

	dbUser, err := h.UserRepo.GetUserByEmail(c.UserContext(), string(userInput.Email))
	if err != nil {
		if err == user.ErrNotFound {
			return &fiber.Error{Code: fiber.ErrNotFound.Code, Message: err.Error()}
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(userInput.Password))
	if err != nil {
		return err
	}

	tokens, err := token.GenerateTokens(dbUser.ID, dbUser.Email)
	if err != nil {
		return err
	}

	err = h.TokenRepo.SaveToken(c.UserContext(), token.TokenSchema{UserId: dbUser.ID, RefreshToken: tokens.Refresh})
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{Name: middleware.CookieJWT, Value: tokens.Refresh,
		HTTPOnly: true, Expires: time.Now().Add(token.RefreshokenExpires), SessionOnly: false})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": tokens.Access, "email": dbUser.Email})
}

func (h *HttpServer) Logout(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	// TODO Refactor: instead of using Getenv we should rely on config.Get, but for now its low priority
	claims, err := token.ClaimModel(userToken.Raw, []byte(os.Getenv("JWT_REFRESH_SECRET")))
	if err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}

	err = h.TokenRepo.RemoveToken(c.UserContext(), token.TokenSchema{UserId: claims.UserId})
	if err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok", "message": "ok"})
}

func (h *HttpServer) Activate(c *fiber.Ctx) error {

	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"status": "error", "message": "Not Implemented"})
}

func (h *HttpServer) Refresh(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)

	// TODO Refactor: instead of using Getenv we should rely on config.Get, but for now its low priority
	claims, err := token.ClaimModel(userToken.Raw, []byte(os.Getenv("JWT_REFRESH_SECRET")))
	if err != nil {
		return err
	}

	dbToken, err := h.TokenRepo.GetToken(c.UserContext(), claims.UserId)
	if err != nil {
		return err
	} else if dbToken.RefreshToken != userToken.Raw {
		return &fiber.Error{Code: fiber.ErrUnauthorized.Code, Message: fiber.ErrUnauthorized.Error()}
	}

	tokens, err := token.GenerateTokens(claims.UserId, claims.Email)
	if err != nil {
		return err
	}

	err = h.TokenRepo.SaveToken(c.UserContext(), token.TokenSchema{UserId: claims.UserId, RefreshToken: tokens.Refresh})
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{Name: middleware.CookieJWT, Value: tokens.Refresh,
		HTTPOnly: true, Expires: time.Now().Add(token.RefreshokenExpires)})
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"token": tokens.Access})
}

func (h *HttpServer) GetUserInfo(c *fiber.Ctx, username string) error {
	_ = username

	userToken := c.Locals("user").(*jwt.Token)

	// TODO Refactor: instead of using Getenv we should rely on config.Get, but for now its low priority
	claims, _ := token.ClaimModel(userToken.Raw, []byte(os.Getenv("JWT_ACCESS_SECRET")))

	user, err := h.UserRepo.GetUserById(c.UserContext(), claims.UserId)
	if err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}

	// TODO: consider returning dedicated model instead of modifying existing one
	user.Password = nil

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *HttpServer) GetAvailableProcessors(c *fiber.Ctx) error {

	info, err := h.MediaTransport.GetWorkersInfo()
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(ports.AvailableProcessorsSchema{AvailableProcessors: info.AvailableProcessors})
}

func (h *HttpServer) UploadImage(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)

	// TODO Refactor: instead of using Getenv we should rely on config.Get, but for now its low priority
	claims, _ := token.ClaimModel(userToken.Raw, []byte(os.Getenv("JWT_ACCESS_SECRET")))

	user, err := h.UserRepo.GetUserById(c.UserContext(), claims.UserId)
	if err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}

	file, err := c.FormFile("document")
	if err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}

	imgId, err := h.MediaRepo.UploadPicture(c.UserContext(), *file, media.ImageMetadata{UserId: user.ID, Header: file.Header})
	if err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}

	err = h.UserRepo.AddImageId(c.UserContext(), user, imgId)
	if err != nil {
		return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: err.Error()}
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"imageid": imgId})
}

// TODO: allow to download only for owner or allower users
func (h *HttpServer) DownloadImage(c *fiber.Ctx, imageId ports.ImageId) error {
	myValidator := validation.XValidator{Validator: validator.New()}
	type ImageInput struct {
		ImageId string `params:"id" validate:"required,len=24"`
	}

	body := ImageInput{
		ImageId: string(imageId),
	}

	if errs := myValidator.Validate(body); len(errs) > 0 && errs[0].Error {

		return validation.GenerateErrorResp(&errs)
	}

	fstream, file, err := h.MediaRepo.DownloadPicture(c.UserContext(), body.ImageId)
	if err != nil {
		if err == fiber.ErrBadRequest {
			return &fiber.Error{Code: fiber.ErrBadRequest.Code, Message: fiber.ErrBadRequest.Error()}
		}
		return &fiber.Error{Code: fiber.ErrInternalServerError.Code, Message: fiber.ErrInternalServerError.Error()}
	}
	// TODO: sets the same mime info that users passes, potential vulnerability
	for key := range file.Metadata.Header {
		c.Set(key, file.Metadata.Header.Get(key))

	}
	return c.SendStream(fstream)
}

func (h *HttpServer) ProcessImage(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	validator := validation.XValidator{Validator: validator.New()}
	var task media.ImageTask

	if err := c.BodyParser(&task); err != nil {
		return err
	}
	if errs := validator.Validate(task); len(errs) > 0 && errs[0].Error {
		err := validation.GenerateErrorResp(&errs)
		return err
	}

	register := func(w *bufio.Writer) {
		h.MediaTransport.ProcessImage(task, w)
	}

	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(register)
	return nil
}
