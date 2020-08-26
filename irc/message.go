package irc

import (
	"bytes"
	"errors"
	"strings"
)

type message struct {
	prefix     *string
	command    string
	parameters []*string
}

func parseMessage(line string) (*message, error) {
	lineSplit := strings.Split(line, " ")
	lineSplitIndex := 0

	if len(lineSplit) <= 0 {
		return nil, errors.New("invalid message")
	}

	var prefix *string
	var command string
	var parameters []*string

	if line[0] == ':' {
		prefix = new(string)
		*prefix = strings.TrimPrefix(lineSplit[lineSplitIndex], ":")
		lineSplitIndex++
	}

	command = lineSplit[lineSplitIndex]
	lineSplitIndex++

	for index, parameter := range lineSplit {
		if index < lineSplitIndex {
			continue
		}

		parameterPointer := new(string)
		parameters = append(parameters, parameterPointer)

		if parameter[0] != ':' {
			*parameterPointer = parameter
		} else {
			//trailing parameter: everything following this is one parameter
			var builder strings.Builder
			builder.WriteString(strings.TrimPrefix(parameter, ":"))

			for innerIndex, innerParameter := range lineSplit {
				if innerIndex <= index {
					continue
				}
				builder.WriteString(" ") //add space from split
				builder.WriteString(innerParameter)
			}

			*parameterPointer = builder.String()
			break
		}
	}

	m := message{
		prefix:     prefix,
		command:    command,
		parameters: parameters,
	}
	return &m, nil
}

func (message message) serialize() string {
	var buffer bytes.Buffer

	if message.prefix != nil {
		buffer.WriteRune(':')
		buffer.WriteString(*message.prefix)
		buffer.WriteRune(' ')
	}

	buffer.WriteString(message.command)

	for _, parameter := range message.parameters {
		if parameter != nil {
			buffer.WriteRune(' ')
			buffer.WriteString(*parameter)
		} else {
			//This shouldn't happen anymore but better be save
			break
		}
	}

	return buffer.String()
}
