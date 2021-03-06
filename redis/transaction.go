package redis

import ()

type pipe struct {
	commands     []command
	fErrCallback errCallbackFunc
}

func (this *pipe) Execute(command command) {
	this.commands = append(this.commands, command)
}

func (this *pipe) errCallback(err error, s string) {
	this.fErrCallback.Call(err, s)
}

func (this Client) piping(callback func(SafeExecutor) bool, queued bool) {
	p := new(pipe)
	p.commands = make([]command, 0, 5)
	p.fErrCallback = this.fErrCallback
	var result bool
	defer func() {
		var bundle []byte
		for _, command := range p.commands {
			comm, err := buildCommand(command.arguments())
			if err != nil {
				this.errCallback(err, "piping")
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

//Pipeline creates an Executor that will force every command issued on it to be sent at the same time (thus saving on network costs).
//It waits until the end of the function to execute them
func (this Client) Pipeline(callback func(SafeExecutor)) {
	this.piping(func(e SafeExecutor) bool {
		callback(e)
		return true
	}, false)
}

//Transaction creates an Executor that will tell redis to queue all of the commands and complete them atomically
//(this prevents other clients from issuing commands in between yours)
func (this Client) Transaction(callback func(SafeExecutor)) {
	this.piping(func(p SafeExecutor) (result bool) {
		NilCommand(p, "MULTI")
		defer func() {
			rec := recover()
			if rec == nil {
				NilCommand(p, "EXEC")
			} else {
				NilCommand(p, "DISCARD")
				result = false
			}
		}()

		callback(p)
		return true
	}, true)
}
