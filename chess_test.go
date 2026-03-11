package main

import (
	"testing"
)

func TestLayoutPrint(t *testing.T) {
	layout := Layout{
		{3, 1}: 'W',
		{2, 1}: ' ', {2, 2}: ' ',
		{1, 1}: ' ', {1, 2}: 'W', {1, 3}: ' ',
		{0, 0}: 'B', {0, 1}: ' ', {0, 2}: 'B', {0, 3}: ' ',
	}
	got := stringLayout(layout)
	want := " ♘  \n □□ \n □♘□\n♞□♞□\n"
	if got != want {
		t.Errorf("stringLayout() =\n%q\nwant\n%q", got, want)
	}
}

func TestLegalMoveActuallyMoves(t *testing.T) {
	board := newChessboard(startLayout, goalLayout)
	ok := board.movePiece(Position{0, 0}, Position{2, 1})
	if !ok {
		t.Fatal("expected movePiece to return true")
	}
	expected := Layout{
		{3, 1}: 'W',
		{2, 1}: 'B', {2, 2}: ' ',
		{1, 1}: ' ', {1, 2}: 'W', {1, 3}: ' ',
		{0, 0}: ' ', {0, 1}: ' ', {0, 2}: 'B', {0, 3}: ' ',
	}
	if !layoutsEqual(board.layout, expected) {
		t.Errorf("layout after move =\n%s\nwant\n%s", stringLayout(board.layout), stringLayout(expected))
	}
}

func TestMoveToBusyFails(t *testing.T) {
	board := newChessboard(startLayout, goalLayout)
	if board.movePiece(Position{0, 0}, Position{1, 2}) {
		t.Error("expected movePiece to return false when destination is occupied")
	}
}

func TestMoveFromEmptyFails(t *testing.T) {
	board := newChessboard(startLayout, goalLayout)
	if board.movePiece(Position{0, 1}, Position{1, 2}) {
		t.Error("expected movePiece to return false when source is empty")
	}
}

func TestMoveToOutsideFails(t *testing.T) {
	board := newChessboard(startLayout, goalLayout)
	if board.movePiece(Position{0, 0}, Position{2, -1}) {
		t.Error("expected movePiece to return false when destination is outside the board")
	}
}

func TestFindBlack(t *testing.T) {
	board := newChessboard(startLayout, goalLayout)
	blacks := board.findPieces('B')
	expected := map[Position]bool{{0, 0}: true, {0, 2}: true}
	if len(blacks) != len(expected) {
		t.Fatalf("findPieces('B') returned %d pieces, want 2", len(blacks))
	}
	for _, pos := range blacks {
		if !expected[pos] {
			t.Errorf("unexpected black piece at %v", pos)
		}
	}
}

func TestFindWhite(t *testing.T) {
	board := newChessboard(startLayout, goalLayout)
	whites := board.findPieces('W')
	expected := map[Position]bool{{1, 2}: true, {3, 1}: true}
	if len(whites) != len(expected) {
		t.Fatalf("findPieces('W') returned %d pieces, want 2", len(whites))
	}
	for _, pos := range whites {
		if !expected[pos] {
			t.Errorf("unexpected white piece at %v", pos)
		}
	}
}

func TestFindNextBoards(t *testing.T) {
	board := newChessboard(startLayout, goalLayout)
	next := board.nextBoards()

	expected1 := Layout{
		{3, 1}: 'W',
		{2, 1}: 'B', {2, 2}: ' ',
		{1, 1}: ' ', {1, 2}: 'W', {1, 3}: ' ',
		{0, 0}: ' ', {0, 1}: ' ', {0, 2}: 'B', {0, 3}: ' ',
	}
	expected2 := Layout{
		{3, 1}: 'W',
		{2, 1}: 'B', {2, 2}: ' ',
		{1, 1}: ' ', {1, 2}: 'W', {1, 3}: ' ',
		{0, 0}: 'B', {0, 1}: ' ', {0, 2}: ' ', {0, 3}: ' ',
	}

	if len(next) != 2 {
		t.Fatalf("nextBoards() returned %d boards, want 2", len(next))
	}

	found1, found2 := false, false
	for _, b := range next {
		if layoutsEqual(b.layout, expected1) {
			found1 = true
		}
		if layoutsEqual(b.layout, expected2) {
			found2 = true
		}
		if len(b.history) != 2 {
			t.Errorf("expected history length 2, got %d", len(b.history))
		}
	}
	if !found1 {
		t.Error("expected board 1 not found in nextBoards()")
	}
	if !found2 {
		t.Error("expected board 2 not found in nextBoards()")
	}
}

