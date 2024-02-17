package app

import (
	"fmt"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/michgur/puncher/app/db"
	"github.com/michgur/puncher/app/model"
	"github.com/pquerna/otp/totp"
)

// generate a TOTP for a client to redeem their slot punch
// use the business id to generate the TOTP
func generateTOTPSecret(businessID string) (secret string) {
	otp, err := totp.Generate(totp.GenerateOpts{
		Issuer:      businessID,
		AccountName: "client",
	})
	if err != nil {
		panic(err)
	}
	return otp.Secret()
}

var validateOpts = totp.ValidateOpts{
	Period: 30,
	Digits: 4,
}

func generateOTP(secret string) (otp string) {
	otp, err := totp.GenerateCodeCustom(secret, time.Now(), validateOpts)
	if err != nil {
		panic(err)
	}
	return otp
}

func validateOTP(secret, otp string) (valid bool) {
	valid, err := totp.ValidateCustom(otp, secret, time.Now(), validateOpts)
	if err != nil {
		valid = false
	}
	return valid
}

var businessIdToSecret = map[string]string{}

func Main() {
	/*
		There will be 2 ways to punch a slot:
		- Physical transactions: business generates an OTP & physically displays it to the client, client enters it
		- Online transactions: business generates a redeem-link, figure out how to make it secure
	*/

	for _, bid := range []string{"1", "2", "3"} {
		businessIdToSecret["business"+bid] = generateTOTPSecret(bid)
	}

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
		api.GET("/generate/:business-id", func(c *gin.Context) {
			businessID := c.Param("business-id")
			if secret, ok := businessIdToSecret[businessID]; ok {
				otp := generateOTP(secret)
				c.JSON(200, gin.H{
					"message": "generate3",
					"otp":     otp,
				})
			} else {
				c.JSON(404, gin.H{
					"message": "business not found",
				})
			}
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
			CardID string `json:"cardId"`
		}
		api.POST("/enroll", func(c *gin.Context) {
			var card EnrollBody
			err := c.BindJSON(&card)
			fmt.Println("oh shit", c.Request.Body)
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
	r.GET("/punch/:business-id", func(c *gin.Context) {
		businessID := c.Param("business-id")
		if secret, ok := businessIdToSecret[businessID]; ok {
			otp := generateOTP(secret)
			c.HTML(200, "index.html", gin.H{
				"otp":        otp,
				"businessId": businessID,
				"eightTimes": [8]struct{}{},
			})
		} else {
			c.JSON(404, gin.H{
				"message": "business not found",
			})
		}
	})

	r.GET("/validate/:business-id/:otp", func(c *gin.Context) {
		time.Sleep(2 * time.Second)

		businessID := c.Param("business-id")
		otp := c.Param("otp")
		if secret, ok := businessIdToSecret[businessID]; ok {
			valid := validateOTP(secret, otp)
			if valid {
				c.HTML(200, "otpInputSuccess.html", gin.H{})
			} else {
				c.HTML(200, "otpInputFail.html", gin.H{
					"value":      otp,
					"businessId": businessID,
				})
			}
		} else {
			c.JSON(404, gin.H{
				"message": "business not found",
			})
		}
	})

	r.GET("/enroll", func(c *gin.Context) {
		c.HTML(200, "enroll.html", gin.H{})
	})

	r.Static("/static", "./static")
	r.Run() // listen and serve on
}