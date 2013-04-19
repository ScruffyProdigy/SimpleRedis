package redis

import (
	"net"
	"strings"
)

type Connection struct {
	net.Conn
	id     int
	client *Client
}

func (this Connection) input(command command) error {
	comm, err := buildCommand(command.arguments())
	if err != nil {
		return err
	}

	_, err = this.Write(comm)
	return err
}

func (this Connection) output(command command) error {
	res, err := getResponse(this)
	if err != nil {
		return err
	}

	return command.callback()(res)
}

func (this Connection) Error(e error, c command) {
	this.client.errCallback(e, strings.Join(c.arguments(), " "))
}

func (this Connection) Execute(command command) {
	err := this.input(command)
	if err != nil {
		this.Error(err, command)
		return
	}

	err = this.output(command)
	if err != nil {
		this.Error(err, command)
	}
}
