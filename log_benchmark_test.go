package lorg

import (
	"bytes"
	stdlog "log"
	"sync"
	"testing"
)

func BenchmarkLog_Printf_Parallel(b *testing.B) {
	const logString = "lorg"
	var buffer bytes.Buffer

	log := NewLog()
	log.SetFormat(NewFormat("%s"))
	log.SetOutput(&buffer)
	log.format.Reset()

	bufferMutex := &sync.Mutex{}

	log.setMutex(bufferMutex)

	b.RunParallel(func(pb *testing.PB) {
		bufferMutex.Lock()
		buffer.Reset()
		bufferMutex.Unlock()
		for pb.Next() {
			log.Printf("%v", logString)
		}
	})
}

func BenchmarkLog_Printf_NoFormat(b *testing.B) {
	const logString = "lorg"
	var buffer bytes.Buffer

	log := NewLog()
	log.SetOutput(&buffer)
	log.SetFormat(NewFormat("%s"))

	for i := 0; i < b.N; i++ {
		buffer.Reset()
		log.Printf("%v", logString)
	}
}

func BenchmarkStdLog_Printf(b *testing.B) {
	const logString = "lorg"
	var buffer bytes.Buffer

	log := stdlog.New(&buffer, "", 0)
	log.SetOutput(&buffer)

	for i := 0; i < b.N; i++ {
		buffer.Reset()
		log.Printf("%v", logString)
	}
}
