package assinador

import "time"

type Operation string

const (
	OperationSign     Operation = "sign"
	OperationValidate Operation = "validate"
)

type LocalRequest struct {
	Operation     Operation
	InputPath     string
	OutputPath    string
	SignaturePath string
	JavaBin       string
	JarPath       string
	Timeout       time.Duration
}

type LocalResult struct {
	Stdout string
	Stderr string
	Code   int
}

type UsageError struct {
	Message string
	Code    int
}

func (e UsageError) Error() string {
	return e.Message
}
