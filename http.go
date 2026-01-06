package httpmatter

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

// trip represent a request and possible responses
// Each trip can only be used once, thus same request needs to be
// added multiple times to handle multiple responses.
type trip struct {
	req       *RequestMatter
	resps     []*ResponseMatter
	responder responder
}

type HTTP struct {
	t          testing.TB
	namespaces []string
	trip       *trip
	trips      []*trip
}

func NewHTTP(t testing.TB, namespaces ...string) *HTTP {
	return &HTTP{
		t:          t,
		namespaces: namespaces,
		trip:       nil,
		trips:      []*trip{},
	}
}

func (h *HTTP) Init() {
	groups := make(map[string][]*trip)

	for _, trip := range h.trips {
		err := makeMatter(trip.req, WithTB(h.t))
		if err != nil {
			h.t.Fatalf("error creating matter for %s: %v", trip.req.Name, err)
		}
		for _, resp := range trip.resps {
			err := makeMatter(resp, WithTB(h.t))
			if err != nil {
				h.t.Fatalf("error creating matter for %s: %v", resp.Name, err)
			}
		}

		key := h.toKey(trip.req)
		groups[key] = append(groups[key], trip)
	}

	for key, group := range groups {
		h.t.Logf("Registering responder for %s with %d trips", key, len(group))
		httpmock.RegisterResponder(
			// Because group always have same Method and URL,
			// we can use the first trip to get the method and url
			group[0].req.Method,
			group[0].req.URL.String(),
			func(r *http.Request) (*http.Response, error) {
				if len(group) == 0 {
					h.t.Fatalf("no more trips for %s", key)
				}
				chosen := group[0].responder(r, group[0].req, group[0].resps)
				// Now when the first call is responded, we need to remove
				// the trip from the group
				before := len(group)
				group = group[1:]
				after := len(group)
				if before == after {
					h.t.Fatalf("group did not change length")
				}

				h.t.Logf("group length changed from %d to %d", before, after)
				return chosen.Response, nil
			},
		)
	}

	httpmock.Activate()
}
func (h *HTTP) Destroy() {
	// Verify that every registered responder was actually invoked at least once.
	// This helps ensure test coverage for multi-step flows (e.g., paginated calls).
	callCounts := httpmock.GetCallCountInfo()
	httpmock.DeactivateAndReset()

	var missing []string
	for _, trip := range h.trips {
		key := h.toKey(trip.req)
		if callCounts[key] == 0 {
			missing = append(missing, trip.req.Name)
		}
	}
	if len(missing) > 0 {
		h.t.Fatalf("uninvoked HTTP mocks: %v", missing)
	}
}

func (h *HTTP) Add(reqname string, respnames ...string) *HTTP {
	h.t.Logf("Adding request %s with %d responses", reqname, len(respnames))
	if h.trip != nil {
		h.t.Fatalf(
			"request set already pending:%v, use Respond() to continue",
			h.trip.req.Name)
	}
	req := h.newRequest(reqname)
	resps := make([]*ResponseMatter, len(respnames))
	for i, respname := range respnames {
		resps[i] = h.newResponse(respname)
	}
	reqSet := &trip{
		req:       req,
		resps:     resps,
		responder: nil,
	}
	h.trip = reqSet
	return h
}

func (h *HTTP) Respond(fn responder) *HTTP {
	if h.trip == nil {
		h.t.Fatalf("no requests found")
	}
	if fn == nil {
		fn = DefaultResponder
	}
	h.trip.responder = fn

	h.trips = append(h.trips, h.trip)
	h.trip = nil
	return h
}

// newRequest check for the request in each namespace
// It gives priority to the first namespace that has the request
func (h *HTTP) newRequest(name string) *RequestMatter {
	for _, namespace := range h.namespaces {
		req := NewRequestMatter(namespace, name)
		if err := req.Validate(); err == nil {
			return req
		} else if ErrReadingFile().Is(err) {
			continue
		}
	}
	h.t.Fatalf("no request found for %s in %v", name, h.namespaces)
	return nil
}

// newResponse check for the response in each namespace
// It gives priority to the first namespace that has the response
func (h *HTTP) newResponse(name string) *ResponseMatter {
	for _, namespace := range h.namespaces {
		resp := NewResponseMatter(namespace, name)
		if err := resp.Validate(); err == nil {
			return resp
		} else if ErrReadingFile().Is(err) {
			continue
		} else {
			h.t.Fatalf("error creating matter: %v", err)
		}
	}
	h.t.Fatalf("no response found for %s in %v", name, h.namespaces)
	return nil
}

func (h *HTTP) toKey(req *RequestMatter) string {
	if req == nil || req.Request == nil || req.Method == "" || req.URL == nil {
		h.t.Fatalf("request has no method or url")
	}
	return fmt.Sprintf("%s %s", req.Method, req.URL.String())
}
