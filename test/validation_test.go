package main

import (
	"fmt"
	"testing"

	"git.garena.com/jiayu.li/entry-task/cmd/protocol"
)

func checkValidUsername(t *testing.T, username string) {
	want := true
	if got := protocol.IsValidUsername(username); got != want {
		t.Errorf("isValidUsername(\"%s\") = %t, want %t", username, got, want)
	}
}

func Test_valid_username(t *testing.T) {
	validNames := [6]string{
		"abcd",
		"a234",
		"a2345678901234567890",
		"abcd123a-b",
		"A123_abc",
		"A_-_abc",
	}
	for _, name := range validNames {
		t.Run(fmt.Sprintf("Check valid username %s", name), func(t *testing.T) {
			checkValidUsername(t, name)
		})
	}
}

func checkInvalidUsername(t *testing.T, username string) {
	want := false
	if got := protocol.IsValidUsername(username); got != want {
		t.Errorf("isValidUsername(\"%s\") = %t, want %t", username, got, want)
	}
}

func Test_invalid_username(t *testing.T) {
	invalidNames := [6]string{
		"abc",
		"1234",
		"a12345678901234567890",
		"-abcde",
		"_1234",
		"A_-,123",
	}
	for _, name := range invalidNames {
		t.Run(fmt.Sprintf("Check invalid username %s", name), func(t *testing.T) {
			checkInvalidUsername(t, name)
		})
	}
}

func checkValidPassword(t *testing.T, password string) {
	want := true
	if got := protocol.IsValidPassword(password); got != want {
		t.Errorf("isValidPassword(\"%s\") = %t, want %t", password, got, want)
	}
}

func Test_valid_password(t *testing.T) {
	validPasswords := [6]string{
		"abcd",
		"f2%^@fds#FW",
		"a2345678901234567890",
		"abcd123a-b",
		"AFEH3$>,32sabc",
		"A_-_abc",
	}
	for _, pass := range validPasswords {
		t.Run(fmt.Sprintf("Check valid password %s", pass), func(t *testing.T) {
			checkValidPassword(t, pass)
		})
	}
}

func checkInvalidPassword(t *testing.T, password string) {
	want := false
	if got := protocol.IsValidPassword(password); got != want {
		t.Errorf("isValidPassword(\"%s\") = %t, want %t", password, got, want)
	}
}

func Test_invalid_password(t *testing.T) {
	invalidPasswords := [3]string{
		"abc",
		"fasf2afcasvf2%^@fds#FW",
		"a12345678901234567890",
	}
	for _, pass := range invalidPasswords {
		t.Run(fmt.Sprintf("Check invalid password %s", pass), func(t *testing.T) {
			checkInvalidPassword(t, pass)
		})
	}
}
