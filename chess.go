package main

import (
	"fmt"
	"strings"
)

// Position represents a (row, col) coordinate on the board
type Position struct {
	row, col int
}

// Layout is the board state: a map from Position to piece ('W', 'B', or ' ')
type Layout map[Position]rune

// pieceMappings maps piece chars to display symbols
var pieceMapping = map[rune]string{
	'W': "♘",
	'B': "♞",
	' ': "□",
}

// knightMoves lists all possible knight move offsets
var knightMoves = []Position{
	{2, -1}, {2, 1}, {1, -2}, {1, 2},
	{-2, -1}, {-2, 1}, {-1, 2}, {-1, -2},
}

var startLayout = Layout{
	{0, 0}: 'B', {0, 1}: ' ', {0, 2}: 'B', {0, 3}: ' ',
	{1, 1}: ' ', {1, 2}: 'W', {1, 3}: ' ',
	{2, 1}: ' ', {2, 2}: ' ',
	{3, 1}: 'W',
}

var goalLayout = Layout{
	{0, 0}: 'W', {0, 1}: ' ', {0, 2}: 'W', {0, 3}: ' ',
	{1, 1}: ' ', {1, 2}: 'B', {1, 3}: ' ',
	{2, 1}: ' ', {2, 2}: ' ',
	{3, 1}: 'B',
}

// stringLayout returns an ASCII representation of a layout, used for display and as a map key
func stringLayout(layout Layout) string {
	var sb strings.Builder
	for row := 3; row >= 0; row-- {
		for col := 0; col < 4; col++ {
			pos := Position{row, col}
			if piece, ok := layout[pos]; ok {
				sb.WriteString(pieceMapping[piece])
			} else {
				sb.WriteString(" ")
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// copyLayout performs a deep copy of a layout
func copyLayout(src Layout) Layout {
	dst := make(Layout, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// layoutsEqual returns true if two layouts are identical
func layoutsEqual(a, b Layout) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

// Chessboard holds the current layout, goal, and move history
type Chessboard struct {
	layout     Layout
	goalLayout Layout
	history    []Layout
}

// newChessboard creates a new board with the given start and goal layouts
func newChessboard(start, goal Layout) *Chessboard {
	layout := copyLayout(start)
	return &Chessboard{
		layout:     layout,
		goalLayout: goal,
		history:    []Layout{copyLayout(layout)},
	}
}

// finished returns true when the board has reached the goal layout
func (cb *Chessboard) finished() bool {
	return layoutsEqual(cb.layout, cb.goalLayout)
}

// findPieces returns all positions containing the given piece
func (cb *Chessboard) findPieces(piece rune) []Position {
	var positions []Position
	for pos, p := range cb.layout {
		if p == piece {
			positions = append(positions, pos)
		}
	}
	return positions
}

// printHistory returns a string of the full board history
func (cb *Chessboard) printHistory() string {
	var sb strings.Builder
	for _, layout := range cb.history {
		sb.WriteString("\n")
		sb.WriteString(stringLayout(layout))
	}
	return sb.String()
}

// movePiece moves a piece from start to end, recording the move in history.
// Returns true on success, false on failure.
func (cb *Chessboard) movePiece(start, end Position) bool {
	if _, ok := cb.layout[start]; !ok {
		return false
	}
	if _, ok := cb.layout[end]; !ok {
		return false
	}
	if cb.layout[end] != ' ' {
		return false
	}
	if cb.layout[start] == ' ' {
		return false
	}
	piece := cb.layout[start]
	cb.layout[start] = ' '
	cb.layout[end] = piece
	cb.history = append(cb.history, copyLayout(cb.layout))
	return true
}

// deepCopy returns a full deep copy of the chessboard
func (cb *Chessboard) deepCopy() *Chessboard {
	historyCopy := make([]Layout, len(cb.history))
	for i, h := range cb.history {
		historyCopy[i] = copyLayout(h)
	}
	return &Chessboard{
		layout:     copyLayout(cb.layout),
		goalLayout: copyLayout(cb.goalLayout),
		history:    historyCopy,
	}
}

// nextBoards generates all valid next board states from the current one
func (cb *Chessboard) nextBoards() []*Chessboard {
	var boards []*Chessboard
	var allPieces []Position
	allPieces = append(allPieces, cb.findPieces('B')...)
	allPieces = append(allPieces, cb.findPieces('W')...)

	for _, startPos := range allPieces {
		for _, move := range knightMoves {
			endPos := Position{startPos.row + move.row, startPos.col + move.col}
			newBoard := cb.deepCopy()
			if newBoard.movePiece(startPos, endPos) {
				boards = append(boards, newBoard)
			}
		}
	}
	return boards
}

// traverse performs BFS over all possible board states, returning boards that
// reach the goal and all visited layouts along the way
func traverse(board *Chessboard) ([]*Chessboard, []Layout) {
	var finishedBoards []*Chessboard
	boards := []*Chessboard{board}
	visitedLayouts := make(map[string]Layout)

	for len(boards) > 0 {
		// Record all current layouts as visited
		for _, b := range boards {
			key := stringLayout(b.layout)
			visitedLayouts[key] = copyLayout(b.layout)
		}

		// Generate next boards, excluding already-visited layouts
		uniqueNext := make(map[string]*Chessboard)
		for _, outerBoard := range boards {
			for _, innerBoard := range outerBoard.nextBoards() {
				key := stringLayout(innerBoard.layout)
				if _, visited := visitedLayouts[key]; visited {
					continue
				}
				if innerBoard.finished() {
					finishedBoards = append(finishedBoards, innerBoard)
				} else {
					// Keep only one board per layout (discard duplicate histories)
					uniqueNext[key] = innerBoard
				}
			}
		}

		boards = make([]*Chessboard, 0, len(uniqueNext))
		for _, b := range uniqueNext {
			boards = append(boards, b)
		}
	}

	visited := make([]Layout, 0, len(visitedLayouts))
	for _, v := range visitedLayouts {
		visited = append(visited, v)
	}
	return finishedBoards, visited
}

func main() {
	board := newChessboard(startLayout, goalLayout)
	completedBoards, visitedLayouts := traverse(board)
	for _, b := range completedBoards {
		fmt.Println("*********************")
		fmt.Printf("In %d moves, visiting %d layouts (excluding the final one):\n",
			len(b.history), len(visitedLayouts))
		fmt.Println(b.printHistory())
		fmt.Println()
	}
}
