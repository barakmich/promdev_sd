package promdev_sd

import (
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"sync"
	"time"
)

type TargetList []TargetSet

type LabelSet map[string]string

var hashPool = &sync.Pool{
	New: func() any {
		return sha1.New()
	},
}

func (l LabelSet) Hash() string {
	h := hashPool.Get().(hash.Hash)
	for k, v := range l {
		h.Write([]byte(k))
		h.Write([]byte{':'})
		h.Write([]byte(v))
		h.Write([]byte{'\n'})
	}
	sum := h.Sum(nil)
	h.Reset()
	hashPool.Put(h)
	return hex.EncodeToString(sum)
}

func (l LabelSet) Clone() LabelSet {
	out := make(LabelSet)
	for k, v := range l {
		out[k] = v
	}
	return out
}

type TargetSet struct {
	Targets []string `json:"targets"`
	Labels  LabelSet `json:"labels"`
}

type Target struct {
	Target       string
	LabelSetHash string
	TouchedAt    time.Time
	Token        string
	Namespace    string
}
