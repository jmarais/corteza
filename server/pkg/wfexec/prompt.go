package wfexec

import (
	"fmt"
	"io"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/expr"
)

type (
	prompted struct {
		// when !nul, assuming waiting for input we're waiting for input
		payload *expr.Vars

		// user to be prompted
		ownerId uint64

		// state to be resumed
		state *State

		sent bool

		// prompt reference; something client can use
		// for orientation, what kind of prompt is expected
		ref string
	}

	PendingPrompt struct {
		Ref       string     `json:"ref"`
		SessionID uint64     `json:"sessionID,string"`
		CreatedAt time.Time  `json:"createdAt"`
		StateID   uint64     `json:"stateID,string"`
		Payload   *expr.Vars `json:"payload"`
		OwnerId   uint64     `json:"-"`

		Original *prompted `json:"-"`
	}

	ResumedPrompt struct {
		StateID uint64 `json:"stateID,string"`
		OwnerId uint64 `json:"-"`
	}
)

func Prompt(ownerId uint64, ref string, payload *expr.Vars) *prompted {
	return &prompted{payload: payload, ref: ref, ownerId: ownerId}
}

func (p *prompted) PPrint(w io.Writer) {
	fmt.Fprintf(w, "<td>%d</td>", p.ownerId)
	fmt.Fprintf(w, "<td>%s</td>", p.ref)
	fmt.Fprintf(w, "<td>%s</td>", p.state.created)
	fmt.Fprintf(w, "<td>%d</td>", p.state.stateId)
	if p.state.parent != nil {
		fmt.Fprintf(w, "<td>%d</td>", p.state.parent.ID())
	} else {
		fmt.Fprintf(w, "<td></td>")
	}
	fmt.Fprintf(w, "<td>%d</td>", p.state.step.ID())
	if len(p.state.next) > 0 {
		fmt.Fprintf(w, "<td>%d</td>", p.state.next[0].ID())
	} else {
		fmt.Fprintf(w, "<td></td>")
	}
	fmt.Fprintf(w, "<td>%s</td>", p.state.action)

}
func (p *prompted) toPending() *PendingPrompt {
	return &PendingPrompt{
		Ref:       p.ref,
		CreatedAt: p.state.created,
		StateID:   p.state.stateId,
		Payload:   p.payload,
		OwnerId:   p.ownerId,
		Original:  p,
	}
}

func (p *prompted) toResumed() *ResumedPrompt {
	return &ResumedPrompt{
		StateID: p.state.stateId,
		OwnerId: p.ownerId,
	}
}

func (p *prompted) MarkSent() {
	p.sent = true
}
