// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package ui

import (
	"fmt"
	"log/slog"

	"github.com/derailed/k9s/internal/model"
	"github.com/derailed/k9s/internal/slogs"
	"github.com/derailed/tview"
)

// Pages represents a stack of view pages.
type Pages struct {
	*tview.Pages
	*model.Stack
}

// NewPages return a new view.
func NewPages() *Pages {
	p := Pages{
		Pages: tview.NewPages(),
		Stack: model.NewStack(),
	}
	p.AddListener(&p)

	return &p
}

// IsTopDialog checks if front page is a dialog.
func (p *Pages) IsTopDialog() bool {
	_, pa := p.GetFrontPage()
	switch pa.(type) {
	case *tview.ModalForm, *ModalList:
		return true
	default:
		return false
	}
}

// Show displays a given page.
func (p *Pages) Show(c model.Component) {
	p.SwitchToPage(componentID(c))
}

// Current returns the current component.
func (p *Pages) Current() model.Component {
	c := p.CurrentPage()
	if c == nil {
		return nil
	}

	return c.Item.(model.Component)
}

// AddAndShow adds a new page and bring it to front.
func (p *Pages) addAndShow(c model.Component) {
	p.add(c)
	p.Show(c)
}

// Add adds a new page.
func (p *Pages) add(c model.Component) {
	p.AddPage(componentID(c), c, true, true)
}

// Delete removes a page.
func (p *Pages) delete(c model.Component) {
	p.RemovePage(componentID(c))
}

// Dump for debug.
func (p *Pages) Dump() {
	slog.Debug("Dumping Pages", slogs.Page, p)
	for i, c := range p.Peek() {
		slog.Debug(fmt.Sprintf("%d -- %s -- %#v", i, componentID(c), p.GetPrimitive(componentID(c))))
	}
}

// Stack Protocol...

// StackPushed notifies a new component was pushed.
func (p *Pages) StackPushed(c model.Component) {
	p.addAndShow(c)
}

// StackPopped notifies a component was removed.
func (p *Pages) StackPopped(o, _ model.Component) {
	p.delete(o)
}

// StackTop notifies a new component is at the top of the stack.
func (p *Pages) StackTop(top model.Component) {
	if top == nil {
		return
	}
	p.Show(top)
}

// Helpers...

func componentID(c model.Component) string {
	if c.Name() == "" {
		slog.Error("Component has no name", slogs.Component, fmt.Sprintf("%T", c))
	}
	return fmt.Sprintf("%s-%p", c.Name(), c)
}
