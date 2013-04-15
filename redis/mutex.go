package redis

import (
	"strconv"
)

//This mutex is useful for making sure that network resources are not being used by any other servers
type Mutex interface {
	Try(action func()) bool
	Force(action func())
}

type redisMutex struct {
	init      String
	processes List
}

func newMutex(client Prefix, key string, count int) Mutex {
	m := new(redisMutex)
	m.init = client.String(key + ":Initialized")
	m.processes = client.List(key)
	m.initialize(count)

	return m
}

func (this *redisMutex) initialize(count int) {
	if <-this.init.Replace("initialized") != "initialized" {
		for i := 0; i < count; i++ {
			this.processes.RightPush(strconv.Itoa(i + 1))
		}
	}
}

func (this *redisMutex) Try(action func()) bool {
	val, available := <-this.processes.LeftPop()
	if !available {
		return false
	}

	defer func() {
		<-this.processes.RightPush(val)
	}()

	action()

	return true
}

func (this *redisMutex) Force(action func()) {
	val := <-this.processes.BlockUntilLeftPop()

	defer func() {
		<-this.processes.RightPush(val)
	}()

	action()
}
