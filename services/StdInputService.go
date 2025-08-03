package services

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type InputService interface {
	AskUserYN(prompt string) (bool, error)
	GetPipedInput() (string, bool)
}

type StdInputService struct {
	reader io.Reader
}

func NewStdInputService() *StdInputService {
	return &StdInputService{reader: os.Stdin}
}

func (s *StdInputService) AskUserYN(prompt string) (bool, error) {
	fmt.Print(prompt + " [y/n]: ")
	scanner := bufio.NewScanner(s.reader)
	if scanner.Scan() {
		answer := scanner.Text()
		answer = strings.TrimSpace(strings.ToLower(answer))
		return answer == "y" || answer == "yes", nil
	}
	return false, scanner.Err()
}

func (s *StdInputService) GetPipedInput() (string, bool) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", false
	}
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", false
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", false
	}
	clean := strings.TrimSpace(string(data))
	if clean == "" {
		return "", false
	}
	return clean, true
}
