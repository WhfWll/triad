//go:build ignore

package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/urfave/cli"
	"smart/services"
	"smart/tools/enums"
)

func main() {
	app := cli.NewApp()
	app.Usage = "auth test tools"
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "server"},
		&cli.StringFlag{Name: "scenario"},
		&cli.IntFlag{Name: "days"},
		&cli.StringFlag{Name: "overrideProductID"},
		&cli.BoolFlag{Name: "post"},
	}
	app.Action = func(c *cli.Context) error {
		server := c.String("server")
		scenario := c.String("scenario")
		days := c.Int("days")
		overrideProductID := c.String("overrideProductID")
		post := c.Bool("post")
		ctx := context.Background()
		var auth services.Auth
		serialNumber, err := auth.GenerateSystemSerialNumber(ctx)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		productID := serialNumber
		if strings.TrimSpace(overrideProductID) != "" {
			productID = overrideProductID
		}
		var generateTime string
		switch scenario {
		case "expired_token":
			generateTime = time.Now().AddDate(0, 0, -100).Format(enums.TimeYMDBarLayout)
		case "malformed_time":
			generateTime = "2025/01/01"
		default:
			generateTime = time.Now().Format(enums.TimeYMDBarLayout)
		}
		var authDays string
		switch scenario {
		case "invalid_days":
			authDays = "abc"
		case "negative_days":
			authDays = "-1"
		case "large_days":
			authDays = "10000"
		default:
			if days == 0 {
				authDays = "7"
			} else {
				authDays = fmt.Sprintf("%d", days)
			}
		}
		var payload map[string]string
		switch scenario {
		case "wrong_product":
			payload = map[string]string{
				"productID":    productID + "_WRONG",
				"generateTime": generateTime,
				"authDays":     authDays,
			}
		default:
			payload = map[string]string{
				"productID":    productID,
				"generateTime": generateTime,
				"authDays":     authDays,
			}
		}
		data, _ := json.Marshal(payload)
		cipher := auth.RsaEncrypt(data, []byte(enums.SWPubKey))
		var code string
		switch scenario {
		case "bad_hex":
			code = "ZZZZ"
		case "empty_code":
			code = ""
		default:
			code = hex.EncodeToString(cipher)
		}
		fmt.Println("scenario:", scenario)
		fmt.Println("productID:", productID)
		fmt.Println("generateTime:", generateTime)
		fmt.Println("authDays:", authDays)
		fmt.Println("authCode:", code)
		if post {
			if strings.TrimSpace(server) == "" {
				fmt.Println("server required")
				return nil
			}
			body, _ := json.Marshal(map[string]string{"authCode": code})
			req, _ := http.NewRequest("POST", strings.TrimRight(server, "/")+"/smart/system/authsave", strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
			defer resp.Body.Close()
			fmt.Println("post status:", resp.StatusCode)
		}
		return nil
	}
	_ = app.Run([]string{})
}
