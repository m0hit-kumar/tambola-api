package main

import (
	"log"
	"net/http"
	"strconv"
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

//---------------- Test server is running ----------------
func (r *Repository) Test(context *fiber.Ctx) error {
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Server Running successfully"})
	return nil
}

//-----------------Test server is running------------------


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
 	user := models.Users{}
 	err := context.BodyParser(&user)
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Username/Email or Password cannot be empty"})
		return err
	}
    if (user.Username == "" && user.Email == "") || user.Password == "" {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Field cannot be empty"})
		return nil
	}
	err = r.DB.Where("(username = ? OR email = ?) AND password = ?", user.Username,user.Email, user.Password).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Incorrect Username/Email or Password"})
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
	user.Token = signedToken
	r.DB.Save(&user)
	cookie := new(fiber.Cookie)
	cookie.Name = "token"
	cookie.Value = signedToken
	cookie.Expires = time.Now().Add(12 * time.Minute)
	context.Cookie(cookie)
	login := &dto.LoginRes{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    signedToken,
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "User Logged in Succefully", "data": login})
	return nil
}

func (r *Repository) SignUp(context *fiber.Ctx) error {
	user := models.Users{}

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
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "User created Successfully"})
	return nil
}

// Auth---------------------end-----------------------

// Tickets -----------------start ---------------------
func (r *Repository) GeTTicketTemplate(context *fiber.Ctx) error {
	ticketDesign := models.TicketDesign{}
	defaultTicketDesign := dto.TicketDesignRes()

	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Default Template", "data": defaultTicketDesign})
		return nil
	}
	userid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Fatal("Error:", err)
		return err
	}
	userID := uint(userid)
	ticketDesign = models.TicketDesign{
		HostName:   defaultTicketDesign.HostName,
		Background: defaultTicketDesign.Background,
		Text:       defaultTicketDesign.Text,
		UserID:     userID,
		Border:     defaultTicketDesign.Border,
	}

	result := r.DB.FirstOrCreate(&ticketDesign, models.TicketDesign{UserID: userID})
	if result.Error != nil {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Internal Server Error"})
		return result.Error
	}
	response := dto.TicketDesign{
		HostName:   ticketDesign.HostName,
		Background: ticketDesign.Background,
		Border:     ticketDesign.Border,
		Text:       ticketDesign.Text,
	}

	if result.RowsAffected == 1 {
		context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Default Template", "data": response})
		return nil
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "User Saved Template", "data": response})
	return nil
}

func (r *Repository) CreateTicketTemplate(context *fiber.Ctx) error {
	cookie := context.Cookies("token")
	if cookie == "" {
		context.Status(http.StatusUnauthorized).JSON(&fiber.Map{"message": "User Must Logged In To perform this task"})
		return nil
	}

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
