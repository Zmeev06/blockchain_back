package main

import (
	"chopcoin/node/routing"
	"chopcoin/shared"
	"chopcoin/shared/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var bc = models.Blockchain{Dir: "blockchain"}

func main() {
	var recipient shared.PublicKey
	if err := json.Unmarshal([]byte(os.Getenv("MINE_RECIPIENT")), &recipient); err != nil {
		fmt.Println(err)
	}
	if os.Getenv("MAIN_NODE") != "" {
		if err := bc.Create(recipient); err != nil {
			fmt.Println("err creating bc")
			fmt.Println(err)
		}
	} else {
		if err := bc.Connect(recipient); err != nil {
			fmt.Println(err)
		}
		go Mine(recipient)
	}
	app := fiber.New()
	app.Use(logger.New())
	routing.Setup(app, bc)
	log.Fatal(app.Listen(os.Getenv("LISTEN_ADDR")))
}
func Mine(recipient shared.PublicKey) {
	addr := os.Getenv("MINE_FOR")
	if addr == "" {
		return
	}
	for {
		mine := fiber.Get(fmt.Sprintf("%s/api/node/mine", addr))
		var tx models.SignedTransaction
		status, bytes, errs := mine.Bytes()
		if len(errs) != 0 {
			fmt.Println("error getting mine")
			fmt.Println(status, errors.Join(errs...).Error())
			return
		}
		json.Unmarshal(bytes, &tx)
		block, err := bc.MineBlock(&tx, recipient)
		if err != nil {
			log.Fatal(err)
		}
		sync := fiber.Post(fmt.Sprintf("%s/api/node/sync", addr))
		sync.JSON(*block)
		{
			status, _, errs := sync.Bytes()
			if len(errs) != 0 {
				fmt.Println("error getting mine")
				fmt.Println(status, errors.Join(errs...).Error())
				return
			}
		}
	}
}
