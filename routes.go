package promdev_sd

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	ns := chi.URLParam(r, "namespace")
	if ns == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var ts TargetSet
	err := json.NewDecoder(r.Body).Decode(&ts)
	if err != nil {
		log.Println("Error decoding body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token := NewToken()

	labels := ts.Labels
	labelhash := labels.Hash()
	now := time.Now()

	s.Lock()
	defer s.Unlock()

	s.labelSets[labelhash] = labels

	for _, t := range ts.Targets {
		target := &Target{
			Target:       t,
			LabelSetHash: labelhash,
			TouchedAt:    now,
			Token:        token,
			Namespace:    ns,
		}
		s.namespaceEntries[ns] = append(s.namespaceEntries[ns], target)
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}

func (s *Server) heartbeat(w http.ResponseWriter, r *http.Request) {
	ns := chi.URLParam(r, "namespace")
	if ns == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token := chi.URLParam(r, "token")
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s.Lock()
	defer s.Unlock()

	updateTime := time.Now()
	cur := s.namespaceEntries[ns]
	for _, c := range cur {
		if c.Token == token {
			c.TouchedAt = updateTime
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) delete(w http.ResponseWriter, r *http.Request) {
	ns := chi.URLParam(r, "namespace")
	if ns == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token := chi.URLParam(r, "token")
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s.Lock()
	defer s.Unlock()

	cur := s.namespaceEntries[ns]
	var new []*Target
	for _, c := range cur {
		if c.Token != token {
			new = append(new, c)
		}
	}
	s.namespaceEntries[ns] = new
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
	if len(tl) == 0 {
		// Hide nils
		tl = make(TargetList, 0)
	}
	w.Header().Add("Content-type", "application/json")
	err := json.NewEncoder(w).Encode(tl)
	if err != nil {
		log.Println("Couldn't encode targets:", err)
	}
}

func (s *Server) convertEntries(entries []*Target) TargetList {
	group := make(map[string][]*Target)
	for _, e := range entries {
		v, ok := group[e.Token]
		if ok {
			group[e.Token] = append(v, e)
		} else {
			group[e.Token] = []*Target{e}
		}
	}
	var out TargetList
	for _, ts := range group {
		var targets []string
		for _, t := range ts {
			targets = append(targets, t.Target)
		}
		v := TargetSet{
			Targets: targets,
			Labels:  s.labelSets[ts[0].LabelSetHash].Clone(),
		}
		out = append(out, v)
	}
	return out
}
