package controllers

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kenjius01/social-sever/database"
	"github.com/kenjius01/social-sever/models"
)

type Post struct {
	ID        int    `json:"id"`
	UserId    int    `json:"userId"`
	Desc      string `json:"desc"`
	Image     string `json:"image"`
	Like      []int  `json:"likes"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func CreateResponsePost(postModel models.Post) Post {
	return Post{
		ID:        postModel.ID,
		UserId:    postModel.UserId,
		Desc:      postModel.Desc,
		Like:      []int{},
		Image:     postModel.Image,
		CreatedAt: postModel.CreatedAt,
		UpdatedAt: postModel.UpdatedAt}
}

func CreatePost(c *fiber.Ctx) error {
	var post models.Post
	if err := c.BodyParser(&post); err != nil {
		return c.Status(400).JSON(err.Error())

	}
	if err := database.DB.Create(&post).Error; err != nil {
		return c.Status(402).JSON(err.Error())
	}
	resPost := CreateResponsePost(post)
	database.DB.Model(&models.Like{}).Select("user_id").Find(&resPost.Like, "post_id = ?", post.ID)
	return c.Status(200).JSON(resPost)
}

//Get a post
func FindPost(id int, post *models.Post) error {
	database.DB.First(&post, "id = ?", id)
	if post.ID == 0 {
		return errors.New("Can not find the post!")
	}
	return nil
}

func GetPost(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	post := models.Post{}
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	if err := FindPost(id, &post); err != nil {
		return c.Status(404).JSON(err.Error())
	}

	return c.Status(200).JSON(post)
}

//Get all post
func GetAllPost(c *fiber.Ctx) error {
	posts := []models.Post{}
	if err := database.DB.Find(&posts).Error; err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON(posts)
}

//Update Post
func UpdatePost(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	post := models.Post{}
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	if err := FindPost(id, &post); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	updatePost := models.Post{}
	if err := c.BodyParser(&updatePost); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	if post.UserId != updatePost.UserId {
		return c.Status(400).JSON("You can update only your post")
	}

	database.DB.Model(&post).Updates(models.Post{
		Desc:  updatePost.Desc,
		Image: updatePost.Image,
	})
	return c.Status(200).JSON(updatePost)
}

func DeletePost(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	post := models.Post{}

	if err := FindPost(id, &post); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	deletePost := models.Post{}
	if err := c.BodyParser(&deletePost); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if post.UserId != deletePost.UserId {
		return c.Status(400).JSON("You can delete only your post")
	}
	if err := database.DB.Delete(&post).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}
	return c.Status(200).JSON("The post has been deleted!")
}

func LikePost(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")

	userLike := models.Like{}
	if err := c.BodyParser(&userLike); err != nil {
		c.Status(400).JSON(err.Error())
	}
	user := models.User{}
	if err := findUser(userLike.UserId, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	if id != userLike.PostId {
		return c.Status(402).JSON("Action forbidden!")

	}
	checkLike := models.Like{}
	database.DB.Where("post_id = ? AND user_id = ?", id, userLike.UserId).First(&checkLike)
	if checkLike.ID == 0 {
		database.DB.Create(&userLike)
		return c.Status(200).JSON("Liked!")

	} else {
		database.DB.Delete(&checkLike)
		return c.Status(200).JSON("unlike!")
	}

}

// Get TimeLine Post
func GetPostTimeLine(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("userId")

	posts := []models.Post{}
	if err := database.DB.Raw("? UNION ?",
		database.DB.Model(&models.Post{}).Joins("inner join followers f on f.user_id = posts.user_id AND f.follower_id = ?", id),
		database.DB.Model(&models.Post{}).Order("updated_at desc").Where("user_id = ?", id),
	).Scan(&posts).Error; err != nil {
		return c.Status(400).JSON(err.Error())
	}
	// result := []int{}

	resPosts := []Post{}
	for _, post := range posts {
		responsePost := CreateResponsePost(post)
		database.DB.Model(&models.Like{}).Select("user_id").Find(&responsePost.Like, "post_id = ?", post.ID)
		resPosts = append(resPosts, responsePost)
	}

	return c.Status(200).JSON(resPosts)

}
