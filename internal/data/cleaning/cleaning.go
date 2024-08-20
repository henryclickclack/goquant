package cleaning

import (
	"math"
	"math/rand"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

type FillStrategy int

const (
    FillWithNull FillStrategy = iota
    FillWithMean
    FillWithMedian
)

// RemoveOutliers removes data points that are outside of the mean ± 3*stdDev for each numerical column in a dataframe.
//
// Parameters:
//   df (dataframe.DataFrame): The input dataframe.
// Returns:
//   dataframe.DataFrame: The dataframe with outliers removed.
func RemoveOutliers(df dataframe.DataFrame) dataframe.DataFrame {
    // Iterate over the dataframe and remove rows where any value is outside of the mean ± 3*stdDev
    for _, colName := range df.Names() {
        if df.Col(colName).Type() == series.Float || df.Col(colName).Type() == series.Int {
            // Calculate mean and stdDev for the column
            col := df.Col(colName)
            meanValue := col.Mean()
            stdDevValue := col.StdDev()

            // Calculate bounds
            lowerBound := meanValue - 2.0*stdDevValue
            upperBound := meanValue + 2.0*stdDevValue

            // Filter out the outliers
            df = df.Filter(
                dataframe.F{
                    Colname:    colName,
                    Comparator: series.GreaterEq,
                    Comparando: lowerBound,
                },
            ).Filter(
                dataframe.F{
                    Colname:    colName,
                    Comparator: series.LessEq,
                    Comparando: upperBound,
                },
            )
        }
    }

    return df
}


// FillMissingData fills missing data in a DataFrame using forward fill.
//
// Parameters:
//   df (dataframe.DataFrame): The input DataFrame.
// Returns:
//   dataframe.DataFrame: The DataFrame with missing data filled.
func FillMissingData(df dataframe.DataFrame, strategy FillStrategy, fillValue float64) dataframe.DataFrame {
    for _, colName := range df.Names() {
        col := df.Col(colName)

        if col.Type() == series.Float { // We only handle float columns here
            var filled series.Series

            switch strategy {
            case FillWithNull:
                // Replace NaN with a specific value (fillValue)
                filled = col.Map(func(val series.Element) series.Element {

                    if val.IsNA(){
						val.Set(fillValue)
                    }
                    return val
                })

            case FillWithMean:
                // Replace NaN with the mean of the column
                meanValue := col.Mean()
				filled = col.Map(func(val series.Element) series.Element {

					if val.IsNA(){
						val.Set(meanValue)
					
					}
					return val
				})

            case FillWithMedian:
                // Replace NaN with the median of the column
                medianValue := col.Median()
				filled = col.Map(func(val series.Element) series.Element {

					if val.IsNA(){
						val.Set(medianValue)
					
					}
					return val
				})
            }

            df = df.Mutate(filled)
        }
    }
    return df
}


// InsertOutlierAndNanTest inserts outliers and NaN values into a DataFrame for testing purposes.
//
// Parameters:
//   df (dataframe.DataFrame): The input DataFrame.
// Returns:
//   dataframe.DataFrame: The DataFrame with outliers and NaN values inserted.
func InsertOutlierAndNanTest(df dataframe.DataFrame) dataframe.DataFrame {
    // Iterate over each column
    for _, colName := range df.Names() {
        col := df.Col(colName)

        if col.Type() == series.Float { // We only modify float columns
            numRows := col.Len()
            data := col.Float()

            // Insert a few NaN values at random positions
            for i := 0; i < numRows/10; i++ { // Add NaN to 10% of the rows
                randIdx := rand.Intn(numRows)
                data[randIdx] = math.NaN()
            }

            // Insert a few outliers at random positions
            maxValue := col.Max()
            minValue := col.Min()
            outlierValueHigh := maxValue * 10 // A very high outlier
            outlierValueLow := minValue / 10  // A very low outlier

            for i := 0; i < numRows/10; i++ { // Add outliers to 10% of the rows
                randIdx := rand.Intn(numRows)
                if rand.Float64() > 0.5 {
                    data[randIdx] = outlierValueHigh
                } else {
                    data[randIdx] = outlierValueLow
                }
            }

            // Create a new series with the modified data
            modifiedSeries := series.New(data, series.Float, colName)

            // Update the dataframe with the modified column
            df = df.Mutate(modifiedSeries)

        }
    }

    return df
}