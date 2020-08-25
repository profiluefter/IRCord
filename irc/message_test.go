package irc

import (
	"reflect"
	"testing"
)

func TestParseMessage(t *testing.T) {
	actual := parseMessage("command")
	expected := message{
		command: "command",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fail()
	}
}

func TestParseMessageWithPrefix(t *testing.T) {
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

func TestParseMessageWithParameters(t *testing.T) {
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

func TestSerialize(t *testing.T) {
	actual := message{
		command: "test",
	}.serialize()

	expected := "test"

	if actual != expected {
		t.Fatalf("%s != %s", actual, expected)
	}
}

func TestSerializeWithPrefix(t *testing.T) {
	prefix := "testrunner"
	actual := message{
		prefix:  &prefix,
		command: "test",
	}.serialize()

	expected := ":testrunner test"

	if actual != expected {
		t.Fatalf("%s != %s", actual, expected)
	}
}

func TestSerializeWithParameters(t *testing.T) {
	firstParameter := "param1"
	secondParameter := "param2"
	actual := message{
		command: "test",
		parameters: [15]*string{
			&firstParameter,
			&secondParameter,
		},
	}.serialize()

	expected := "test param1 param2"

	if actual != expected {
		t.Fatalf("%s != %s", actual, expected)
	}
}
