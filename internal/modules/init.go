package modules

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"payment-sandbox/internal/middleware"
	"payment-sandbox/internal/modules/admin"
	"payment-sandbox/internal/modules/auth"
	"payment-sandbox/internal/modules/invoice"
	"payment-sandbox/internal/modules/payment"
	"payment-sandbox/internal/modules/refund"
	"payment-sandbox/internal/modules/wallet"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Init(r *gin.RouterGroup, db *gorm.DB) {

	authService := auth.NewService(db)
	walletService := wallet.NewService(db)
	invoiceService := invoice.NewService(db)
	paymentService := payment.NewService(db)
	refundService := refund.NewService(db)
	adminService := admin.NewService(db)

	authHandler := auth.NewHandler(authService)
	walletHandler := wallet.NewHandler(walletService)
	paymentHandler := payment.NewHandler(paymentService)
	invoiceHandler := invoice.NewHandler(invoiceService)
	refundHandler := refund.NewHandler(refundService)
	adminHandler := admin.NewHandler(adminService)

	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)
	r.POST("/auth/refresh", authHandler.Refresh)

	r.GET("/pay/:token", paymentHandler.GetInvoice)
	r.POST("/pay/:token", paymentHandler.Pay)

	// MERCHANT (JWT, Idempotency, Role)
	merchant := r.Group("/merchant")
	merchant.Use(middleware.JWT())
	merchant.Use(middleware.Idempotency(db))
	merchant.Use(middleware.Role("MERCHANT"))

	merchant.GET("/wallet", walletHandler.GetBalance)
	merchant.POST("/topup", walletHandler.TopUp)

	merchant.POST("/invoice", invoiceHandler.Create)
	merchant.GET("/invoice", invoiceHandler.List)
	merchant.PUT("/invoice/:id", invoiceHandler.Update)
	merchant.DELETE("/invoice/:id", invoiceHandler.Delete)

	merchant.POST("/refund", refundHandler.Request)

	merchant.GET("/profile", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "merchant only"})
	})

	// ADMIN
	admin := r.Group("/admin")
	admin.Use(middleware.JWT())
	admin.Use(middleware.Role("ADMIN"))

	admin.PUT("/topup/:id", walletHandler.UpdateTopUpStatus)
	admin.PUT("/payment/:id", paymentHandler.Update)

	admin.PUT("/refund/:id/status", refundHandler.UpdateStatus)
	admin.PUT("/refund/:id/process", refundHandler.Process)

	admin.GET("/dashboard", adminHandler.Dashboard)

	// SWAGGER
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
