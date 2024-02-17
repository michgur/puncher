package otp

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"encoding/base32"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

var validateOpts totp.ValidateOpts

func init() {
	// get otp lifetime from environment variable OTP_LIFETIME_SECONDS
	otpLifetime, err := strconv.Atoi(os.Getenv("OTP_LIFETIME_SECONDS"))
	if err != nil {
		fmt.Println("error getting OTP_LIFETIME_SECONDS from environment", err)
		otpLifetime = 60
	}
	digits, err := strconv.Atoi(os.Getenv("OTP_DIGITS"))
	if err != nil {
		fmt.Println("error getting OTP_DIGITS from environment", err)
		digits = 4
	}
	// setup the validation options
	validateOpts = totp.ValidateOpts{
		Period: uint(otpLifetime),
		Digits: otp.Digits(digits),
	}
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

func tryRemoveTransaction(key TransactionKey) {
	err := RemoveTransaction(key)
	if err != nil {
		println(err)
	}
}

func ValidateOTP(cardID string, cardSecret string, otp string) bool {
	// get all transaction keys for the card
	keys, err := GetTransactions(cardID)
	if err != nil {
		return false
	}

	// validate the OTP using the card's secret and the unique transaction ID
	for _, key := range keys {
		if key.UnixTime+int64(validateOpts.Period) < time.Now().Unix() {
			go tryRemoveTransaction(key)
		}
		secret := cardSecret + intToB32(key.TransactionID)
		valid, err := totp.ValidateCustom(otp, secret, time.Now(), validateOpts)
		if err == nil && valid {
			// asynchonously remove the transaction key
			go tryRemoveTransaction(key)
			return true
		}
	}
	return false
}
