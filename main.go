package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	c "go.mod/controllers"
)

func main() {
	e := echo.New()

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/users", echo.HandlerFunc(c.GetAllUsers), c.Authenticate)
	e.POST("/users", echo.HandlerFunc(c.InsertUser), c.Authenticate)
	e.PUT("/users/:userID", echo.HandlerFunc(c.UpdateUser), c.Authenticate)
	e.DELETE("/users/:userID", echo.HandlerFunc(c.DeleteUser), c.Authenticate)

	e.POST("/login", c.Login)
	e.POST("/logout", c.Logout)

	s := http.Server{
		Addr:    ":8080",
		Handler: e,
	}

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Println("Connected to port 8080")
	} else {
		log.Fatal(err)
	}
}
