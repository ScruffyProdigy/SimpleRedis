package redis

import ()

type pipe struct {
	commands    []command
	errCallback errCallback
}

func (this *pipe) Execute(command command) error {
	this.commands = append(this.commands, command)
	return nil
}

func (this *pipe) ErrCallback(err error, s string) {
	this.errCallback.Call(err, s)
}

func (this Client) piping(callback func(Executor) bool, queued bool) {
	p := new(pipe)
	p.commands = make([]command, 0, 5)
	p.errCallback = this.errCallback
	var result bool
	defer func() {
		var bundle []byte
		for _, command := range p.commands {
			comm, err := buildCommand(command.arguments())
			if err != nil {
				this.errCallback.Call(err, "piping")
			}
			bundle = append(bundle, comm...)
		}
		this.useConnection(func(c *Connection) {
			c.Write(bundle)
			if !result {
				//everything was discarded - just get basic result and don't bother waiting for everything else
				getResponse(c)
				return
			}
			if queued {
				//get rid of all of the "queued" responses
				for i := 0; i < len(p.commands)-1; i++ {
					getResponse(c)
				}
				//the first reply is going to be a multi-bulk, with all of the other replies as subresponses
				//get rid of the multi-bulk, and just get the other replies as normal
				//(this is a little bit hacky, perhaps I'll make it less so in future versions)
				getString(c)
				p.commands = p.commands[1 : len(p.commands)-1]
			}
			for _, command := range p.commands {
				c.output(command)
			}
		})
	}()
	result = callback(p)
}

func (this Client) Pipeline(callback func(Executor)) {
	this.piping(func(e Executor) bool {
		callback(e)
		return true
	}, false)
}

func (this Client) Transaction(callback func(Executor)) {
	multi, _ := newNilCommand([]string{"MULTI"})
	exec, _ := newNilCommand([]string{"EXEC"})
	discard, _ := newNilCommand([]string{"DISCARD"})
	this.piping(func(p Executor) (result bool) {
		p.Execute(multi)
		defer func() {
			rec := recover()
			if rec == nil {
				print("Executed!\n")
				p.Execute(exec)
			} else {
				print("Discarded!\n")
				p.Execute(discard)
				result = false
			}
		}()

		callback(p)
		return true
	}, true)
}
