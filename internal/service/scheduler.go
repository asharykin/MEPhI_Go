package service

import (
	"context"
	"log"
	"time"
)

type Scheduler struct {
	creditPaymentService *CreditPaymentService
	stopChan             chan struct{}
}

func NewScheduler(creditPaymentService *CreditPaymentService) *Scheduler {
	return &Scheduler{
		creditPaymentService: creditPaymentService,
		stopChan:             make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				s.processPayments()
			case <-s.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.stopChan)
}

func (s *Scheduler) processPayments() {
	ctx := context.Background()
	payments, err := s.creditPaymentService.GetPendingPayments(ctx)
	if err != nil {
		log.Printf("Error getting pending payments: %v", err)
		return
	}

	for _, payment := range payments {
		if payment.DueDate.Before(time.Now()) {
			err := s.creditPaymentService.ProcessPayment(ctx, payment.ID)
			if err != nil {
				log.Printf("Error processing payment %d: %v", payment.ID, err)
				continue
			}
			log.Printf("Successfully processed payment %d", payment.ID)
		}
	}
}
