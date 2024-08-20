package strategies

import (
	"fmt"
	"math/rand"

	backtest_types "goquant/pkg/backtest"

	"github.com/go-gota/gota/dataframe"
)

// MarkovChainStrategy holds the transition matrix and implements the Markov Chain strategy.
type MarkovChainStrategy struct {
	TransitionMatrix map[string]map[string]float64
	States           []string
}

// NewMarkovChainStrategy initializes a new MarkovChainStrategy.
func NewMarkovChainStrategy() *MarkovChainStrategy {
	return &MarkovChainStrategy{
		TransitionMatrix: make(map[string]map[string]float64),
		States:           []string{"Up", "Down", "Unchanged"},
	}
}

// Build constructs the transition matrix based on the entire historical dataframe.
func (mcs *MarkovChainStrategy) Build(df dataframe.DataFrame) {
	mcs.TransitionMatrix = calculateTransitionMatrix(df, mcs.States)
}

// Run applies the Markov Chain strategy to predict the next state and generate trading signals.
func (mcs *MarkovChainStrategy) Run(df dataframe.DataFrame) backtest_types.StrategyAction {
	if df.Nrow() < 2 {
		return "Hold"
	}
	// Get the current state
	currentState := getCurrentState(df)

	// Predict the next state
	nextState := predictNextState(currentState, mcs.States, mcs.TransitionMatrix)

	// Generate a signal based on the predicted next state
	switch nextState {
	case "Up":
		return "Buy"
	case "Down":
		return "Sell"
	default:
		return "Hold"
	}
}

// calculateTransitionMatrix calculates the transition matrix from historical data.
func calculateTransitionMatrix(df dataframe.DataFrame, states []string) map[string]map[string]float64 {
	transitionMatrix := make(map[string]map[string]float64)

	for _, state := range states {
		transitionMatrix[state] = make(map[string]float64)
		for _, nextState := range states {
			transitionMatrix[state][nextState] = 0
		}
	}

	totalTransitions := make(map[string]int)
	prices := df.Col("Close").Float()

	for i := 1; i < len(prices)-1; i++ { // Adjusted to avoid out-of-bounds error
		currentState := getState(prices[i-1], prices[i])
		nextState := getState(prices[i], prices[i+1])

		transitionMatrix[currentState][nextState]++
		totalTransitions[currentState]++
	}

	// Normalize the transition matrix to get probabilities
	for _, state := range states {
		for _, nextState := range states {
			if totalTransitions[state] > 0 {
				transitionMatrix[state][nextState] /= float64(totalTransitions[state])
			}
		}
	}
	fmt.Println(transitionMatrix)
	return transitionMatrix
}

// getCurrentState determines the current state based on the last two price points.
func getCurrentState(df dataframe.DataFrame) string {
	prices := df.Col("Close").Float()
	return getState(prices[len(prices)-2], prices[len(prices)-1])
}

// getState determines the state based on price movements.
func getState(prevPrice, currPrice float64) string {
	if currPrice > prevPrice {
		return "Up"
	} else if currPrice < prevPrice {
		return "Down"
	}
	return "Unchanged"
}

// predictNextState predicts the next state based on the current state and the transition matrix.
func predictNextState(currentState string, states []string, transitionMatrix map[string]map[string]float64) string {
	prob := rand.Float64()
	cumulativeProb := 0.0

	for _, nextState := range states {
		cumulativeProb += transitionMatrix[currentState][nextState]
		if prob <= cumulativeProb {
			return nextState
		}
	}

	// Default to current state if prediction fails
	return currentState
}
