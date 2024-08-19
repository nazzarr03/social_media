package user

import (
	"context"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/social-media/config"
	"github.com/nazzarr03/social-media/middleware"
	"github.com/nazzarr03/social-media/models"
	"github.com/nazzarr03/social-media/utils"
	"github.com/streadway/amqp"
)

func Register(c *fiber.Ctx) error {
	var request RegisterRequest
	var user models.User
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	utils.DTOtoJSON(request, user)
	user.Password = hashedPassword
	user.CreatedAt = time.Now()

	if err := config.Db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
	})
}

func Login(c *fiber.Ctx) error {
	var request LoginRequest
	var user models.User
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err := config.Db.Model(&models.User{}).Where("username = ?", request.Username).First(&user).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if err := utils.CheckPassword(request.Password, user.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid password",
		})
	}

	token, err := middleware.GenerateToken(&user, 24)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := config.Rdb.Set(context.Background(), strconv.Itoa(int(user.UserID)), token, 24*time.Hour).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}

func CreateProfileImage(c *fiber.Ctx) error {
	var user models.User
	userID := c.Params("id")
	if err := config.Db.Model(&models.User{}).Where("user_id = ?", userID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	tempFilePath := fmt.Sprintf("./uploads/%s", file.Filename)
	if err := c.SaveFile(file, tempFilePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	resp, err := utils.UploadToCloudinary(config.Cld, tempFilePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	os.Remove(tempFilePath)
	user.ImageURL = resp.SecureURL
	user.UpdatedAt = time.Now()
	if err := config.Db.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ch, err := config.RabbitMQConn.Channel()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"notification",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	email := user.Email
	subject := "Profile Image Uploaded"
	body := fmt.Sprintf("Your profile image has been uploaded successfully. You can view it by visiting this link: %s", user.ImageURL)
	message := fmt.Sprintf("%s\n%s\n\n%s", email, subject, body)
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Image uploaded successfully",
	})
}

func UpdateProfileImage(c *fiber.Ctx) error {
	var user models.User
	userID := c.Params("id")
	if err := config.Db.Model(&models.User{}).Where("user_id = ?", userID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	lastPart := path.Base(user.ImageURL)
	publicID := strings.Split(lastPart, ".")[0]
	if err := utils.DeleteFromCloudinary(config.Cld, publicID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	tempFilePath := fmt.Sprintf("./uploads/%s", file.Filename)
	if err := c.SaveFile(file, tempFilePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	resp, err := utils.UploadToCloudinary(config.Cld, tempFilePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	os.Remove(tempFilePath)
	user.ImageURL = resp.SecureURL
	user.UpdatedAt = time.Now()
	if err := config.Db.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Image updated successfully",
	})
}
