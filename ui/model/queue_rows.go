package model

// The queue view lists tracks in play order. While shuffle is on, a view
// row and a track index can differ. m.queueCursor holds a track index, so each
// cursor action acts on the track under the cursor. m.queueScroll holds a view
// row. The helpers below are the one place that converts between the two.

// shuffleMoveWarning tells the user why Shift+Up and Shift+Down do not move
// a track while shuffle is on.
const shuffleMoveWarning = "Turn off shuffle to move tracks"
const shuffleReverseWarning = "Turn off shuffle to reverse the queue"
const pagingReverseWarning = "Wait for the queue to finish loading"

// shuffleHeadersWarning tells the user why Ctrl+H does not show the album
// headers while shuffle is on.
const shuffleHeadersWarning = "Turn off shuffle to show album headers"

// queueShowsAlbumHeaders reports whether the queue view draws album headers.
// In the shuffled play order, the tracks of one album are not next to each
// other, so the headers would split the album. The headers come back when
// shuffle is off.
func (m Model) queueShowsAlbumHeaders() bool {
	return m.showAlbumHeaders && !m.queue.Shuffled()
}

// queueCursorRow returns the view row of the playlist cursor.
func (m Model) queueCursorRow() int {
	if row := m.queue.OrderPosition(m.queueCursor); row >= 0 {
		return row
	}
	return m.queueCursor
}

// setQueueCursorRow puts the queue cursor on the track at the view row and
// scrolls the view to show the cursor.
func (m *Model) setQueueCursorRow(row int) {
	if indices, _ := m.queue.OrderWindow(row, 1); len(indices) == 1 {
		m.queueCursor = indices[0]
	} else {
		m.queueCursor = row
	}
	m.adjustScroll()
}

// keepQueueCursorRow runs change, then puts the queue cursor back on its
// view row. While shuffle is on, a change to the tracks can move the track
// under the cursor to a different row.
func (m *Model) keepQueueCursorRow(change func()) {
	row := m.queueCursorRow()
	change()
	if n := m.queue.Len(); n > 0 {
		m.setQueueCursorRow(min(row, n-1))
	}
}
