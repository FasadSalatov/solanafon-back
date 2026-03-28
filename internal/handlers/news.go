package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fasad/solanafon-back/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// NewsHandler handles /api/news/* endpoints
type NewsHandler struct {
	db *gorm.DB
}

func NewNewsHandler(db *gorm.DB) *NewsHandler {
	return &NewsHandler{db: db}
}

// GetFeed — GET /api/news
func (h *NewsHandler) GetFeed(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	var total int64
	h.db.Model(&models.NewsPost{}).Count(&total)

	var posts []models.NewsPost
	h.db.Preload("App").Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts)

	result := make([]fiber.Map, len(posts))
	for i, p := range posts {
		var isLiked bool
		h.db.Model(&models.NewsLike{}).Where("post_id = ? AND user_id = ?", p.ID, userID).
			Select("count(*) > 0").Scan(&isLiked)

		result[i] = fiber.Map{
			"id": fmt.Sprintf("post_%d", p.ID), "appId": fmt.Sprintf("app_%d", p.AppID),
			"appName": p.App.Title, "appIcon": p.App.IconURL,
			"text": p.Text, "imageUrl": p.ImageURL,
			"commentsCount": p.CommentsCount, "likesCount": p.LikesCount,
			"sharesCount": p.SharesCount, "isLiked": isLiked,
			"createdAt": p.CreatedAt,
		}
	}

	return c.JSON(fiber.Map{"success": true, "posts": result, "pagination": fiber.Map{"page": page, "total": total}})
}

// LikePost — POST /api/news/:postId/like
func (h *NewsHandler) LikePost(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	postID, _ := strconv.Atoi(c.Params("postId"))

	var existing models.NewsLike
	if h.db.Where("post_id = ? AND user_id = ?", postID, userID).First(&existing).Error == nil {
		// Unlike
		h.db.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&models.NewsLike{})
		h.db.Model(&models.NewsPost{}).Where("id = ?", postID).UpdateColumn("likes_count", gorm.Expr("likes_count - 1"))
		var post models.NewsPost
		h.db.First(&post, postID)
		return c.JSON(fiber.Map{"success": true, "liked": false, "likesCount": post.LikesCount})
	}

	h.db.Create(&models.NewsLike{PostID: uint(postID), UserID: userID, CreatedAt: time.Now()})
	h.db.Model(&models.NewsPost{}).Where("id = ?", postID).UpdateColumn("likes_count", gorm.Expr("likes_count + 1"))
	var post models.NewsPost
	h.db.First(&post, postID)
	return c.JSON(fiber.Map{"success": true, "liked": true, "likesCount": post.LikesCount})
}

// SharePost — POST /api/news/:postId/share
func (h *NewsHandler) SharePost(c *fiber.Ctx) error {
	postID, _ := strconv.Atoi(c.Params("postId"))
	h.db.Model(&models.NewsPost{}).Where("id = ?", postID).UpdateColumn("shares_count", gorm.Expr("shares_count + 1"))
	var post models.NewsPost
	h.db.First(&post, postID)
	return c.JSON(fiber.Map{"success": true, "sharesCount": post.SharesCount})
}

// GetComments — GET /api/news/:postId/comments
func (h *NewsHandler) GetComments(c *fiber.Ctx) error {
	postID := c.Params("postId")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset := (page - 1) * limit

	var comments []models.NewsComment
	h.db.Where("post_id = ?", postID).Preload("User").Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&comments)

	result := make([]fiber.Map, len(comments))
	for i, cm := range comments {
		result[i] = fiber.Map{
			"id": fmt.Sprintf("comment_%d", cm.ID), "userId": fmt.Sprintf("user_%d", cm.UserID),
			"userName": cm.User.GetDisplayName(), "userAvatar": cm.User.GetAvatarURL(),
			"text": cm.Text, "createdAt": cm.CreatedAt,
		}
	}
	return c.JSON(fiber.Map{"success": true, "comments": result})
}

// PostComment — POST /api/news/:postId/comments
func (h *NewsHandler) PostComment(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	postID, _ := strconv.Atoi(c.Params("postId"))

	var input struct {
		Text string `json:"text"`
	}
	if err := c.BodyParser(&input); err != nil || input.Text == "" {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "VALIDATION_ERROR", "message": "text is required"}})
	}

	comment := models.NewsComment{PostID: uint(postID), UserID: userID, Text: input.Text}
	h.db.Create(&comment)
	h.db.Model(&models.NewsPost{}).Where("id = ?", postID).UpdateColumn("comments_count", gorm.Expr("comments_count + 1"))

	return c.Status(201).JSON(fiber.Map{"success": true, "comment": fiber.Map{
		"id": fmt.Sprintf("comment_%d", comment.ID), "text": comment.Text, "createdAt": comment.CreatedAt,
	}})
}

// CreatePost — POST /api/developer/apps/:appId/news
func (h *NewsHandler) CreatePost(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	appID, _ := strconv.Atoi(c.Params("appId"))

	var app models.MiniApp
	if err := h.db.Where("id = ? AND creator_id = ?", appID, userID).First(&app).Error; err != nil {
		return c.Status(403).JSON(fiber.Map{"error": fiber.Map{"code": "FORBIDDEN", "message": "Not your app"}})
	}

	var input struct {
		Text     string `json:"text"`
		ImageURL string `json:"imageUrl"`
	}
	if err := c.BodyParser(&input); err != nil || input.Text == "" {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "VALIDATION_ERROR", "message": "text is required"}})
	}

	post := models.NewsPost{AppID: uint(appID), Text: input.Text, ImageURL: input.ImageURL}
	h.db.Create(&post)
	return c.Status(201).JSON(fiber.Map{"success": true, "post": fiber.Map{
		"id": fmt.Sprintf("post_%d", post.ID), "text": post.Text, "createdAt": post.CreatedAt,
	}})
}
