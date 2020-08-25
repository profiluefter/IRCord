package irc

import (
	"reflect"
	"testing"
)

func TestParseMessageWithoutParameters(t *testing.T) {
	actual := parseMessage(":prefix command")

	var expectedPrefix = "prefix"
	expected := message{
		prefix:  &expectedPrefix,
		command: "command",
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fail()
	}
}

func TestParseMessage(t *testing.T) {
	actual := parseMessage("command1234 parameter1 parameter2")

	var (
		parameter1 = "parameter1"
		parameter2 = "parameter2"
	)
	expected := message{
		command: "command1234",
		parameters: [15]*string{
			&parameter1,
			&parameter2,
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fail()
	}
}

func TestParseMessageWithTrailingParameter(t *testing.T) {
	actual := parseMessage("command1234 parameter1 :parameter2 with spaces")

	var (
		parameter1 = "parameter1"
		parameter2 = "parameter2 with spaces"
	)
	expected := message{
		command: "command1234",
		parameters: [15]*string{
			&parameter1,
			&parameter2,
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fail()
	}
}

func TestParseMessageWithOnlyCommand(t *testing.T) {
	actual := parseMessage("command")
	expected := message{
		command: "command",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fail()
	}
}
