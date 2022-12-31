package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/kamalshkeir/kenv"
	"github.com/kamalshkeir/kmux"
	"github.com/kamalshkeir/korm"
	"github.com/kamalshkeir/ksbus"
	"github.com/kamalshkeir/pgdriver"
)

type User struct {
	Id               uint      `korm:"pk"`
	Firstname        string    `korm:"text"`
	Lastname         string    `korm:"text"`
	Email            string    `korm:"text"`
	Creationdate     time.Time `korm:"now"`
	Isserviceaccount bool
}

type DbConfig struct {
	Name     string `kenv:"DB_NAME"`
	Host     string `kenv:"DB_HOST"`
	Port     string `kenv:"DB_PORT"`
	Username string `kenv:"DB_USER"`
	Password string `kenv:"DB_PASS"`
}

func main() {
	// load env
	kenv.Load(".env")
	dbConf := &DbConfig{}
	err := kenv.Fill(dbConf)
	if err != nil {
		log.Fatal(err)
	}
	// connect to db
	pgdriver.Use()
	err = korm.New(korm.POSTGRES, dbConf.Name, fmt.Sprintf("%s:%s@%s:%s", dbConf.Username, dbConf.Password, dbConf.Host, dbConf.Port))
	if err != nil {
		log.Fatal(err)
	}
	bus := korm.WithBus(ksbus.NewServer())
	err = bus.App.LocalTemplates("templates")
	if err != nil {
		log.Fatal(err)
	}
	// handle Index
	bus.App.GET("/", func(c *kmux.Context) {
		c.Html("index.html", nil)
	})

	// list all users html page
	bus.App.GET("/user", func(c *kmux.Context) {
		users, err := korm.Table("\"User\"").Query("select * from \"User\"").All()
		if err != nil {
			fmt.Println("error GET:", err)
			c.Html("user.html", nil)
			return
		}
		c.Html("user.html", map[string]any{
			"users": users,
		})
	})

	// delete user given id
	bus.App.DELETE("/user/id:int", func(c *kmux.Context) {
		id := c.Param("id")
		if id == "" {
			c.Status(400).Json(map[string]any{
				"error": "id is empty",
			})
			return
		}

		idToDelete, err := strconv.Atoi(id)
		if err != nil {
			c.Status(400).Json(map[string]any{
				"error": "id not a number",
			})
			return
		}

		_, err = korm.Table("\"User\"").Where("id = ?", idToDelete).Delete()
		if err != nil {
			fmt.Println("error delete:", err)
			c.Status(400).Json(map[string]any{
				"error": err.Error(),
			})
			return
		}

		c.Json(map[string]any{
			"success": id + " deleted",
			"id":      id,
		})
	})

	bus.App.POST("/user", func(c *kmux.Context) {
		body := c.BodyJson()
		userToAdd := map[string]any{}
		if v, ok := body["firstname"]; ok {
			userToAdd["firstname"] = v
		}
		if v, ok := body["lastname"]; ok {
			userToAdd["lastname"] = v
		}
		if v, ok := body["email"]; ok {
			userToAdd["email"] = v
		}

		if len(userToAdd) == 0 {
			c.Status(http.StatusBadRequest).Json(map[string]any{
				"error": "expected user in body to be string not empty",
			})
			return
		}

		_, err = korm.Table("\"User\"").Insert(userToAdd)
		if err != nil {
			fmt.Println("error:", err)
			c.Status(http.StatusBadRequest).Json(map[string]any{
				"error": err.Error(),
			})
			return
		}
		c.Json(map[string]any{
			"success": "done",
		})
	})

	// run server
	bus.Run("localhost:9313")
	// shutdown db
	err = korm.Shutdown("User")
	if err != nil {
		log.Fatal(err)
	}
}
