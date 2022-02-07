package main

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

func main() {
	//test()
	//standardClaims()
	//customClaimsExample()
	//parsingWithBitfieldChecks()
	//parsingHmac()
	testRun()
}

func test() {
	mySigningKey := []byte("AllYourBase")

	type MyCustomClaims struct {
		Foo string `json:"foo"`
		jwt.StandardClaims
	}

	// Create the Claims
	claims := MyCustomClaims{
		"bar",
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	fmt.Printf("%v %v", ss, err)
}

func customClaimsExample() {

	signingKey := []byte("AllYourBase")

	type CustomClaims struct {
		Foo string `json:"foo"`
		jwt.StandardClaims
	}

	claims := CustomClaims{
		Foo: "bar",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(signingKey)
	fmt.Printf("%v %v\n", ss, err)
}

func standardClaims() {
	mySigningKey := []byte("AllYourBase")

	// Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: 15000,
		Issuer:    "test",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	fmt.Printf("%v %v", ss, err)
}

func testRun() {
	signingKey := []byte("AllYourBase")

	fmt.Printf("--- Making Token\n")

	t, err := makeToken(signingKey)
	fmt.Printf("token string: %v,  err: %v\n", t, err)
	if err != nil {
		return
	}

	fmt.Printf("--- Parsing Token\n")

	//t += "g"



	token, err := parseToken(t, signingKey)
	fmt.Printf("token: %v, err: %v\n",  token, err)
	if err != nil {
		return
	}
}

type CustomClaims struct {
	UUID string `json:"uuid"`
	jwt.StandardClaims
}

func makeToken(signingKey []byte) (string, error) {

	claims := CustomClaims{
		UUID: "thisistheuuid",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + time.Minute.Milliseconds() * 0,
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tss, err := token.SignedString(signingKey)
	fmt.Printf("%v %v\n", tss, err)
	return tss, err
}

func parseToken(ts string, stringingKey []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(ts, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return stringingKey, nil
	})

	if token == nil {
		fmt.Printf("error parsing token: %v\n", err)
		return nil, err
	}

	if token.Valid {
		return token, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		var es string
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			es = "That's not even a token"
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			es = "Timing is everything"
		} else {
			es = fmt.Sprintf("Couldn't handle this token:", err)
		}
		return nil, errors.New(es)
	}
	return nil, err
}

// An example of parsing the error types using bitfield checks
func parsingWithBitfieldChecks() {
	// Token from another example.  This token is expired
	var tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c"

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})

	if token == nil {
		fmt.Printf("error parsing token: %v\n", err)
		return
	}

	fmt.Printf("raw token: %v\n", token.Raw)
	fmt.Printf("token claims: %v\n", token.Claims)
	fmt.Printf("token signature: %v\n", token.Signature)
	fmt.Printf("token header: %v\n", token.Header)
	fmt.Printf("is token valid: %v\n", token.Valid)
	fmt.Printf("token method: %v\n", token.Method)

	if token.Valid {
		fmt.Println("You look nice today")
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			fmt.Println("Timing is everything")
		} else {
			fmt.Println("Couldn't handle this token:", err)
		}
	}

	if token.Valid {

	} else if ve, ok := err.(*jwt.ValidationError); ok {
		fmt.Printf("ve.Errors: %v = %v\n", ve.Errors, strconv.FormatInt(int64(ve.Errors), 2))
		fmt.Printf("ValidationErrorMalformed: %v\n", jwt.ValidationErrorMalformed)
		fmt.Printf("%v\n", ve.Errors&jwt.ValidationErrorMalformed)
		fmt.Printf("jwt.ValidationErrorExpired = %v\n", strconv.FormatInt(int64(jwt.ValidationErrorExpired), 2))
		fmt.Printf("jwt.ValidationErrorNotValidYet = %v\n", strconv.FormatInt(int64(jwt.ValidationErrorNotValidYet), 2))
		fmt.Printf("jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet = %v\n", strconv.FormatInt(int64(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet), 2))
		fmt.Printf("ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) = %v\n", strconv.FormatInt(int64(ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet)), 2))

	}
}

func parsingHmac() {
	// sample token string taken from the New example
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

	hmacSampleSecret := []byte("my_secret_key")

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})

	fmt.Printf("token: %v, err: %v\n", token, err)

}
