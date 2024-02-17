package otp

import (
	"fmt"
	"time"

	"github.com/pquerna/otp/totp"
)

var validateOpts = totp.ValidateOpts{
	Period: 60,
	Digits: 4,
}

func GenerateOTP(cardID string, cardSecret string) (otp string, err error) {
	// generate a unique transaction key
	key, err := AddTransaction(cardID)
	if err != nil {
		panic(err)
	}

	// generate the OTP using the card's secret and the unique transaction ID
	secret := cardSecret + fmt.Sprint(key.TransactionID)
	otp, err = totp.GenerateCodeCustom(secret, time.Now(), validateOpts)
	if err != nil {
		panic(err)
	}
	return
}

func ValidateOTP(cardID string, cardSecret string, otp string) (valid bool, err error) {
	// get all transaction keys for the card
	keys, err := GetTransactions(cardID)
	if err != nil {
		panic(err)
	}

	// validate the OTP using the card's secret and the unique transaction ID
	for _, key := range keys {
		secret := cardSecret + fmt.Sprint(key.TransactionID)
		valid, err = totp.ValidateCustom(otp, secret, time.Now(), validateOpts)
		if err != nil {
			valid = false
		}
		if valid {
			// asynchonously remove the transaction key
			go func() {
				err = RemoveTransaction(key)
				if err != nil {
					panic(err)
				}
			}()
			return
		}
	}
	return
}
