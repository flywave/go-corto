package corto

type PredictionType uint32

const (
	PREDICTION_DIFF      PredictionType = 0x0
	PREDICTION_ESTIMATED PredictionType = 0x1
	PREDICTION_BORDER    PredictionType = 0x2
)

type FormatType uint32

const (
	FORMAT_UINT32 FormatType = 0
	FORMAT_INT32  FormatType = 1
	FORMAT_UINT16 FormatType = 2
	FORMAT_INT16  FormatType = 3
	FORMAT_UINT8  FormatType = 4
	FORMAT_INT8   FormatType = 5
	FORMAT_FLOAT  FormatType = 6
	FORMAT_DOUBLE FormatType = 7
)

type StrategyType uint32

const (
	PARALLEL   StrategyType = 0x1
	CORRELATED StrategyType = 0x2
)

type EntropyType uint32

const (
	ENTROPY_NONE     EntropyType = 0
	ENTROPY_TUNSTALL EntropyType = 1
)
