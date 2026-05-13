package auth

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/term"
)

func defaultPrompt(prompt string, confirm bool) (string, error) {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		return "", errors.New("passphrase required: no terminal attached")
	}

	fmt.Fprint(os.Stderr, prompt)
	first, err := term.ReadPassword(fd)
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", err
	}
	if !confirm {
		return string(first), nil
	}

	fmt.Fprint(os.Stderr, "Confirm passphrase: ")
	second, err := term.ReadPassword(fd)
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", err
	}
	if string(first) != string(second) {
		return "", ErrPassphraseMismatch
	}
	return string(first), nil
}
