package irc

import (
	"bytes"
	"strings"
)

type message struct {
	prefix     *string
	command    string
	parameters [15]*string
}

func parseMessage(line string) message {
	lineSplit := strings.Split(line, " ")
	lineSplitIndex := 0

	var prefix *string
	var command string
	var parameters [15]*string

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

		parametersIndex := index - lineSplitIndex
		parameters[parametersIndex] = new(string)

		if parameter[0] != ':' {
			*parameters[parametersIndex] = parameter
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

			*parameters[parametersIndex] = builder.String()
			break
		}
	}

	return message{
		prefix:     prefix,
		command:    command,
		parameters: parameters,
	}
}

func (message *message) serialize() string {
	var buffer bytes.Buffer

	if message.prefix != nil {
		buffer.WriteString(*message.prefix)
		buffer.WriteByte(' ')
	}

	buffer.WriteString(message.command)

	for _, parameter := range message.parameters {
		if parameter != nil {
			buffer.WriteByte(' ')
			buffer.WriteString(*parameter)
		} else {
			break
		}
	}

	return buffer.String()
}
