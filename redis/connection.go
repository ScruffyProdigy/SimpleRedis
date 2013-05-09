package redis

import (
	"net"
	"strings"
)

//A Connection is a single connection to a Redis Instance.
//Each client typically has a pool of these to work with
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
		command.callback()(nil)
		return err
	}

	return command.callback()(res)
}

//Error is how an error gets reported.
//Since The redis code operates in a separate goroutine, errors can't always be reported directly
func (this Connection) Error(e error, c command) {
	this.client.errCallback(e, strings.Join(c.arguments(), " "))
}

//Execute allows a command to be executed on a specific connection
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
