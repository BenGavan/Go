package PASETO_example_1
//
//import (
//	"github.com/stretchr/testify/require"
//	"math/rand"
//	"testing"
//	"time"
//)
//
//func TestJWTMaker(t *testing.T) {
//	maker, err := NewJWTMaker(RandString(32))
//
//	require.NoError(t, err)
//
//	username := "Username here"
//	duration := time.Minute
//
//	issuedAt := time.Now()
//	expiredAt := issuedAt.Add(duration)
//
//	token, err := maker.CreateToken(username, duration)
//	require.NoError(t, err)
//	require.NotEmpty(t, token)
//
//	payload, err := maker.VerifyToken(token)
//	require.NoError(t, err)
//	require.NotEmpty(t, payload)
//
//	require.NotZero(t, payload.ID)
//	require.Equal(t, username, payload.Username)
//	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
//	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
//}
//
//
//
//
//
//
//
//
//
///////////   Random string Generator   /////////
//const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
//const (
//	letterIdxBits = 6                    // 6 bits to represent a letter index
//	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
//	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
//)
//
//func RandString(length int) string {
//	b := make([]byte, length)
//	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
//	for i, cache, remain := length-1, rand.Int63(), letterIdxMax; i >= 0; {
//		if remain == 0 {
//			cache, remain = rand.Int63(), letterIdxMax
//		}
//		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
//			b[i] = letterBytes[idx]
//			i--
//		}
//		cache >>= letterIdxBits
//		remain--
//	}
//
//	return string(b)
//}