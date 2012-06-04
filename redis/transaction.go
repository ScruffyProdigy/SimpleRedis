package redis

import ()

type pipe struct {
	commands []command
}

func (this *pipe) Execute(command command) {
	this.commands = append(this.commands, command)
}

func (this Client) piping(callback func(Executor), queued bool) {
	p := new(pipe)
	p.commands = make([]command, 0, 5)
	defer func() {
		var bundle []byte
		for _, command := range p.commands {
			bundle = append(bundle, buildCommand(command.arguments())...)
		}
		this.useConnection(func(c Connection) {
			c.Write(bundle)
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
	callback(p)
}

func (this Client) Pipeline(callback func(Executor)) {
	this.piping(callback, false)
}

func (this Client) Transaction(callback func(Executor)) {
	multi, _ := newNilCommand([]string{"MULTI"})
	exec, _ := newNilCommand([]string{"EXEC"})
	discard, _ := newNilCommand([]string{"DISCARD"})
	this.piping(func(p Executor) {
		p.Execute(multi)
		defer func() {
			rec := recover()
			if rec == nil {
				p.Execute(exec)
			} else {
				p.Execute(discard)
			}
		}()

		callback(p)
	}, true)
}
