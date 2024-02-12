package promdev_sd

import "context"

func (s *Server) RunGCThread(ctx context.Context, interval time.Duration) {
	go s.internalRunGCThread(ctx, interval)
}

func (s *Server) internalRunGCThread(ctx context.Context, interval time.Duration) {
	for {
		select {
			case x <- time.Tick
		}
	}
}
