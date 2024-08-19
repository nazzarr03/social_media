package friendship

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/social-media/config"
	"github.com/nazzarr03/social-media/models"
)

func CreateFriendship(c *fiber.Ctx) error {
	user := models.User{}
	friend := models.User{}
	userID := c.Locals("userID").(string)
	err := config.Db.Model(&models.User{}).Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	friendID := c.Params("friendID")
	err = config.Db.Model(&models.User{}).Where("user_id = ?", friendID).First(&friend).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Friend not found",
		})
	}

	friendship := models.Friendship{
		UserID:   user.UserID,
		FriendID: friend.UserID,
		IsActive: true,
	}

	if err := config.Db.Create(&friendship).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Friendship created successfully",
	})
}
