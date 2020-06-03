package cmd_utils

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

func RequestPasswordWithRepeat() (string, error) {
	var err error
	fmt.Printf("New password: ")
	newPass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}

	fmt.Printf("Repeat password: ")
	newPass2, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}

	if string(newPass) != string(newPass2) {
		return "", fmt.Errorf("passwords does mot match")
	}

	return string(newPass), err
}


func RequestPassword() (string, error) {
	var err error
	fmt.Printf("Password: ")
	newPass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}

	return string(newPass), err
}
