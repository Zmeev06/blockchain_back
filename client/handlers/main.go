package handlers

import (
	"chopcoin/client/models"
	"encoding/json"

	"os"
	"path"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const USERS = "users"
const DATA = "useless_crap"

var JWT_SECRET []byte
var NODE_ADDR string

type creds struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func Init() {
	os.MkdirAll(USERS, 0700)
	os.MkdirAll(DATA, 0777)
	JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))
	NODE_ADDR = os.Getenv("NODE_ADDR")
}
func makeCrapPath(name string) string {
	return path.Join(DATA, name)
}
func getUselessCrapByName(name string) (crap models.Crap, err error) {
	bytes, err := os.ReadFile(makeCrapPath(name))
	if err != nil {
		return
	}
	if err = json.Unmarshal(bytes, &crap); err != nil {
		return
	}

	return
}
func getUserByName(name string) (user models.User, err error) {
	bytes, err := os.ReadFile(makeUserPath(name))
	if err != nil {
		return
	}
	if err = json.Unmarshal(bytes, &user); err != nil {
		return
	}
	return
}
func getUserFromJwt(ctx *fiber.Ctx) (user models.User, err error) {
	token := ctx.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user, err = getUserByName(claims["identity"].(string))
	if err != nil {
		return
	}
	return
}
func makeToken(user models.User) (str string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["identity"] = user.Login
	claims["exp"] = time.Now().Add(time.Hour * 256).Unix()
	str, err = token.SignedString(JWT_SECRET)
	if err != nil {
		return
	}
	return
}
func makeUserPath(name string) string {
	return path.Join(USERS, name)
}
