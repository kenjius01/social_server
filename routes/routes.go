package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kenjius01/social-sever/controllers"
	"github.com/kenjius01/social-sever/middlewares"
)

func Setup(app *fiber.App) {

	// Auth endpoint
	auth := app.Group("/auth")
	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)
	auth.Post("/logout", controllers.Logout)

	//User enpoint
	user := app.Group("/user")
	user.Get("/", controllers.GetAllUsers)
	user.Get("/:id", controllers.GetUser)
	user.Put("/:id", middlewares.VerifyUser, controllers.UpdateUser)
	user.Delete("/:id", middlewares.VeryfyAdmin, controllers.DeleteUser)
	user.Post("/follow/:id", controllers.FollowUser)
	user.Delete("/unfollow/:id", controllers.Unfollow)
	user.Get("/follower/:id", controllers.GetAllFollowers)
	user.Get("/following/:id", controllers.GetAllFollwing)
	user.Get("/follow/:userId", controllers.GetNumFollow)

	//Post enpoint
	post := app.Group("/post")
	post.Post("/", middlewares.VerifyUserPost, controllers.CreatePost)
	post.Get("/", controllers.GetAllPost)
	post.Get("/:id", controllers.GetPost)
	post.Put("/:id", middlewares.VerifyUserPost, controllers.UpdatePost)
	post.Delete("/:id", middlewares.VerifyUserPost, controllers.DeletePost)
	post.Post("/:id/like", controllers.LikePost)
	post.Get("/:userId/timeline", controllers.GetPostTimeLine)

}
