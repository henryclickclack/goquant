package strategies

import (
	backtest_types "goquant/pkg/backtest"

	"github.com/go-gota/gota/dataframe"
)

// StrategyFunc is a function type that represents a trading strategy.
type StrategyFunc func(df dataframe.DataFrame) backtest_types.StrategyAction

// EnsembleStrategy represents a strategy that combines multiple strategies.
type EnsembleStrategy struct {
	Strategies []StrategyFunc
	Weights    []float64
}

// NewEnsembleStrategy creates a new EnsembleStrategy instance.
//
// Parameters:
// - strategies: A slice of StrategyFunc representing the individual strategies to combine.
// - weights: A slice of float64 representing the weights for each strategy. If empty or nil, equal weights will be used.
// Returns a pointer to the newly created EnsembleStrategy.
func NewEnsembleStrategy(strategies []StrategyFunc, weights []float64) *EnsembleStrategy {
	// If no weights are provided, assign equal weights to all strategies.
	if len(weights) == 0 || weights == nil {
		weights = make([]float64, len(strategies))
		for i := range weights {
			weights[i] = 1.0 / float64(len(strategies))
		}
	}

	return &EnsembleStrategy{
		Strategies: strategies,
		Weights:    weights,
	}
}

// Run applies the ensemble strategy to make a decision based on the combined strategies.
func (es *EnsembleStrategy) Run(df dataframe.DataFrame) backtest_types.StrategyAction {
	actionScores := map[backtest_types.StrategyAction]float64{
		"Buy":  0,
		"Sell": 0,
		"Hold": 0,
	}

	// Aggregate the weighted decisions from each strategy.
	for i, strategy := range es.Strategies {
		action := strategy(df)
		actionScores[action] += es.Weights[i]
	}

	// Determine the action with the highest score. if mmultiple actions have the same highest score, choose hold first
	var finalAction backtest_types.StrategyAction
	maxScore := -1.0
	for action, score := range actionScores {
		if score > maxScore {
			maxScore = score
			finalAction = action
		}
	}

	return finalAction
}
