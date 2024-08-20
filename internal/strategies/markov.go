package strategies

import (
	"fmt"
	"math/rand"
	"strings"

	backtest_types "goquant/pkg/backtest"

	"github.com/go-gota/gota/dataframe"
)

// MarkovChainStrategy holds the transition matrix and implements the Markov Chain strategy.
type MarkovChainStrategy struct {
	TransitionMatrix map[string]map[string]float64
	States           []string
	Depth            int
}

// NewMarkovChainStrategy initializes a new MarkovChainStrategy with a specified depth.
func NewMarkovChainStrategy(depth int) *MarkovChainStrategy {
	return &MarkovChainStrategy{
		TransitionMatrix: make(map[string]map[string]float64),
		States:           []string{"Up", "Down", "Unchanged"},
		Depth:            depth,
	}
}

// Build constructs the transition matrix based on the entire historical dataframe.
func (mcs *MarkovChainStrategy) Build(df dataframe.DataFrame) {
	mcs.TransitionMatrix = calculateTransitionMatrix(df, mcs.States, mcs.Depth)
}

// Run applies the Markov Chain strategy to predict the next state and generate trading signals.
func (mcs *MarkovChainStrategy) Run(df dataframe.DataFrame) backtest_types.StrategyAction {
	if df.Nrow() < mcs.Depth {
		return "Hold"
	}
	// Get the current sequence of states
	currentSequence := getCurrentSequence(df, mcs.Depth)

	// Predict the next state
	nextState := predictNextState(currentSequence, mcs.States, mcs.TransitionMatrix)

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

// calculateTransitionMatrix calculates the transition matrix from historical data with depth n.
func calculateTransitionMatrix(df dataframe.DataFrame, states []string, depth int) map[string]map[string]float64 {
	transitionMatrix := make(map[string]map[string]float64)

	for _, state := range states {
		for _, nextState := range states {
			transitionMatrix[state] = make(map[string]float64)
			fmt.Printf("State: %s, Next State: %s\n", state, nextState)
		}
	}

	prices := df.Col("Close").Float()
	totalTransitions := make(map[string]int)

	// Build the transition matrix
	for i := depth; i < len(prices)-1; i++ {
		currentSequence := getStateSequence(prices[i-depth : i])
		nextState := getState(prices[i], prices[i+1])

		if transitionMatrix[currentSequence] == nil {
			transitionMatrix[currentSequence] = make(map[string]float64)
		}
		transitionMatrix[currentSequence][nextState]++
		totalTransitions[currentSequence]++
	}

	// Normalize the transition matrix to get probabilities
	for sequence := range transitionMatrix {
		for _, nextState := range states {
			if totalTransitions[sequence] > 0 {
				transitionMatrix[sequence][nextState] /= float64(totalTransitions[sequence])
			}
		}
	}

	return transitionMatrix
}

// getCurrentSequence determines the current sequence of states based on the last n price points.
func getCurrentSequence(df dataframe.DataFrame, depth int) string {
	prices := df.Col("Close").Float()
	return getStateSequence(prices[len(prices)-depth:])
}

// getStateSequence converts a slice of prices into a sequence of states.
func getStateSequence(prices []float64) string {
	states := []string{}
	for i := 1; i < len(prices); i++ {
		states = append(states, getState(prices[i-1], prices[i]))
	}
	return strings.Join(states, "-")
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

// predictNextState predicts the next state based on the current sequence and the transition matrix.
func predictNextState(currentSequence string, states []string, transitionMatrix map[string]map[string]float64) string {
	prob := rand.Float64()
	cumulativeProb := 0.0

	for _, nextState := range states {
		cumulativeProb += transitionMatrix[currentSequence][nextState]
		if prob <= cumulativeProb {
			return nextState
		}
	}

	// Default to current state if prediction fails
	return strings.Split(currentSequence, "-")[0] // Return the first state in the sequence
}

// predictNextStateMax predicts the next state with the highest probability based on the current sequence and the transition matrix.
//
// currentSequence - The current sequence of states.
// states - A list of possible next states.
// transitionMatrix - A map of state transitions with their corresponding probabilities.
// Returns the next state with the highest probability.
func predictNextStateMax(currentSequence string, states []string, transitionMatrix map[string]map[string]float64) string {
	maxProb := 0.0
	bestState := states[0]

	for _, nextState := range states {
		if transitionMatrix[currentSequence][nextState] > maxProb {
			maxProb = transitionMatrix[currentSequence][nextState]
			bestState = nextState
		}
	}

	return bestState
}