func TestFindNextBoardsHarder(t *testing.T) {
	start := Layout{
		{3, 1}: 'W',
		{2, 1}: ' ', {2, 2}: ' ',
		{1, 1}: ' ', {1, 2}: ' ', {1, 3}: 'B',
		{0, 0}: 'W', {0, 1}: ' ', {0, 2}: 'B', {0, 3}: ' ',
	}
	expectedLayouts := []Layout{
		{{3, 1}: ' ', {2, 1}: ' ', {2, 2}: ' ', {1, 1}: ' ', {1, 2}: 'W', {1, 3}: 'B', {0, 0}: 'W', {0, 1}: ' ', {0, 2}: 'B', {0, 3}: ' '},
		{{3, 1}: 'W', {2, 1}: 'W', {2, 2}: ' ', {1, 1}: ' ', {1, 2}: ' ', {1, 3}: 'B', {0, 0}: ' ', {0, 1}: ' ', {0, 2}: 'B', {0, 3}: ' '},
		{{3, 1}: 'W', {2, 1}: 'B', {2, 2}: ' ', {1, 1}: ' ', {1, 2}: ' ', {1, 3}: 'B', {0, 0}: 'W', {0, 1}: ' ', {0, 2}: ' ', {0, 3}: ' '},
		{{3, 1}: 'W', {2, 1}: 'B', {2, 2}: ' ', {1, 1}: ' ', {1, 2}: ' ', {1, 3}: ' ', {0, 0}: 'W', {0, 1}: ' ', {0, 2}: 'B', {0, 3}: ' '},
		{{3, 1}: 'W', {2, 1}: ' ', {2, 2}: ' ', {1, 1}: ' ', {1, 2}: ' ', {1, 3}: ' ', {0, 0}: 'W', {0, 1}: 'B', {0, 2}: 'B', {0, 3}: ' '},
		{{3, 1}: 'W', {2, 1}: ' ', {2, 2}: ' ', {1, 1}: ' ', {1, 2}: 'W', {1, 3}: 'B', {0, 0}: ' ', {0, 1}: ' ', {0, 2}: 'B', {0, 3}: ' '},
	}

	board := newChessboard(start, goalLayout)
	next := board.nextBoards()

	if len(next) != len(expectedLayouts) {
		t.Fatalf("nextBoards() returned %d boards, want %d", len(next), len(expectedLayouts))
	}

	for _, exp := range expectedLayouts {
		found := false
		for _, b := range next {
			if layoutsEqual(b.layout, exp) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected board not found:\n%s", stringLayout(exp))
		}
	}
}

func TestEasyTraverseHasTwoSolutions(t *testing.T) {
	easyStart := Layout{
		{3, 1}: 'B',
		{2, 1}: 'W', {2, 2}: ' ',
		{1, 1}: ' ', {1, 2}: 'B', {1, 3}: ' ',
		{0, 0}: ' ', {0, 1}: ' ', {0, 2}: 'W', {0, 3}: ' ',
	}
	board := newChessboard(easyStart, goalLayout)
	completed, _ := traverse(board)
	if len(completed) != 2 {
		t.Errorf("traverse() returned %d solutions, want 2", len(completed))
	}
}

func TestTraverseFindsTheTwoSolutions(t *testing.T) {
	board := newChessboard(startLayout, goalLayout)
	completed, _ := traverse(board)
	if len(completed) != 2 {
		t.Errorf("traverse() returned %d solutions, want 2", len(completed))
	}
}
