package otp

import (
	"fmt"
	"time"

	"encoding/base32"

	"github.com/pquerna/otp/totp"
)

var validateOpts = totp.ValidateOpts{
	Period: 60,
	Digits: 4,
}

var b32NoPadding = base32.StdEncoding.WithPadding(base32.NoPadding)

func intToB32(i int) string {
	return b32NoPadding.EncodeToString([]byte(fmt.Sprint(i)))
}

func GenerateOTP(cardID string, cardSecret string) (otp string, err error) {
	// generate a unique transaction key
	key, err := AddTransaction(cardID)
	if err != nil {
		return
	}

	// generate the OTP using the card's secret and the unique transaction ID
	secret := cardSecret + intToB32(key.TransactionID)
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
		secret := cardSecret + intToB32(key.TransactionID)
		valid, err := totp.ValidateCustom(otp, secret, time.Now(), validateOpts)
		if err == nil && valid {
			// asynchonously remove the transaction key
			go func() {
				err = RemoveTransaction(key)
				if err != nil {
					println(err)
				}
			}()
			return true
		}
	}
	return false
}
