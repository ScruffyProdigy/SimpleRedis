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
	i := 0
	var lockNextRead func()
	lockNextRead = func() {
		if i < rw.readers {
			i++
			rw.Read.Force(lockNextRead)
		} else {
			finalAction()
		}
	}
	return lockNextRead
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
