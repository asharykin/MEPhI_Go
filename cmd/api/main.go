package main

import (
	"banksystem/internal/config"
	"banksystem/internal/handler"
	"banksystem/internal/logger"
	"banksystem/internal/middleware"
	"banksystem/internal/repository"
	"banksystem/internal/service"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	logger.Init()
	logger.Info("Starting Bank Service API")

	cfg := config.LoadConfig()

	storage, err := repository.NewStorage(cfg.Database.ConnectionString)
	if err != nil {
		logger.Fatal("Failed to connect to database: %v", err)
	}
	defer storage.Close()

	userRepo := repository.NewUserRepository(storage)
	accountRepo := repository.NewAccountRepository(storage)
	cardRepo := repository.NewCardRepository(storage)
	transactionRepo := repository.NewTransactionRepository(storage)
	creditRepo := repository.NewCreditRepository(storage)
	paymentScheduleRepo := repository.NewPaymentScheduleRepository(storage)

	keyRateProvider := service.NewCentralBankKeyRateProvider()
	emailService := service.NewEmailService(
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.User,
		cfg.SMTP.Pass,
	)

	userService := service.NewUserService(userRepo, cfg.JWT.Secret)
	accountService := service.NewAccountService(accountRepo, creditRepo, paymentScheduleRepo)
	cardService := service.NewCardService(cardRepo, accountRepo, cfg.Security.HMACSecret, cfg.Security.PGPPublicKey, cfg.Security.PGPPrivateKey)

	transactionService := service.NewTransactionService(transactionRepo, accountRepo, storage)

	creditService := service.NewCreditService(
		creditRepo,
		paymentScheduleRepo,
		accountRepo,
		transactionRepo,
		userRepo,
		keyRateProvider,
		emailService,
		storage,
	)

	analyticsService := service.NewAnalyticsService(transactionRepo, creditRepo, accountRepo, paymentScheduleRepo)

	userHandler := handler.NewUserHandler(userService)
	accountHandler := handler.NewAccountHandler(accountService)
	cardHandler := handler.NewCardHandler(cardService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	creditHandler := handler.NewCreditHandler(creditService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	r := mux.NewRouter()

	r.Use(middleware.LoggingMiddleware)

	r.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/api/login", userHandler.Login).Methods("POST")

	authRouter := r.PathPrefix("/").Subrouter()
	authRouter.Use(middleware.AuthMiddleware([]byte(cfg.JWT.Secret)))

	authRouter.HandleFunc("/api/accounts", accountHandler.CreateAccount).Methods("POST")
	authRouter.HandleFunc("/api/accounts", accountHandler.GetAccounts).Methods("GET")
	authRouter.HandleFunc("/api/accounts/{id}/deposit", accountHandler.Deposit).Methods("POST")
	authRouter.HandleFunc("/api/accounts/{id}/withdraw", accountHandler.Withdraw).Methods("POST")
	authRouter.HandleFunc("/api/accounts/{id}/predict", accountHandler.PredictBalance).Methods("GET")

	authRouter.HandleFunc("/api/cards", cardHandler.CreateCard).Methods("POST")
	authRouter.HandleFunc("/api/cards", cardHandler.GetCards).Methods("GET")

	authRouter.HandleFunc("/api/transfer", transactionHandler.Transfer).Methods("POST")

	authRouter.HandleFunc("/api/credits", creditHandler.CreateCredit).Methods("POST")
	authRouter.HandleFunc("/api/credits/{id}/schedule", creditHandler.GetSchedule).Methods("GET")

	authRouter.HandleFunc("/api/analytics", analyticsHandler.GetAnalytics).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("Server starting on %s", addr))

	go func() {
		ticker := time.NewTicker(12 * time.Hour)
		defer ticker.Stop()

		logger.Info("Scheduler started for processing scheduled payments every 12 hours")

		for range ticker.C {
			logger.Info("Running scheduled payment processing...")
			ctx := context.Background()
			err := creditService.ProcessScheduledPayments(ctx)
			if err != nil {
				logger.Error("Error during scheduled payment processing", "error", err)
			} else {
				logger.Info("Scheduled payment processing completed successfully")
			}
		}
	}()

	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
