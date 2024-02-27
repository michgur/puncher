package app

import (
	"fmt"

	"github.com/a-h/templ/examples/integration-gin/gintemplrenderer"
	"github.com/gin-gonic/gin"
	"github.com/michgur/puncher/app/db"
	"github.com/michgur/puncher/app/model"
	"github.com/michgur/puncher/app/templ"
)

func validateRequest(c *gin.Context, err error) bool {
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"message": "bad request",
		})
		return false
	}
	return true
}

var R *gin.Engine

func init() {
	/*
		There will be 2 ways to punch a slot:
		- Physical transactions: business generates an OTP & physically displays it to the client, client enters it
		- Online transactions: business generates a redeem-link, figure out how to make it secure
	*/
	R = gin.Default()
	R.HTMLRender = gintemplrenderer.Default

	// add gzip middleware
	// R.Use(gzip.Gzip(gzip.DefaultCompression))

	R.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong3",
		})
	})

	api := R.Group("/api")
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

	// // r.LoadHTMLGlob("templates/*")
	// R.GET("/punch/:card-id", func(c *gin.Context) {
	// 	cardID := c.Param("card-id")
	// 	cardDetails, err := db.GetCardDetails(cardID)
	// 	if err != nil {
	// 		c.JSON(404, gin.H{
	// 			"message": "card not found",
	// 		})
	// 		fmt.Println(err)
	// 		return
	// 	}

	// 	otp, err := otp.GenerateOTP(cardID, cardDetails.Secret)
	// 	if err != nil {
	// 		c.JSON(500, gin.H{
	// 			"message": "failed to generate OTP",
	// 		})
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	c.HTML(200, "", gin.H{
	// 		"otp":        otp,
	// 		"cardID":     cardID,
	// 		"cardName":   cardDetails.Name,
	// 		"eightTimes": [8]struct{}{},
	// 	})
	// })

	// R.GET("/validate/:card-id/:otp", func(c *gin.Context) {
	// 	time.Sleep(2 * time.Second)

	// 	enteredOtp := c.Param("otp")
	// 	cardID := c.Param("card-id")
	// 	cardDetails, err := db.GetCardDetails(cardID)
	// 	if err != nil {
	// 		c.JSON(404, gin.H{
	// 			"message": "card not found",
	// 		})
	// 		fmt.Println(err)
	// 		return
	// 	}

	// 	valid := otp.ValidateOTP(cardID, cardDetails.Secret, enteredOtp)
	// 	if valid {
	// 		c.HTML(200, "otpInputSuccess.html", gin.H{})
	// 	} else {
	// 		c.HTML(200, "otpInputFail.html", gin.H{
	// 			"value":  enteredOtp,
	// 			"cardID": cardID,
	// 		})
	// 	}
	// })

	R.GET("/cards/all", func(c *gin.Context) {
		cards, err := db.GetAllCardDetails()
		if err != nil {
			c.JSON(500, gin.H{
				"message": "error fetching cards",
			})
			return
		}
		c.HTML(200, "", templ.AllCards(cards))
	})

	// R.GET("/enroll", func(c *gin.Context) {
	// 	c.HTML(200, "enroll.html", gin.H{})
	// })

	// R.POST("/customize/:card-id", func(c *gin.Context) {
	// 	cardID := c.Param("card-id")

	// 	var body map[string]interface{}
	// 	err := c.BindJSON(&body)
	// 	if !validateRequest(c, err) {
	// 		return
	// 	}

	// 	var card model.CardDetails
	// 	card.Name = body["name"].(string)
	// 	body["textureOpacity"], err = strconv.Atoi(body["textureOpacity"].(string))
	// 	if !validateRequest(c, err) {
	// 		return
	// 	}
	// 	s, err := json.Marshal(body)
	// 	if !validateRequest(c, err) {
	// 		return
	// 	}
	// 	err = json.Unmarshal(s, &card.Design)
	// 	if !validateRequest(c, err) {
	// 		return
	// 	}

	// 	err = db.SetCardNameAndDesign(cardID, card.Name, card.Design)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		c.JSON(500, gin.H{
	// 			"message": "failed to update card",
	// 		})
	// 		return
	// 	}
	// 	c.HTML(200, "", templ.Card(card, true))
	// })

	// R.GET("/customize/:card-id", func(c *gin.Context) {
	// 	cardID := c.Param("card-id")
	// 	cardDetails, err := db.GetCardDetails(cardID)
	// 	if err != nil {
	// 		c.JSON(404, gin.H{
	// 			"message": "card not found",
	// 		})
	// 		fmt.Println(err)
	// 		return
	// 	}

	// 	conf, err := design.ReadDesignConfig()
	// 	if err != nil {
	// 		c.JSON(500, gin.H{
	// 			"message": "failed to load customization options",
	// 		})
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	c.HTML(200, "", templ.CustomizeCard(cardDetails, conf))
	// })

	// R.GET("/new", func(c *gin.Context) {
	// 	c.HTML(200, "newCard.html", gin.H{})
	// })

	// R.GET("/hello", func(c *gin.Context) {
	// 	c.HTML(200, "", templ.Card(model.CardDetails{
	// 		Name: "Yael's Fan Club",
	// 		Design: design.CardDesign{
	// 			Color:          "citron",
	// 			Font:           "font-pacifico",
	// 			Pattern:        "bubbles.svg",
	// 			Texture:        "noise-dark.png",
	// 			TextureOpacity: 30,
	// 		},
	// 	}, false))
	// })

	R.Static("/static", "./static")
}
