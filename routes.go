package promdev_sd

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) heartbeat(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) discovery(w http.ResponseWriter, r *http.Request) {
	ns := chi.URLParam(r, "namespace")
	if ns == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s.Lock()
	defer s.Unlock()
	entries := s.namespaceEntries[ns]
	// Could be empty/non-existant. Prometheus documentation suggests
	// that returning 200 is required, and empty list is acceptable
	// so the fallthrough will work just fine.
	tl := s.convertEntries(entries)
	w.Header().Add("Content-type", "application/json")
	err := json.NewEncoder(w).Encode(tl)
	if err != nil {
		log.Println("Couldn't encode targets:", err)
	}
}

func (s *Server) convertEntries(entries []*Target) TargetList {
	labelset := make(map[string][]*Target)
	for _, e := range entries {
		v, ok := labelset[e.LabelSetHash]
		if ok {
			labelset[e.LabelSetHash] = append(v, e)
		} else {
			labelset[e.LabelSetHash] = []*Target{e}
		}
	}
	var out TargetList
	for l, ts := range labelset {
		var targets []string
		for _, t := range ts {
			targets = append(targets, t.Target)
		}
		v := TargetSet{
			Targets: targets,
			Labels:  s.labelSets[l].Clone(),
		}
		out = append(out, v)
	}
	return out
}
