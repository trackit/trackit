package routes

// Documentation decorates a handler to document it. Summary, Description and
// Tags will be set on the documentation if not zero.
type Documentation HandlerDocumentationBody

func (d Documentation) Decorate(h Handler) Handler {
	n := h.Documentation
	if d.Summary != "" {
		n.Summary = d.Summary
	}
	if d.Description != "" {
		n.Description = d.Description
	}
	for k, v := range d.Tags {
		if n.Tags == nil {
			n.Tags = make(Tags)
		}
		t := n.Tags[k]
		n.Tags[k] = append(t, v...)
	}
	h.Documentation = n
	return h
}
