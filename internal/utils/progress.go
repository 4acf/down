package utils

import (
	"fmt"
)

type ProgressBar struct {
	inProgressMessage string
	completeMessage   string
	iterations        int
}

const (
	PROGRESS_BAR_LENGTH = 30
)

func NewProgressBar(inProgressMessage string, completeMessage string, iterations int) ProgressBar {
	return ProgressBar{
		inProgressMessage: inProgressMessage,
		completeMessage:   completeMessage,
		iterations:        iterations,
	}
}

func (p *ProgressBar) UpdateConsole(finishedIteration int) {

	fraction := float64(finishedIteration) / float64(p.iterations)

	filledLength := int(float64(PROGRESS_BAR_LENGTH) * fraction)

	bar := fmt.Sprintf("%s%s", repeat("=", filledLength), repeat("-", PROGRESS_BAR_LENGTH-filledLength))

	layout := fmt.Sprintf("[%s] %.2f%% (%d/%d)", bar, fraction*100, finishedIteration, p.iterations)

	var message string
	if finishedIteration == p.iterations {
		message = p.completeMessage
	} else {
		message = p.inProgressMessage
	}

	fmt.Printf("\r%s %s", message, layout)

}

func repeat(str string, n int) string {
	result := ""
	for range n {
		result += str
	}
	return result
}
