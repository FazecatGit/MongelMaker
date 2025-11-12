package utils

import (
	"fmt"
	"time"
)

type RetryConfig struct {
	MaxRetries int           // How many times to retry (e.g., 3)
	Delay      time.Duration // Initial delay between retries (e.g., 2 seconds)
	Backoff    float64       // Multiplier for exponential backoff (e.g., 2.0)
}

func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries: 3,
		Delay:      time.Second * 2,
		Backoff:    2.0,
	}
}

func RetryWithBackoff(operation func() error, config *RetryConfig) error {
	delay := config.Delay
	for i := 0; i < config.MaxRetries; i++ {
		err := operation()
		if err == nil {
			return nil
		}
		if i < config.MaxRetries-1 {
			fmt.Printf("⚠️  Attempt %d failed: %v. Retrying in %s...\n", i+1, err, delay)
			time.Sleep(delay)
			delay = time.Duration(float64(delay) * config.Backoff)
		}
	}
	return fmt.Errorf("operation failed after %d attempts", config.MaxRetries)
}

func TestRetryLogic() {
	fmt.Println("🧪 Testing retry logic...")

	failCount := 0
	config := DefaultRetryConfig()

	err := RetryWithBackoff(func() error {
		failCount++
		if failCount < 3 {
			return fmt.Errorf("simulated failure %d", failCount)
		}
		fmt.Println("✅ Success on attempt", failCount)
		return nil
	}, config)

	if err != nil {
		fmt.Printf("❌ Test failed: %v\n", err)
	} else {
		fmt.Println("✅ Retry test completed successfully!")
	}
}
