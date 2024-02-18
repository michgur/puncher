package app

import (
	"fmt"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/michgur/puncher/app/db"
	"github.com/michgur/puncher/app/design"
	"github.com/michgur/puncher/app/model"
	"github.com/michgur/puncher/app/otp"
)

func Main() {
	/*
		There will be 2 ways to punch a slot:
		- Physical transactions: business generates an OTP & physically displays it to the client, client enters it
		- Online transactions: business generates a redeem-link, figure out how to make it secure
	*/

	r := gin.Default()

	// add gzip middleware
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	api := r.Group("/api")
	{
		api.POST("/cards/new", func(c *gin.Context) {
			var card model.CardDetails
			err := c.BindJSON(&card)
			if err != nil {
				c.JSON(400, gin.H{
					"message": "bad request",
				})
				return
			}
			err = db.NewCard(card)
			if err != nil {
				c.JSON(500, gin.H{
					"message": "failed to create card (maybe it already exists?)",
				})
				return
			}
			c.JSON(200, gin.H{
				"message": "card created",
			})
		})

		type CustomizeBody struct {
			CardID string `json:"cardID"`
			Design string `json:"design"`
		}
		api.POST("/cards/customize", func(c *gin.Context) {
			var body CustomizeBody
			err := c.BindJSON(&body)
			if err != nil {
				fmt.Println(err)
				c.JSON(400, gin.H{
					"message": "bad request",
				})
				return
			}
			cd, err := design.ParseCardDesign(body.Design)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"message": "failed to parse JSON",
				})
				return
			}
			err = db.SetCardDesign(body.CardID, cd)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"message": "failed to create card (maybe it already exists?)",
				})
				return
			}
			c.JSON(200, gin.H{
				"message": "card created",
			})
		})

		api.GET("/cards", func(c *gin.Context) {
			cards, err := db.GetAllCardInstances()
			if err != nil {
				c.JSON(500, gin.H{
					"message": "error fetching cards",
				})
				return
			}
			c.JSON(200, gin.H{
				"cards": cards,
			})
		})

		type EnrollBody struct {
			CardID string `json:"cardID"`
		}
		api.POST("/enroll", func(c *gin.Context) {
			var card EnrollBody
			err := c.BindJSON(&card)
			if err != nil {
				c.JSON(400, gin.H{
					"message": "bad request",
				})
				return
			}
			instance := model.CardInstance{
				CardID: card.CardID,
			}
			err = db.InsertCardInstance(instance)
			if err != nil {
				c.JSON(500, gin.H{
					"message": "error inserting card",
				})
				return
			}
			c.JSON(200, gin.H{
				"message": "enrolled",
			})
		})
	}

	r.LoadHTMLGlob("templates/*")
	r.GET("/punch/:card-id", func(c *gin.Context) {
		cardID := c.Param("card-id")
		cardDetails, err := db.GetCardDetails(cardID)
		if err != nil {
			c.JSON(404, gin.H{
				"message": "card not found",
			})
			fmt.Println(err)
			return
		}

		otp, err := otp.GenerateOTP(cardID, cardDetails.Secret)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "failed to generate OTP",
			})
			fmt.Println(err)
			return
		}
		c.HTML(200, "index.html", gin.H{
			"otp":        otp,
			"cardID":     cardID,
			"cardName":   cardDetails.Name,
			"eightTimes": [8]struct{}{},
		})
	})

	r.GET("/validate/:card-id/:otp", func(c *gin.Context) {
		time.Sleep(2 * time.Second)

		enteredOtp := c.Param("otp")
		cardID := c.Param("card-id")
		cardDetails, err := db.GetCardDetails(cardID)
		if err != nil {
			c.JSON(404, gin.H{
				"message": "card not found",
			})
			fmt.Println(err)
			return
		}

		valid := otp.ValidateOTP(cardID, cardDetails.Secret, enteredOtp)
		if valid {
			c.HTML(200, "otpInputSuccess.html", gin.H{})
		} else {
			c.HTML(200, "otpInputFail.html", gin.H{
				"value":  enteredOtp,
				"cardID": cardID,
			})
		}
	})

	r.GET("/cards/all", func(c *gin.Context) {
		cards, err := db.GetAllCardDetails()
		if err != nil {
			c.JSON(500, gin.H{
				"message": "error fetching cards",
			})
			return
		}
		println(cards)
		c.HTML(200, "cards.html", gin.H{
			"cards":      cards,
			"eightTimes": [8]struct{}{},
		})
	})

	r.GET("/enroll", func(c *gin.Context) {
		c.HTML(200, "enroll.html", gin.H{})
	})

	r.GET("/customize/:card-id", func(c *gin.Context) {
		cardID := c.Param("card-id")
		cardDetails, err := db.GetCardDetails(cardID)
		if err != nil {
			c.JSON(404, gin.H{
				"message": "card not found",
			})
			fmt.Println(err)
			return
		}

		c.HTML(200, "customizeCard.html", gin.H{
			"card": cardDetails,
		})
	})

	r.GET("/new", func(c *gin.Context) {
		c.HTML(200, "newCard.html", gin.H{})
	})

	r.Static("/static", "./static")
	r.Run() // listen and serve on
}
