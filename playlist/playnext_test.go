package playlist

import (
	"slices"
	"testing"
)

func TestQueueAndDequeue(t *testing.T) {
	p := makePlaylist(5, false)

	p.AddPlayNext(2)
	p.AddPlayNext(4)

	if p.PlayNextLen() != 2 {
		t.Fatalf("PlayNextLen() = %d, want 2", p.PlayNextLen())
	}

	// RemovePlayNext first item
	if !p.RemovePlayNext(2) {
		t.Fatal("RemovePlayNext(2) = false, want true")
	}
	if p.PlayNextLen() != 1 {
		t.Fatalf("PlayNextLen() after dequeue = %d, want 1", p.PlayNextLen())
	}

	// RemovePlayNext non-existent
	if p.RemovePlayNext(2) {
		t.Fatal("RemovePlayNext(2) again = true, want false")
	}
}

func TestQueuePosition(t *testing.T) {
	p := makePlaylist(5, false)

	p.AddPlayNext(1)
	p.AddPlayNext(3)
	p.AddPlayNext(4)

	if pos := p.PlayNextPosition(1); pos != 1 {
		t.Errorf("PlayNextPosition(1) = %d, want 1", pos)
	}
	if pos := p.PlayNextPosition(3); pos != 2 {
		t.Errorf("PlayNextPosition(3) = %d, want 2", pos)
	}
	if pos := p.PlayNextPosition(4); pos != 3 {
		t.Errorf("PlayNextPosition(4) = %d, want 3", pos)
	}
	if pos := p.PlayNextPosition(0); pos != 0 {
		t.Errorf("PlayNextPosition(0) = %d, want 0 (not queued)", pos)
	}
}

func TestQueueTracks(t *testing.T) {
	p := makePlaylist(5, false)

	p.AddPlayNext(0) // A
	p.AddPlayNext(2) // C
	p.AddPlayNext(4) // E

	qt := p.PlayNextTracks()
	if len(qt) != 3 {
		t.Fatalf("PlayNextTracks len = %d, want 3", len(qt))
	}
	if qt[0].Title != "A" || qt[1].Title != "C" || qt[2].Title != "E" {
		t.Errorf("PlayNextTracks = [%s, %s, %s], want [A, C, E]",
			qt[0].Title, qt[1].Title, qt[2].Title)
	}
}

func TestTracksAndQueue(t *testing.T) {
	p := makePlaylist(3, false)
	p.AddPlayNext(2)
	p.AddPlayNext(0)

	tracks, queue := p.TracksAndPlayNext()
	if len(tracks) != 3 || tracks[0].Title != "A" || tracks[2].Title != "C" {
		t.Fatalf("tracks = %+v, want A, B, C", tracks)
	}
	if !slices.Equal(queue, []int{2, 0}) {
		t.Fatalf("queue = %v, want [2 0]", queue)
	}
	tracks[0].Title = "Changed"
	queue[0] = 1
	if tracks, queue := p.TracksAndPlayNext(); tracks[0].Title != "A" || queue[0] != 2 {
		t.Fatalf("after a change to the returned slices: tracks[0] = %q, queue[0] = %d; want A, 2", tracks[0].Title, queue[0])
	}
}

func TestQueueWindow(t *testing.T) {
	p := makePlaylist(5, false)
	for i := range 5 {
		p.AddPlayNext(i)
	}

	tests := []struct {
		name       string
		start      int
		limit      int
		wantTitles string
	}{
		{name: "middle", start: 1, limit: 2, wantTitles: "BC"},
		{name: "negative start", start: -2, limit: 2, wantTitles: "AB"},
		{name: "limit past end", start: 3, limit: 10, wantTitles: "DE"},
		{name: "start at end", start: 5, limit: 1},
		{name: "zero limit", start: 0, limit: 0},
		{name: "negative limit", start: 0, limit: -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracks := p.PlayNextWindow(tt.start, tt.limit)
			got := ""
			for _, track := range tracks {
				got += track.Title
			}
			if got != tt.wantTitles {
				t.Fatalf("PlayNextWindow(%d, %d) titles = %q, want %q", tt.start, tt.limit, got, tt.wantTitles)
			}
		})
	}

	window := p.PlayNextWindow(1, 1)
	window[0].Title = "changed"
	if got := p.PlayNextWindow(1, 1)[0].Title; got != "B" {
		t.Fatalf("mutating PlayNextWindow result changed playlist title to %q", got)
	}
}

func TestQueueOutputsOwnProviderMeta(t *testing.T) {
	p := New()
	p.Replace([]Track{{ProviderMeta: map[string]string{"provider.id": "original"}}})
	p.AddPlayNext(0)

	window := p.PlayNextWindow(0, 1)
	window[0].ProviderMeta["provider.id"] = "window"
	queued := p.PlayNextTracks()
	queued[0].ProviderMeta["provider.id"] = "queue"

	if got := p.Tracks()[0].ProviderMeta["provider.id"]; got != "original" {
		t.Fatalf("metadata after queue output mutation = %q, want original", got)
	}
}

