package editor

import (
	"github.com/piodrs/sea/internal/buffer"

	"github.com/charmbracelet/bubbles/viewport"
)

type Direction int

const (
	DirectionHorizontal Direction = iota
	DirectionVertical
)

type pane struct {
	buffer    *buffer.Buffer
	cursor    int
	left      int
	viewport  viewport.Model
	children  [2]*pane
	direction Direction
	width     int
	height    int
}

func (p *pane) panes() []*pane {
	if p.buffer != nil {
		return []*pane{p}
	}

	return append(p.children[0].panes(), p.children[1].panes()...)
}

func (p *pane) contains(target *pane) bool {
	if p.buffer != nil {
		return p == target
	}

	return p.children[0].contains(target) || p.children[1].contains(target)
}

func (p *pane) remove(target *pane) *pane {
	if p == target {
		return nil
	}

	if p.buffer != nil {
		return p
	}

	for index, child := range p.children {
		p.children[index] = child.remove(target)

		if p.children[index] == nil {
			return p.children[1-index]
		}
	}

	return p
}

func (p *pane) resize(width, height int, focused *pane) {
	p.width = width
	p.height = height

	if p.buffer != nil {
		p.viewport.Width = max(1, width)
		p.viewport.Height = max(1, height-1)
		p.refresh(p == focused)

		return
	}

	first, second := p.children[0], p.children[1]
	tooNarrow := p.direction == DirectionVertical && width < 3
	tooShort := p.direction == DirectionHorizontal && height < 5

	if width == 0 || height == 0 || tooNarrow || tooShort {
		visible, hidden := first, second

		if second.contains(focused) {
			visible, hidden = second, first
		}

		visible.resize(width, height, focused)
		hidden.resize(0, 0, focused)

		return
	}

	if p.direction == DirectionVertical {
		available := width - 1
		first.resize(available/2, height, focused)
		second.resize(available-available/2, height, focused)

		return
	}

	available := height - 1
	first.resize(width, available/2, focused)
	second.resize(width, available-available/2, focused)
}
