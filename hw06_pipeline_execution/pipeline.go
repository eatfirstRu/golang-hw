package hw06pipelineexecution

import "sync"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if done != nil {
		_, ok := <-done
		if !ok {
			cls := make(Bi)
			defer close(cls)
			return cls
		}
	}

	out := in
	var wg sync.WaitGroup
	for _, stage := range stages {
		wg.Add(1)
		go func() {
			out = stage(out)
			wg.Done()
		}()
		wg.Wait()
	}

	return out
}
