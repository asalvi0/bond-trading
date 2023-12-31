package main

import (
	"fmt"
	"log"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/asalvi0/bond-trading/internal/authn"
	. "github.com/asalvi0/bond-trading/internal/model"
	"github.com/asalvi0/bond-trading/internal/order"
	"github.com/asalvi0/bond-trading/internal/user"
	"github.com/asalvi0/bond-trading/internal/utils/storage/sqlite3"
)

func main() {
	app := fiber.New(fiber.Config{
		// Prefork:     true,
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	setupMiddleware(app)

	authn.RegisterRoutes(app)
	user.RegisterRoutes(app)
	order.RegisterRoutes(app)

	log.Fatal(app.Listen(":8080"))
}

func setupMiddleware(app *fiber.App) {
	storage := sqlite3.New()

	app.Use(favicon.New())
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		TimeFormat: "02-Jan-2006 15:04:05",
	}))

	app.Use(idempotency.New())
	app.Use(requestid.New())
	app.Use(etag.New())

	app.Use(pprof.New())
	app.Get("/metrics", monitor.New())

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(limiter.New(limiter.Config{
		Storage:    storage,
		Max:        1000,
		Expiration: 1 * time.Minute,
	}))

	app.Use(helmet.New())
	app.Use(csrf.New())
}

func orderDemo() {
	order := NewOrder(1, 12, 123, 100, 99.50, Buy, Open)
	if err := order.Validate(); err != nil {
		fmt.Println("Validation error:", err)
		return
	}

	// Marshal the Order struct into JSON with indentation
	prettyJSON, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the prettified JSON
	fmt.Println(string(prettyJSON))
}
