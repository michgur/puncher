package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/a-h/templ/examples/integration-gin/gintemplrenderer"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/michgur/puncher/app/db"
	"github.com/michgur/puncher/app/design"
	"github.com/michgur/puncher/app/model"
	"github.com/michgur/puncher/app/otp"
	"github.com/michgur/puncher/app/templ"

	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
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

var ginLambda *ginadapter.GinLambda
var r *gin.Engine

func init() {
	/*
		There will be 2 ways to punch a slot:
		- Physical transactions: business generates an OTP & physically displays it to the client, client enters it
		- Online transactions: business generates a redeem-link, figure out how to make it secure
	*/

	r = gin.Default()
	r.HTMLRender = gintemplrenderer.Default

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

	// r.LoadHTMLGlob("templates/*")
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
		c.HTML(200, "", gin.H{
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
		c.HTML(200, "", templ.AllCards(cards))
	})

	r.GET("/enroll", func(c *gin.Context) {
		c.HTML(200, "enroll.html", gin.H{})
	})

	r.POST("/customize/:card-id", func(c *gin.Context) {
		cardID := c.Param("card-id")

		var body map[string]interface{}
		err := c.BindJSON(&body)
		if !validateRequest(c, err) {
			return
		}

		var card model.CardDetails
		card.Name = body["name"].(string)
		body["textureOpacity"], err = strconv.Atoi(body["textureOpacity"].(string))
		if !validateRequest(c, err) {
			return
		}
		s, err := json.Marshal(body)
		if !validateRequest(c, err) {
			return
		}
		err = json.Unmarshal(s, &card.Design)
		if !validateRequest(c, err) {
			return
		}

		err = db.SetCardNameAndDesign(cardID, card.Name, card.Design)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"message": "failed to update card",
			})
			return
		}
		c.HTML(200, "", templ.Card(card, true))
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

		conf, err := design.ReadDesignConfig()
		if err != nil {
			c.JSON(500, gin.H{
				"message": "failed to load customization options",
			})
			fmt.Println(err)
			return
		}
		c.HTML(200, "", templ.CustomizeCard(cardDetails, conf))
	})

	r.GET("/new", func(c *gin.Context) {
		c.HTML(200, "newCard.html", gin.H{})
	})

	r.GET("/hello", func(c *gin.Context) {
		c.HTML(200, "", templ.Card(model.CardDetails{
			Name: "Yael's Fan Club",
			Design: design.CardDesign{
				Color:          "citron",
				Font:           "font-pacifico",
				Pattern:        "bubbles.svg",
				Texture:        "noise-dark.png",
				TextureOpacity: 30,
			},
		}, false))
	})

	r.Static("/static", "./static")

	ginLambda = ginadapter.New(r)
}

func Handler(context context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(context, request)
}

func Main() {
	if os.Getenv("ENV") == "lambda" {
		lambda.Start(Handler)
	} else {
		r.Run(":8080")
	}
}
