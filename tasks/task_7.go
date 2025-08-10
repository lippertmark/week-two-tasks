package tasks

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	//"strconv"
)

type contextKey string

var jwtKey = contextKey("jwtKey")

var secretKey = []byte("secret-key")

var notAuthorizedError = errors.New("User is not authorized")
var notValidToken = errors.New("Not valid token")

func AddJWTToContext(ctx context.Context, userID int) (context.Context, error) {
	userIDStr := strconv.Itoa(userID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{Subject: userIDStr})
	tokenString, err := token.SignedString([]byte("secret-key"))
	if err != nil {
		return ctx, err
	}
	fmt.Println("mytoken:", tokenString)
	ctx = context.WithValue(ctx, jwtKey, tokenString)
	return ctx, nil
}

func ExtractUserIDFromContext(ctx context.Context) (int, error) {
	jwtToken := ctx.Value(jwtKey)
	if jwtToken == nil {
		return 0, notAuthorizedError
	}
	token, err := jwt.Parse(fmt.Sprintf("%v", jwtToken), func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, notValidToken
	}
	userIdStr, err2 := token.Claims.GetSubject()
	if err2 != nil {
		return 0, err2
	}
	userId, err3 := strconv.Atoi(userIdStr)
	if err3 != nil {
		return 0, err3
	}
	return userId, nil
}

func Task7() {
	ctx, err := AddJWTToContext(context.Background(), 1233)
	if err != nil {
		fmt.Println(err)
	}
	userId, err2 := ExtractUserIDFromContext(ctx)
	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Println(userId)
}
