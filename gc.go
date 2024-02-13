package promdev_sd

import (
	"context"
	"log"
	"time"
)

func (s *Server) RunGCThread(ctx context.Context) {
	go s.internalRunGCThread(ctx, s.config.GCPeriod)
}

func (s *Server) internalRunGCThread(ctx context.Context, interval time.Duration) {
outer:
	for {
		select {
		case <-ctx.Done():
			break outer
		case <-time.Tick(interval):
			s.cleanup()
		}
	}
}

func (s *Server) cleanup() {
	targetTime := time.Now().Add(-s.config.TimeoutPeriod)
	s.Lock()
	defer s.Unlock()
	useableLabels := make(map[string]bool)
	for ns, ts := range s.namespaceEntries {
		var cleaned []*Target
		for _, t := range ts {
			if t.TouchedAt.After(targetTime) {
				cleaned = append(cleaned, t)
				useableLabels[t.LabelSetHash] = true
			} else {
				log.Printf("Cleaning up target %s/%s token %s", t.Namespace, t.Target, t.Token)
			}
		}
		s.namespaceEntries[ns] = cleaned
	}
	for k := range s.labelSets {
		if !useableLabels[k] {
			delete(s.labelSets, k)
		}
	}
}
