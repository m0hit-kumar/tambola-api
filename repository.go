package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/m0hit-kumar/tambola/dto"
	"github.com/m0hit-kumar/tambola/models"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

// Books   --------------- start-------------------
func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}

	err := context.BodyParser(&book)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
		return err
	}
	err = r.DB.Create(&book).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not create book"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Book created successfully"})

	return nil
}

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Books{}
	err := r.DB.Find(bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Book not found"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Book fetched", "data": bookModels})
	return nil
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookModel := models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Id cannot be empty"})
		return nil
	}
	err := r.DB.Delete(bookModel, id)
	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not delete book"})
		return err.Error
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "book deleted succesfully"})

	return nil
}

func (r *Repository) GetBookByID(context *fiber.Ctx) error {
	id := context.Params("id")
	bookModel := &models.Books{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Id cannot be empty"})
		return nil
	}
	fmt.Println("the id is ", id)
	err := r.DB.Where("id=?", id).First(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not fetch the book"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Book id feteched correctly", "data": bookModel})
	return nil
}

// Books ------------------end----------------------

// Auth -------------------start---------------------
func (r *Repository) Login(context *fiber.Ctx) error {
	req := dto.LoginReq{}
	user := models.Users{}
	err := context.BodyParser(&req)
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Request Failed"})
		return err
	}
	if req.Username == "" || req.Password == "" {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Field cannot be empty"})
		return nil
	}
	err = r.DB.Where("username = ? AND password = ?", req.Username, req.Password).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Incorrect Username or Password"})
			return nil
		}
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Internal Server Error"})
		return err
	}

	jwt_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Id":        user.ID,
		"Username":  user.Username,
		"ExpiredAt": time.Now().Add(time.Minute * 5).Unix(),
	})
	signedToken, err := jwt_token.SignedString([]byte("my_secret_key"))
	if err != nil {
		log.Fatal(err)
		return err
	}
	login := &dto.LoginRes{
		Username: user.Username,
		Email:    user.Email,
		Token:    signedToken,
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "User Logged in Succefully", "data": login})
	return nil
}

func (r *Repository) SignUp(context *fiber.Ctx) error {
	user := dto.SignUpReq{}

	err := context.BodyParser(&user)
	if user.Username == "" || user.Password == "" || user.Email == "" {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Field cannot be empty"})
		return nil
	}
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "Request failed"})
		return err
	}
	err = r.DB.Create(&user).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not create user"})
		return err
	}
	r.CreateSampleTemplate(context)
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "User created Successfully"})
	return nil
}

// Auth---------------------end-----------------------

// Tickets -----------------start ---------------------
func (r *Repository) GeTTicketTemplate(context *fiber.Ctx) error {
	ticketDesign := dto.TicketDesign{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Id cannot be empty", "data": "s"})
		return nil
	}
	err := r.DB.Where("userID = ?", id).First(&ticketDesign).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Default Template", "data": dto.TicketDesignRes()})
		return nil
	}
	context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Default Template", "data": ticketDesign})
	return nil
}

func (r *Repository) CreateSampleTemplate(context *fiber.Ctx) error {
	err := r.DB.Create(dto.TicketDesignRes()).Error
	if err != nil {
		log.Fatal("cannot create the sample ticket")
		return err
	}
	print("Sample ticket template created succesfully")
	return nil
}

func (r *Repository) CreateTicketTemplate(context *fiber.Ctx) error {
	ticketDesign := models.TicketDesign{}
	err := context.BodyParser(&ticketDesign)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
		return err
	}
	err = r.DB.Create(&ticketDesign).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not Ticket Template"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Ticket Template created successfully"})

	return nil
}
