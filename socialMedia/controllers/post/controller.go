package post

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/social-media/config"
	"github.com/nazzarr03/social-media/models"
	"github.com/nazzarr03/social-media/utils"
)

func CreatePost(c *fiber.Ctx) error {
	var postRequest CreatePostRequest
	var post models.Post
	userId := c.Params("userID")
	err := config.Db.Model(&models.User{}).Where("user_id = ?", userId).First(&models.User{}).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if err := c.BodyParser(&postRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	utils.DTOtoJSON(postRequest, post)
	post.UserID = utils.StringToInt(userId)
	post.CreatedAt = time.Now()
	if err := config.Db.Create(&post).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	shortKey := utils.GenerateShortKey()
	longURL := fmt.Sprintf("http://postgres:8081/posts/%s", strconv.Itoa(post.PostID))
	_, err = url.ParseRequestURI(longURL)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL",
		})
	}

	if err := config.Rdb.Set(context.Background(), shortKey, longURL, 0).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Post created successfully",
	})
}

func CreateImageByPostID(c *fiber.Ctx) error {
	var post models.Post
	postId := c.Params("postID")
	err := config.Db.Model(&models.Post{}).Where("post_id = ?", postId).First(&post).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
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
	post.ImageURL = resp.SecureURL
	post.UpdatedAt = time.Now()
	if err := config.Db.Save(&post).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Post image uploaded successfully",
	})
}

func UpdatePost(c *fiber.Ctx) error {
	var updatePostRequest UpdatePostRequest
	var post models.Post
	postID := c.Params("postID")
	if err := config.Db.Model(&models.Post{}).Where("post_id = ?", postID).First(&post).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
	}

	if post.CreatedAt.Add(5 * time.Minute).Before(time.Now()) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You can only update your post within 5 minutes of creation",
		})
	}

	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	utils.DTOtoJSON(updatePostRequest, post)
	post.UpdatedAt = time.Now()
	if err := config.Db.Save(&post).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Post updated successfully",
	})
}

func UpdatePostImage(c *fiber.Ctx) error {
	var post models.Post
	postID := c.Params("postID")
	err := config.Db.Model(&models.Post{}).Where("post_id = ?", postID).First(&post).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
	}

	if post.CreatedAt.Add(5 * time.Minute).Before(time.Now()) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You can only update your post image within 5 minutes of creation",
		})
	}

	lastPart := path.Base(post.ImageURL)
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
	post.ImageURL = resp.SecureURL
	post.UpdatedAt = time.Now()
	if err := config.Db.Save(&post).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Post image updated successfully",
	})
}

func DeletePost(c *fiber.Ctx) error {
	var post models.Post
	postID := c.Params("postID")
	err := config.Db.Model(&models.Post{}).Where("post_id = ?", postID).First(&post).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
	}

	lastPart := path.Base(post.ImageURL)
	publicID := strings.Split(lastPart, ".")[0]
	if err := utils.DeleteFromCloudinary(config.Cld, publicID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := config.Db.Delete(&post).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Post deleted successfully",
	})
}

func GetPostShortURL(c *fiber.Ctx) error {
	var post models.Post
	postID := c.Params("postID")
	err := config.Db.Model(&models.Post{}).Where("post_id = ?", postID).First(&post).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
	}

	shortKey := utils.GenerateShortKey()
	longURL := fmt.Sprintf("http://postgres:8081/posts/%s", postID)
	_, err = url.ParseRequestURI(longURL)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL",
		})
	}

	if err := config.Rdb.Set(context.Background(), shortKey, longURL, 0).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"post": post,
	})
}
