package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func makeChan(out In, done In) Out {
	midChan := make(Bi)

	go func(inChnl In, doneChnl In, out Bi) {
		defer func() {
			if out != nil {
				close(out)
			}
		}()

		for {
			select {
			case <-doneChnl:
				close(out)
				out = nil
				for {
					_, ok := <-inChnl
					if !ok {
						break
					}
				}
				return
			case d, ok := <-inChnl:
				if !ok {
					close(out)
					out = nil
					return
				}
				out <- d
			}
		}
	}(out, done, midChan)

	return midChan
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	midChan := make(Bi)

	go func(inChnl In, doneChnl In, out Bi) {
		defer close(out)

		for {
			select {
			case <-doneChnl:
				return
			case d, ok := <-inChnl:
				if !ok {
					return
				}
				out <- d
			}
		}
	}(in, done, midChan)

	currIn := In(midChan)

	for _, stage := range stages {
		out := stage(currIn)
		currIn = makeChan(out, done)
	}

	return currIn
}
