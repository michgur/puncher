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
		return
	}

	// generate the OTP using the card's secret and the unique transaction ID
	secret := cardSecret + fmt.Sprint(key.TransactionID)
	otp, err = totp.GenerateCodeCustom(secret, time.Now(), validateOpts)
	if err != nil {
		return
	}
	return
}

func ValidateOTP(cardID string, cardSecret string, otp string) bool {
	// get all transaction keys for the card
	keys, err := GetTransactions(cardID)
	if err != nil {
		return false
	}

	// validate the OTP using the card's secret and the unique transaction ID
	for _, key := range keys {
		secret := cardSecret + fmt.Sprint(key.TransactionID)
		valid, err := totp.ValidateCustom(otp, secret, time.Now(), validateOpts)
		if err == nil && valid {
			// asynchonously remove the transaction key
			go func() {
				err = RemoveTransaction(key)
				// handle error
			}()
			return true
		}
	}
	return false
}
