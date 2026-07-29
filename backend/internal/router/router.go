package router

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	redisStorage "github.com/gofiber/storage/redis/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kickpick/backend/internal/config"
	"github.com/kickpick/backend/internal/handler"
	appmw "github.com/kickpick/backend/internal/middleware"
)

// newLimiterConfig builds rate-limit config with Redis-backed storage when
// REDIS_URL is set, so counters stay consistent across multiple server
// instances (Section 19 PRD). Falls back to Fiber's default in-memory storage
// otherwise — fine for single-instance/local dev, but each instance would
// count independently in a multi-instance deployment without Redis.
//
// name namespaces the storage key (Fiber's default KeyGenerator is bare
// c.IP(), so without a prefix every limiter instance sharing one Redis
// backend collides on the same key — general page traffic counted by the
// global 100/min limiter was silently exhausting the auth-specific 10/min
// budget before this fix, since both wrote to the same "<ip>" key).
func newLimiterConfig(cfg *config.Config, name string, max int, expiration time.Duration) limiter.Config {
	lc := limiter.Config{
		Max:        max,
		Expiration: expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			return name + ":" + c.IP()
		},
	}

	if cfg.RedisURL != "" {
		store := redisStorage.New(redisStorage.Config{URL: cfg.RedisURL})
		lc.Storage = store
	}

	return lc
}

func New(cfg *config.Config, pool *pgxpool.Pool) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "KickPick API",
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(helmet.New(helmet.Config{
		ContentSecurityPolicy: "default-src 'none'; frame-ancestors 'none'",
		XFrameOptions:         "DENY",
		ContentTypeNosniff:    "nosniff",
		HSTSMaxAge:            31536000,
		HSTSExcludeSubdomains: false,
		ReferrerPolicy:        "strict-origin-when-cross-origin",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://127.0.0.1:3000, " + cfg.CORSAllowedOrigin,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		AllowMethods:     "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
	}))
	if cfg.RedisURL != "" {
		log.Println("rate limiting: using Redis-backed storage")
	} else {
		log.Println("rate limiting: REDIS_URL not set, using in-memory storage (fine for single-instance/local dev only)")
	}

	app.Use(limiter.New(newLimiterConfig(cfg, "global", 100, 1*time.Minute)))

	app.Get("/health", handler.Health)

	if pool == nil {
		api := app.Group("/api")
		api.Get("/*", handler.NotImplemented)
		return app
	}

	authHandler := handler.NewAuthHandler(pool, cfg)
	productHandler := handler.NewProductHandler(pool)
	wishlistHandler := handler.NewWishlistHandler(pool)
	notificationHandler := handler.NewNotificationHandler(pool)
	reviewHandler := handler.NewReviewHandler(pool)
	profileHandler := handler.NewProfileHandler(pool)
	redirectHandler := handler.NewRedirectHandler(pool)
	exchangeRateHandler := handler.NewExchangeRateHandler(pool)

	requireAuth := appmw.RequireAuth(cfg.JWTAccessSecret)
	requireCSRF := appmw.RequireCSRFHeader()

	// Separate limiter instances per concern — sharing one instance across auth
	// AND redirect would mean a burst of affiliate clicks could lock a user out
	// of login, and vice versa (found via E2E tests all logging in against a
	// shared IP tripping each other's budget).
	authLimiter := limiter.New(newLimiterConfig(cfg, "auth", 10, 1*time.Minute))
	redirectLimiter := limiter.New(newLimiterConfig(cfg, "redirect", 30, 1*time.Minute))

	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", authLimiter, authHandler.Register)
	auth.Post("/login", authLimiter, authHandler.Login)
	auth.Post("/refresh", requireCSRF, authHandler.Refresh)
	auth.Post("/logout", requireCSRF, authHandler.Logout)
	auth.Post("/verify-email", authHandler.VerifyEmail)
	auth.Post("/forgot-password", authLimiter, authHandler.ForgotPassword)
	auth.Post("/reset-password", authHandler.ResetPassword)

	api.Get("/products", productHandler.ListProducts)
	api.Get("/products/trending", productHandler.ListTrending)
	api.Get("/products/price-drops", productHandler.ListPriceDrops)
	api.Get("/products/:slug", productHandler.GetProductBySlug)
	api.Get("/products/:slug/price-history", productHandler.GetPriceHistory)
	api.Get("/products/:slug/size-conversion", productHandler.GetSizeConversion)
	api.Get("/brands", productHandler.ListBrands)
	api.Get("/brands/:slug", productHandler.GetBrandBySlug)
	api.Get("/search/autocomplete", productHandler.Autocomplete)
	api.Get("/exchange-rate", exchangeRateHandler.GetIDRToUSD)

	api.Post("/redirect/:offer_id", redirectLimiter, requireCSRF, redirectHandler.GoToOffer)

	api.Post("/reviews", requireAuth, requireCSRF, reviewHandler.Create)
	api.Post("/reviews/:id/report", requireAuth, requireCSRF, reviewHandler.Report)

	api.Get("/wishlist", requireAuth, wishlistHandler.List)
	api.Post("/wishlist", requireAuth, requireCSRF, wishlistHandler.Add)
	api.Delete("/wishlist/:id", requireAuth, requireCSRF, wishlistHandler.Remove)
	api.Patch("/wishlist/:id/alert", requireAuth, requireCSRF, wishlistHandler.SetAlert)

	api.Get("/notifications", requireAuth, notificationHandler.List)
	api.Get("/notifications/unread-count", requireAuth, notificationHandler.UnreadCount)
	api.Patch("/notifications/:id/read", requireAuth, requireCSRF, notificationHandler.MarkRead)

	api.Get("/profile", requireAuth, profileHandler.Get)
	api.Patch("/profile", requireAuth, requireCSRF, profileHandler.Update)
	api.Delete("/profile", requireAuth, requireCSRF, profileHandler.Delete)

	return app
}
