package model

import tea "charm.land/bubbletea/v2"

// normalizePlayNextOverlay keeps the selected queue row and scroll window valid
// after queue mutations, including mutations made while the overlay is hidden.
func (m *Model) normalizePlayNextOverlay() {
	if m.queue == nil {
		m.playNext.cursor = 0
		m.playNext.scroll = 0
		return
	}
	count := m.queue.PlayNextLen()
	if count == 0 {
		m.playNext.cursor = 0
		m.playNext.scroll = 0
		return
	}
	m.playNext.cursor = min(max(0, m.playNext.cursor), count-1)
	visible := m.effectiveQueueVisible()
	if visible <= 0 {
		m.playNext.scroll = min(max(0, m.playNext.scroll), count-1)
		if m.playNext.cursor < m.playNext.scroll {
			m.playNext.scroll = m.playNext.cursor
		}
		return
	}
	clampScroll(&m.playNext.cursor, &m.playNext.scroll, count, visible)
}

// handlePlayNextKey processes key presses while the queue manager overlay is open.
func (m *Model) handlePlayNextKey(msg tea.KeyPressMsg) tea.Cmd {
	qLen := m.queue.PlayNextLen()

	switch msg.String() {
	case "ctrl+x":
		m.toggleExpandedView()
		m.normalizePlayNextOverlay()
	case "up", "k":
		if m.playNext.cursor > 0 {
			m.playNext.cursor--
		} else if qLen > 0 {
			m.playNext.cursor = qLen - 1
		}
		m.normalizePlayNextOverlay()

	case "down", "j":
		if m.playNext.cursor < qLen-1 {
			m.playNext.cursor++
		} else if qLen > 0 {
			m.playNext.cursor = 0
		}
		m.normalizePlayNextOverlay()
	case "shift+up":
		moved := false
		if m.playNext.cursor > 0 {
			if m.queue.MovePlayNext(m.playNext.cursor, m.playNext.cursor-1) {
				m.playNext.cursor--
				moved = true
			}
		}
		m.normalizePlayNextOverlay()
		if moved {
			return m.rearmStalePreload()
		}
	case "shift+down":
		moved := false
		if m.playNext.cursor < qLen-1 {
			if m.queue.MovePlayNext(m.playNext.cursor, m.playNext.cursor+1) {
				m.playNext.cursor++
				moved = true
			}
		}
		m.normalizePlayNextOverlay()
		if moved {
			return m.rearmStalePreload()
		}
	case "d":
		removed := false
		if qLen > 0 {
			snapshot := m.queue.Snapshot()
			m.queue.RemovePlayNextAt(m.playNext.cursor)
			m.recordQueueUndo(queueUndo{snapshot: snapshot})
			removed = true
			m.status.Show("Removed from play next (Ctrl+Z to undo)", statusTTLDefault)
		} else {
			m.status.Warning("Play next is empty.", statusTTLDefault)
		}
		m.playNext.confirmClear = false
		m.normalizePlayNextOverlay()
		if removed {
			return m.rearmStalePreload()
		}
	case "c":
		if qLen == 0 {
			m.status.Warning("Play next is empty.", statusTTLDefault)
			return nil
		}
		if qLen == 1 || m.playNext.confirmClear {
			snapshot := m.queue.Snapshot()
			m.queue.ClearPlayNext()
			m.recordQueueUndo(queueUndo{snapshot: snapshot})
			m.playNext.confirmClear = false
			m.normalizePlayNextOverlay()
			m.status.Show("Cleared play next (Ctrl+Z to undo)", statusTTLDefault)
			return m.rearmStalePreload()
		}
		m.playNext.confirmClear = true
		m.status.Showf(statusTTLDefault, "Clear %d tracks from play next? Press c again to confirm, Esc to cancel.", qLen)
	case "esc", "A", "q":
		m.playNext.confirmClear = false
		m.playNext.visible = false
	}
	return nil
}
