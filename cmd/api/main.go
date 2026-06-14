package main

import (
	"banksystem/internal/config"
	"banksystem/internal/handlers"
	"banksystem/internal/repositories"
	"banksystem/internal/services"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Подключение к базе данных
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Инициализация логгера
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Инициализация репозиториев
	userRepo := repositories.NewUserRepository(db)
	accountRepo := repositories.NewAccountRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	creditRepo := repositories.NewCreditRepository(db)
	creditPaymentRepo := repositories.NewCreditPaymentRepository(db)
	cardRepo := repositories.NewCardRepository(db)

	// Инициализация SMTP сервиса
	smtpService := services.NewSMTPService(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUsername, cfg.SMTPPassword)

	// Инициализация JWT сервиса
	jwtService := services.NewJWTService(cfg.JWTSecret)

	// Инициализация сервисов
	authService := services.NewAuthService(db, userRepo, jwtService)
	accountService := services.NewAccountService(db, accountRepo, transactionRepo, userRepo, smtpService)
	cardService := services.NewCardService(cardRepo, nil)
	creditService := services.NewCreditService(db, creditRepo, accountRepo, creditPaymentRepo, transactionRepo)
	creditPaymentService := services.NewCreditPaymentService(db, creditPaymentRepo, creditRepo, accountRepo)

	// Инициализация обработчиков
	authHandler := handlers.NewAuthHandler(authService)
	accountHandler := handlers.NewAccountHandler(accountService)
	cardHandler := handlers.NewCardHandler(cardService)
	creditHandler := handlers.NewCreditHandler(creditService)
	creditPaymentHandler := handlers.NewCreditPaymentHandler(creditPaymentService)

	// Инициализация middleware
	authMiddleware := handlers.NewAuthMiddleware(jwtService, logger)

	// Настройка маршрутизации
	mux := http.NewServeMux()

	// Публичные маршруты
	mux.HandleFunc("/api/register", authHandler.Register)
	mux.HandleFunc("/api/login", authHandler.Login)

	// Защищенные маршруты
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/api/accounts/create", accountHandler.CreateAccount)
	protectedMux.HandleFunc("/api/accounts/list", accountHandler.GetUserAccounts)
	protectedMux.HandleFunc("/api/accounts/balance", accountHandler.GetBalance)
	protectedMux.HandleFunc("/api/accounts/deposit", accountHandler.Deposit)
	protectedMux.HandleFunc("/api/accounts/withdraw", accountHandler.Withdraw)
	protectedMux.HandleFunc("/api/accounts/transfer", accountHandler.Transfer)

	protectedMux.HandleFunc("/api/cards/create", cardHandler.CreateCard)
	protectedMux.HandleFunc("/api/cards/list", cardHandler.GetUserCards)
	protectedMux.HandleFunc("/api/cards/get", cardHandler.GetCard)

	protectedMux.HandleFunc("/api/credits/create", creditHandler.CreateCredit)
	protectedMux.HandleFunc("/api/credits/list", creditHandler.GetUserCredits)
	protectedMux.HandleFunc("/api/credits/get", creditHandler.GetCredit)
	protectedMux.HandleFunc("/api/credits/schedule", creditHandler.GetPaymentSchedule)

	protectedMux.HandleFunc("/api/payments/create", creditPaymentHandler.CreatePayment)
	protectedMux.HandleFunc("/api/payments/process", creditPaymentHandler.ProcessPayment)
	protectedMux.HandleFunc("/api/payments/list", creditPaymentHandler.GetPaymentsByCreditID)
	protectedMux.HandleFunc("/api/payments/pending", creditPaymentHandler.GetPendingPayments)

	// Применяем middleware к защищенным маршрутам
	mux.Handle("/api/", authMiddleware.Middleware(protectedMux))

	// Запуск сервера
	logger.Printf("Starting server on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
