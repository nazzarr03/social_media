package comment

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/social-media/config"
	"github.com/nazzarr03/social-media/models"
	"github.com/nazzarr03/social-media/utils"
)

func CreateCommentToPost(c *fiber.Ctx) error {
	var request CreateCommentRequest
	var comment models.Comment
	var post models.Post
	postId := c.Params("postID")
	err := config.Db.Model(&models.Post{}).Where("post_id = ?", postId).First(&post).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	utils.DTOtoJSON(request, comment)
	comment.PostID = utils.StringToInt(postId)
	comment.UserID = utils.StringToInt(c.Locals("userID").(string))

	file, err := c.FormFile("image")
	if err == nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Image not found",
		})
	}

	if file != nil {
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
		comment.ImageURL = resp.SecureURL
	}

	comment.CreatedAt = time.Now()
	if err := config.Db.Create(&comment).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Comment created successfully",
	})
}

func CreateCommentToComment(c *fiber.Ctx) error {
	var request CreateCommentRequest
	var comment models.Comment
	var parentComment models.Comment
	var post models.Post
	parentCommentId := c.Params("commentID")
	postId := c.Params("postID")

	err := config.Db.Model(&models.Comment{}).Where("comment_id = ?", parentCommentId).First(&parentComment).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Comment not found",
		})
	}

	err = config.Db.Model(&models.Post{}).Where("post_id = ?", postId).First(&post).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	utils.DTOtoJSON(request, comment)
	comment.PostID = utils.StringToInt(postId)
	comment.UserID = utils.StringToInt(c.Locals("userID").(string))
	comment.ParentCommentID = &parentComment.CommentID

	file, err := c.FormFile("image")
	if err == nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Image not found",
		})
	}

	if file != nil {
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
		comment.ImageURL = resp.SecureURL
	}

	comment.CreatedAt = time.Now()
	if err := config.Db.Save(&comment).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Comment created successfully",
	})
}
