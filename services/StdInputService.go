package services

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type InputService interface {
	AskUser(prompt string, options []string) (string, error)
	GetPipedInput() (string, bool)
}

type StdInputService struct {
	reader io.Reader
}

func NewStdInputService() *StdInputService {
	return &StdInputService{reader: os.Stdin}
}

func (s *StdInputService) AskUser(prompt string, options []string) (string, error) {
	// opening straight terminal unix interface
	tty, err := os.Open("/dev/tty")
	if err != nil {
		return "", fmt.Errorf("failed to open /dev/tty: %w", err)
	}
	defer func() {
		if err := tty.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error closing tty: %v\n", err)
		}
	}()

	if _, err := fmt.Fprintf(os.Stdout, "%s [%s]: ", prompt, strings.Join(options, "/")); err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(tty)
	if scanner.Scan() {
		answer := strings.TrimSpace(strings.ToLower(scanner.Text()))

		for _, opt := range options {
			if answer == strings.ToLower(opt) {
				return answer, nil
			}
		}

		return "", fmt.Errorf("invalid input: '%s' (allowed: %s)", answer, strings.Join(options, ", "))
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read from /dev/tty: %w", err)
	}
	return "", fmt.Errorf("no input provided")
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
