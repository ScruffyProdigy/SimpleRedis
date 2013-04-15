package redis

type ReadWriteMutex struct {
	readers int
	write   Mutex
	Read    Mutex
	Write   Mutex
}

type writeMutex struct {
	*ReadWriteMutex
}

func lockAllReads(rw *ReadWriteMutex, finalAction func()) func() {
	return func() {
		in := make(chan bool)
		out := make(chan bool)
		for i := 0; i < rw.readers; i++ {
			go func(j int) {
				rw.Read.Force(func() {
					in <- true
					<-out //this one is here to make sure the final action has been completed before releasing any of the readers
				})
				<-out //this one is here to make sure all of the readers have been released before returning control to the calling function
			}(i)
		}

		for i := 0; i < rw.readers; i++ {
			<-in
		}
		finalAction()
		for i := 0; i < 2*rw.readers; i++ {
			out <- true
		}
	}
}

func (this writeMutex) Try(action func()) bool {
	return this.write.Try(lockAllReads(this.ReadWriteMutex, action))
}

func (this writeMutex) Force(action func()) {
	this.write.Force(lockAllReads(this.ReadWriteMutex, action))
}

func newRWMutex(client Prefix, key string, readers int) *ReadWriteMutex {
	rw := new(ReadWriteMutex)
	rw.readers = readers
	rw.write = client.Mutex(key + ":Write")
	rw.Read = client.Semaphore(key+":Read", readers)
	rw.Write = writeMutex{rw}
	return rw
}
