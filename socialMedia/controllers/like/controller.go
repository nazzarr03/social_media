package like

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/social-media/config"
	"github.com/nazzarr03/social-media/models"
	"github.com/nazzarr03/social-media/utils"
)

func CreateLikeToPost(c *fiber.Ctx) error {
	var like models.Like
	userID := c.Locals("userID").(int)
	postID := c.Params("postID")
	err := config.Db.Model(&models.Post{}).Where("post_id = ?", postID).First(&models.Post{}).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
	}

	like.IsLiked = true
	like.UserID = userID
	like.PostID = utils.StringToInt(postID)

	if err := config.Db.Create(&like).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Like created successfully",
	})
}

func CreateLikeToComment(c *fiber.Ctx) error {
	var like models.Like
	userID := c.Locals("userID").(int)
	postID := c.Params("postID")
	commentID := c.Params("commentID")
	err := config.Db.Model(&models.Post{}).Where("post_id = ?", postID).First(&models.Post{}).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
	}

	err = config.Db.Model(&models.Comment{}).Where("comment_id = ?", commentID).First(&models.Comment{}).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Comment not found",
		})
	}

	like.IsLiked = true
	like.UserID = userID
	like.PostID = utils.StringToInt(postID)

	if err := config.Db.Create(&like).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Like created successfully",
	})
}
