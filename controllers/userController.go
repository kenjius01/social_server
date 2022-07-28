package controllers

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kenjius01/social-sever/database"
	"github.com/kenjius01/social-sever/models"
)

type User struct {
	ID           int    `json:"id"`
	FirstName    string `json:"firstName"`
	Username     string `json:"username"`
	LastName     string `json:"lastName"`
	Avatar       string `json:"avatar"`
	CoverImage   string `json:"coverImg"`
	Description  string `json:"desc"`
	Address      string `json:"address"`
	WorkAt       string `json:"workAt"`
	Relationship string `json:"relationship"`
}

type Follower struct {
	ID           int  `json:"id"`
	UserInfo     User `json:"userInfo"`
	FollowerInfo User `json:"followerInfo"`
	UserId       int  `json:"userId"`
	FollowerId   int  `json:"followerId"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func CreateResponseUser(userModel models.User) User {
	return User{ID: userModel.ID,
		Username:     userModel.Username,
		FirstName:    userModel.FirstName,
		LastName:     userModel.LastName,
		Avatar:       userModel.Avatar,
		CoverImage:   userModel.CoverImage,
		Description:  userModel.Description,
		Address:      userModel.Address,
		WorkAt:       userModel.WorkAt,
		Relationship: userModel.Relationship}
}

func CreateResponseFollowUser(followUser models.Follower, user User, follower User) Follower {
	return Follower{
		ID:           followUser.ID,
		UserId:       followUser.UserId,
		FollowerId:   followUser.FollowerId,
		UserInfo:     user,
		FollowerInfo: follower,
		CreatedAt:    followUser.CreatedAt,
		UpdatedAt:    followUser.UpdatedAt}
}

func findUser(id int, user *models.User) error {
	database.DB.First(&user, "id=?", id)
	if user.ID == 0 {
		return errors.New("User does not exits!")
	}
	return nil
}

func GetAllUsers(c *fiber.Ctx) error {
	users := []models.User{}
	if err := database.DB.Find(&users).Error; err != nil {
		return c.Status(400).JSON(err.Error())
	}
	responseUsers := []User{}
	for _, user := range users {
		responseUser := CreateResponseUser(user)
		responseUsers = append(responseUsers, responseUser)
	}
	return c.Status(200).JSON(responseUsers)

}

func GetUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	user := models.User{}
	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer!")
	}
	if err := findUser(id, &user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	responseUser := CreateResponseUser(user)
	return c.Status(200).JSON(responseUser)
}

func UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	user := models.User{}
	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer!")
	}
	if err := findUser(id, &user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	type UpdateUser struct {
		ID           int    `json:"id"`
		FirstName    string `json:"firstName"`
		LastName     string `json:"lastName"`
		Avatar       string `json:"avatar"`
		CoverImage   string `json:"coverImg"`
		Description  string `json:"desc"`
		Address      string `json:"address"`
		WorkAt       string `json:"workAt"`
		Relationship string `json:"relationship"`
	}
	var updateData UpdateUser
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	database.DB.Model(&user).Updates(models.User{
		ID:           updateData.ID,
		FirstName:    updateData.FirstName,
		LastName:     updateData.LastName,
		Avatar:       updateData.Avatar,
		CoverImage:   updateData.CoverImage,
		Description:  updateData.Description,
		Address:      updateData.Address,
		WorkAt:       updateData.WorkAt,
		Relationship: updateData.Relationship})
	return c.Status(200).JSON(updateData)
}

// func ChangePassword(c *fiber.Ctx) error {
// 	id, err := c.ParamsInt("id")
// 	user := models.User{}
// 	if err != nil {
// 		return c.Status(400).JSON("Please ensure that id is an integer!")
// 	}
// 	if err := findUser(id, &user); err != nil {
// 		return c.Status(400).JSON(fiber.Map{
// 			"message": err.Error(),
// 		})
// 	}
// }

func DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	user := models.User{}
	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer!")
	}
	if err := findUser(id, &user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}
	return c.Status(200).JSON("user " + user.Username + " has been deleted")

}

// Follow a User

func FollowUser(c *fiber.Ctx) error {
	userFollow := models.Follower{}
	if err := c.BodyParser(&userFollow); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	id, err := c.ParamsInt("id")
	user := models.User{}
	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer!")
	}
	if err := findUser(id, &user); err != nil {
		return c.Status(404).JSON(err.Error())
	}
	follower := models.User{}
	if err := findUser(userFollow.FollowerId, &follower); err != nil {
		return c.Status(402).JSON(err.Error())
	}

	if id == follower.ID {
		return c.Status(403).JSON("You cannot follow yourseft!")
	}
	checkDupplicate := models.Follower{}
	database.DB.Where("user_id = ? AND follower_id = ?", id, follower.ID).Find(&checkDupplicate)
	if checkDupplicate.ID != 0 {
		return c.Status(400).JSON("You cannot follower one person twice!")
	}
	userFollow.UserId = user.ID

	if err := database.DB.Create(&userFollow).Error; err != nil {
		return c.Status(403).JSON(err.Error())
	}
	// responseUser := CreateResponseUser(user)
	// responseFollower := CreateResponseUser(follower)
	// userFollowResponse := CreateResponseFollowUser(userFollow, responseUser, responseFollower)

	return c.Status(200).JSON("Followed success!")

}

func Unfollow(c *fiber.Ctx) error {
	userFollow := models.Follower{}
	if err := c.BodyParser(&userFollow); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	id, err := c.ParamsInt("id")
	user := models.User{}
	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer!")
	}
	if err := findUser(id, &user); err != nil {
		return c.Status(404).JSON(err.Error())
	}

	if id == userFollow.FollowerId {
		return c.Status(403).JSON("You cannot Unfollow yourseft!")
	}

	checkDupplicate := models.Follower{}
	database.DB.Where("user_id = ? AND follower_id = ?", id, userFollow.FollowerId).Find(&checkDupplicate)
	if checkDupplicate.ID == 0 {
		return c.Status(400).JSON("You didn't follow this user!")
	}

	if err := database.DB.Delete(&checkDupplicate).Error; err != nil {
		return c.Status(403).JSON(err.Error())
	}
	// responseUser := CreateResponseUser(user)
	// responseFollower := CreateResponseUser(follower)
	// userFollowResponse := CreateResponseFollowUser(userFollow, responseUser, responseFollower)

	return c.Status(200).JSON("Unfollowed success!")
}

// Get All Follwers

func GetAllFollowers(c *fiber.Ctx) error {
	type ListFollower struct {
		Id        string `json:"id"`
		Username  string `json:"username"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Avatar    string `json:"avatar"`
	}
	result := []ListFollower{}
	id, err := c.ParamsInt("id")
	user := models.User{}
	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer!")
	}
	if err := findUser(id, &user); err != nil {
		return c.Status(404).JSON(err.Error())
	}
	database.DB.Model(&models.User{}).Select("users.id, users.username,users.first_name,users.last_name,users.avatar").Joins("left join followers on users.id = followers.follower_id").Where("	user_id = ?", id).Scan(&result)
	return c.Status(200).JSON(result)

}

func GetAllFollwing(c *fiber.Ctx) error {

	result := []int{}
	id, err := c.ParamsInt("id")
	user := models.User{}
	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer!")
	}
	if err := findUser(id, &user); err != nil {
		return c.Status(404).JSON(err.Error())
	}
	database.DB.Model(&models.Follower{}).Select("follower_id").Find(&result, "user_id = ?", id)
	return c.Status(200).JSON(result)

}

func GetNumFollow(c *fiber.Ctx) error {

	type NumFollow struct {
		NumFollower  int `json:"numberFollower"`
		NumFollowing int `json:"numberFollowing"`
	}
	id, err := c.ParamsInt("userId")
	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer!")
	}

	result := NumFollow{}

	database.DB.Model(&models.Follower{}).Select("count(follower_id)").Where("user_id = ?", id).Find(&result.NumFollower)
	database.DB.Model(&models.Follower{}).Select("count(user_id)").Where("follower_id = ?", id).Find(&result.NumFollowing)
	return c.Status(200).JSON(result)
}
