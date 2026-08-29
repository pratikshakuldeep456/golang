package chess

// ============================================================
// Core types
// ============================================================

type Color int

const (
	White Color = iota
	Black
)

type Position struct {
	Row, Col int
}

func (p Position) InBounds() bool {
	return p.Row >= 0 && p.Row < 8 && p.Col >= 0 && p.Col < 8
}

type Piece interface {
	Color() Color
	Symbol() string
	// ValidMoves returns every square this piece COULD move to,
	// ignoring whether it would leave the player's own king in check
	// (that check happens one layer up, in Game.MakeMove).
	ValidMoves(b *Board, from Position) []Position
}

// ============================================================
// Board — owns piece positions, and the logic pieces need to
// check for blocking/capture along their movement paths.
// ============================================================

type Board struct {
	grid [8][8]Piece // nil = empty square
}

func (b *Board) PieceAt(pos Position) Piece {
	if !pos.InBounds() {
		return nil
	}
	return b.grid[pos.Row][pos.Col]
}

func (b *Board) IsEmpty(pos Position) bool {
	return b.PieceAt(pos) == nil
}

// IsEnemy — true if there's a piece at pos AND it's a different color.
func (b *Board) IsEnemy(pos Position, color Color) bool {
	p := b.PieceAt(pos)
	return p != nil && p.Color() != color
}

func (b *Board) MovePiece(from, to Position) {
	b.grid[to.Row][to.Col] = b.grid[from.Row][from.Col]
	b.grid[from.Row][from.Col] = nil
}

// ------------------------------------------------------------
// slidingMoves — shared logic for Rook/Bishop/Queen: walk in each
// given direction until hitting the edge, a friendly piece (stop,
// don't include), or an enemy piece (stop, DO include — it's a
// legal capture).
// ------------------------------------------------------------
func slidingMoves(b *Board, from Position, color Color, directions [][2]int) []Position {
	var moves []Position
	for _, dir := range directions {
		next := Position{Row: from.Row + dir[0], Col: from.Col + dir[1]}
		for next.InBounds() {
			if b.IsEmpty(next) {
				moves = append(moves, next)
			} else if b.IsEnemy(next, color) {
				moves = append(moves, next) // capture, then stop
				break
			} else {
				break // friendly piece blocks further movement
			}
			next = Position{Row: next.Row + dir[0], Col: next.Col + dir[1]}
		}
	}
	return moves
}

// ------------------------------------------------------------
// fixedMoves — shared logic for King/Knight: a fixed set of
// offsets, no sliding, just check each landing square once.
// ------------------------------------------------------------
func fixedMoves(b *Board, from Position, color Color, offsets [][2]int) []Position {
	var moves []Position
	for _, off := range offsets {
		next := Position{Row: from.Row + off[0], Col: from.Col + off[1]}
		if !next.InBounds() {
			continue
		}
		if b.IsEmpty(next) || b.IsEnemy(next, color) {
			moves = append(moves, next)
		}
	}
	return moves
}

// ============================================================
// Rook — horizontal/vertical, any distance
// ============================================================

type Rook struct{ color Color }

func (r *Rook) Color() Color   { return r.color }
func (r *Rook) Symbol() string { return "R" }

func (r *Rook) ValidMoves(b *Board, from Position) []Position {
	directions := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	return slidingMoves(b, from, r.color, directions)
}

// ============================================================
// Bishop — diagonal, any distance
// ============================================================

type Bishop struct{ color Color }

func (bp *Bishop) Color() Color   { return bp.color }
func (bp *Bishop) Symbol() string { return "B" }

func (bp *Bishop) ValidMoves(b *Board, from Position) []Position {
	directions := [][2]int{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
	return slidingMoves(b, from, bp.color, directions)
}

// ============================================================
// Queen — Rook + Bishop combined
// ============================================================

type Queen struct{ color Color }

func (q *Queen) Color() Color   { return q.color }
func (q *Queen) Symbol() string { return "Q" }

func (q *Queen) ValidMoves(b *Board, from Position) []Position {
	directions := [][2]int{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1}, // rook directions
		{-1, -1}, {-1, 1}, {1, -1}, {1, 1}, // bishop directions
	}
	return slidingMoves(b, from, q.color, directions)
}

// ============================================================
// Knight — fixed L-shape, jumps over pieces (no path blocking)
// ============================================================

type Knight struct{ color Color }

func (n *Knight) Color() Color   { return n.color }
func (n *Knight) Symbol() string { return "N" }

func (n *Knight) ValidMoves(b *Board, from Position) []Position {
	offsets := [][2]int{
		{-2, -1}, {-2, 1}, {2, -1}, {2, 1},
		{-1, -2}, {-1, 2}, {1, -2}, {1, 2},
	}
	return fixedMoves(b, from, n.color, offsets)
}

// ============================================================
// King — one square any direction (castling deliberately omitted)
// ============================================================

type King struct{ color Color }

func (k *King) Color() Color   { return k.color }
func (k *King) Symbol() string { return "K" }

func (k *King) ValidMoves(b *Board, from Position) []Position {
	offsets := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}
	return fixedMoves(b, from, k.color, offsets)
}

// ============================================================
// Pawn — the special case: forward moves can't capture, diagonal
// moves can ONLY capture. Direction depends on color.
// (En passant and promotion deliberately omitted for now.)
// ============================================================

type Pawn struct{ color Color }

func (p *Pawn) Color() Color   { return p.color }
func (p *Pawn) Symbol() string { return "P" }

func (p *Pawn) ValidMoves(b *Board, from Position) []Position {
	var moves []Position

	dir := -1 // White moves "up" (decreasing row)
	startRow := 6
	if p.color == Black {
		dir = 1 // Black moves "down" (increasing row)
		startRow = 1
	}

	// one step forward — only if empty
	oneStep := Position{Row: from.Row + dir, Col: from.Col}
	if oneStep.InBounds() && b.IsEmpty(oneStep) {
		moves = append(moves, oneStep)

		// two steps forward — only from starting row, and only if
		// BOTH squares are empty (can't jump over a blocker)
		twoStep := Position{Row: from.Row + 2*dir, Col: from.Col}
		if from.Row == startRow && b.IsEmpty(twoStep) {
			moves = append(moves, twoStep)
		}
	}

	// diagonal captures — ONLY if an enemy piece is actually there
	for _, dc := range []int{-1, 1} {
		diag := Position{Row: from.Row + dir, Col: from.Col + dc}
		if diag.InBounds() && b.IsEnemy(diag, p.color) {
			moves = append(moves, diag)
		}
	}

	return moves
}
