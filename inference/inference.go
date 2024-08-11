package inference

import (
	"errors"
	"sync"
	"time"
)

// Structure to hold provider data and performance metrics
type Provider struct {
	responses       map[string]string
	defaultResponse string
	responseTime    time.Duration
	errorRate       float64
}

// Fixed dataset for each provider
var providers = []*Provider{
	{
		responses: map[string]string{
			"What is your name?": "I am Provider 1.",
			"How are you?":       "Provider 1 is functioning well.",
			"Hi":                 "Hello from Provider 1",
		},
		defaultResponse: "Provider 1: I'm sorry, I don't understand the question.",
	},
	{
		responses: map[string]string{
			"What is your name?": "I am Provider 2.",
			"How are you?":       "Provider 2 is operational.",
			"Hi":                 "Greetings from Provider 2",
		},
		defaultResponse: "Provider 2: I cannot process that request.",
	},
	{
		responses: map[string]string{
			"What is your name?": "I am Provider 3.",
			"How are you?":       "Provider 3 is up and running.",
			"Hi":                 "Hi from Provider 3",
		},
		defaultResponse: "Provider 3: Sorry, I don't have a response for that.",
	},
}

var currentProviderIndex = 0
var providerMu sync.Mutex

// GetResponse from the current provider, monitor performance, and switch if needed
func GetResponse(prompt string) (string, error) {
	providerMu.Lock()
	defer providerMu.Unlock()

	for attempts := 0; attempts < len(providers); attempts++ {
		provider := providers[currentProviderIndex]
		start := time.Now()

		response, ok := provider.responses[prompt]
		provider.responseTime = time.Since(start)

		if !ok {
			// If key does not match, return the default response immediately
			return provider.defaultResponse, nil
		}

		// Check if the response time exceeds the acceptable threshold
		if provider.responseTime > 2*time.Second {
			provider.errorRate += 1.0
			if shouldSwitchProvider() {
				switchProvider()
				continue
			} else {
				return "", errors.New("response time exceeded acceptable limits")
			}
		}

		// Reset the error rate if the provider responded within the acceptable time
		provider.errorRate = 0.0
		return response, nil
	}

	return "", errors.New("all providers failed")
}

// Decide whether to switch providers based on performance metrics
func shouldSwitchProvider() bool {
	provider := providers[currentProviderIndex]
	return provider.errorRate > 0.5
}

// Switch to the next provider in the list
func switchProvider() {
	currentProviderIndex = (currentProviderIndex + 1) % len(providers)
}