func TestQueuePositionsStayCoherentAcrossMutations(t *testing.T) {
	p := makePlaylist(5, false)
	p.AddPlayNext(1)
	p.AddPlayNext(3)
	p.AddPlayNext(1)
	p.AddPlayNext(4)
	assertQueuePositions(t, p, map[int]int{1: 1, 3: 2, 4: 4})

	if !p.MovePlayNext(3, 0) {
		t.Fatal("MovePlayNext(3, 0) = false, want true")
	}
	assertQueuePositions(t, p, map[int]int{1: 3, 3: 2, 4: 1})

	p.RemovePlayNextAt(1)
	assertQueuePositions(t, p, map[int]int{1: 2, 4: 1})

	if !p.RemovePlayNext(1) {
		t.Fatal("RemovePlayNext(1) = false, want true")
	}
	assertQueuePositions(t, p, map[int]int{1: 2, 4: 1})

	if !p.Move(1, 2) {
		t.Fatal("Move(1, 2) = false, want true")
	}
	assertQueuePositions(t, p, map[int]int{2: 2, 4: 1})
	snapshot := p.Snapshot()

	if !p.Remove(0) {
		t.Fatal("Remove(0) = false, want true")
	}
	assertQueuePositions(t, p, map[int]int{1: 2, 3: 1})

	p.ClearPlayNext()
	assertQueuePositions(t, p, nil)
	p.Restore(snapshot)
	assertQueuePositions(t, p, map[int]int{2: 2, 4: 1})

	if track, ok := p.Next(); !ok || track.Title != "E" {
		t.Fatalf("Next() = (%q, %t), want (E, true)", track.Title, ok)
	}
	assertQueuePositions(t, p, map[int]int{2: 1})

	p.Replace([]Track{{Title: "replacement"}})
	assertQueuePositions(t, p, nil)
}

func assertQueuePositions(t *testing.T, p *Playlist, want map[int]int) {
	t.Helper()
	for idx := range p.Len() {
		if got := p.PlayNextPosition(idx); got != want[idx] {
			t.Fatalf("PlayNextPosition(%d) = %d, want %d", idx, got, want[idx])
		}
	}
}

func TestRestoreSnapshotRestoresTracksAndQueue(t *testing.T) {
	p := makePlaylist(3, false)
	p.AddPlayNext(0)
	p.AddPlayNext(2)
	snapshot := p.Snapshot()

	p.Remove(0)
	p.ClearPlayNext()
	p.Restore(snapshot)

	tracks := p.Tracks()
	if len(tracks) != 3 || tracks[0].Title != "A" {
		t.Fatalf("tracks after restore = %#v, want original tracks", tracks)
	}
	queued := p.PlayNextTracks()
	if len(queued) != 2 || queued[0].Title != "A" || queued[1].Title != "C" {
		t.Fatalf("queue after restore = %#v, want [A C]", queued)
	}
}

func TestClearQueue(t *testing.T) {
	p := makePlaylist(3, false)
	p.AddPlayNext(0)
	p.AddPlayNext(1)

	p.ClearPlayNext()

	if p.PlayNextLen() != 0 {
		t.Fatalf("PlayNextLen() after clear = %d, want 0", p.PlayNextLen())
	}
}

func TestRemoveQueueAt(t *testing.T) {
	p := makePlaylist(5, false)
	p.AddPlayNext(0) // A
	p.AddPlayNext(2) // C
	p.AddPlayNext(4) // E

	// Remove middle entry (C)
	p.RemovePlayNextAt(1)

	qt := p.PlayNextTracks()
	if len(qt) != 2 {
		t.Fatalf("PlayNextTracks after remove = %d, want 2", len(qt))
	}
	if qt[0].Title != "A" || qt[1].Title != "E" {
		t.Errorf("PlayNextTracks = [%s, %s], want [A, E]", qt[0].Title, qt[1].Title)
	}
}

func TestRemoveQueueAtBounds(t *testing.T) {
	p := makePlaylist(3, false)
	p.AddPlayNext(0)

	// Out of bounds should be no-op
	p.RemovePlayNextAt(-1)
	p.RemovePlayNextAt(5)

	if p.PlayNextLen() != 1 {
		t.Fatalf("PlayNextLen() = %d, want 1", p.PlayNextLen())
	}
}

func TestQueueBoundsCheck(t *testing.T) {
	p := makePlaylist(3, false)

	// Queue out-of-bounds track indices
	p.AddPlayNext(-1)
	p.AddPlayNext(5)

	if p.PlayNextLen() != 0 {
		t.Fatalf("PlayNextLen() = %d, want 0 (invalid indices ignored)", p.PlayNextLen())
	}
}
