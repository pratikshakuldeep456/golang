package chess

import "errors"

var (
	ErrNotYourTurn  = errors.New("not your turn")
	ErrNoPiece      = errors.New("no piece at source square")
	ErrIllegalMove  = errors.New("illegal move")
	ErrOwnKingCheck = errors.New("move would leave your own king in check")
)

type GameStatus int

const (
	Active GameStatus = iota
	Check
	Checkmate
	Stalemate
)

type Game struct {
	Board       *Board
	CurrentTurn Color
	Status      GameStatus
}

// MakeMove — the full validation pipeline for a single move.
func (g *Game) MakeMove(from, to Position) error {
	piece := g.Board.PieceAt(from)

	// 1. is there actually a piece here?
	if piece == nil {
		return ErrNoPiece
	}

	// 2. is it this player's piece?
	if piece.Color() != g.CurrentTurn {
		return ErrNotYourTurn
	}

	// 3. is `to` among the piece's own valid moves (pattern + blocking)?
	if !contains(piece.ValidMoves(g.Board, from), to) {
		return ErrIllegalMove
	}

	// 4. simulate the move — does it leave OUR OWN king in check?
	//    this is why ValidMoves alone isn't enough: a move can be
	//    legal for the piece's pattern but still illegal because it
	//    exposes your own king.
	if g.wouldLeaveKingInCheck(from, to, piece.Color()) {
		return ErrOwnKingCheck
	}

	// 5. all checks passed — actually apply the move
	g.Board.MovePiece(from, to)

	// 6. flip turn
	if g.CurrentTurn == White {
		g.CurrentTurn = Black
	} else {
		g.CurrentTurn = White
	}

	// 7. update status for the NEXT player (did we just check them?)
	g.updateStatus()

	return nil
}

// wouldLeaveKingInCheck — the simulate-then-check-then-undo pattern.
// This is the standard way to validate "does this move expose my king,"
// since you can't know without actually trying the move on the board.
func (g *Game) wouldLeaveKingInCheck(from, to Position, color Color) bool {
	captured := g.Board.PieceAt(to) // remember what was there, to restore it

	g.Board.MovePiece(from, to) // simulate

	inCheck := g.isKingInCheck(color)

	// undo the simulation — restore board to exactly how it was
	g.Board.MovePiece(to, from)
	g.Board.grid[to.Row][to.Col] = captured

	return inCheck
}

// isKingInCheck — find this color's king, then check if ANY enemy
// piece's valid moves include the king's square.
func (g *Game) isKingInCheck(color Color) bool {
	kingPos, found := g.findKing(color)
	if !found {
		return false // shouldn't happen in a real game
	}

	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			p := g.Board.grid[row][col]
			if p == nil || p.Color() == color {
				continue // skip empty squares and our own pieces
			}
			enemyMoves := p.ValidMoves(g.Board, Position{row, col})
			if contains(enemyMoves, kingPos) {
				return true
			}
		}
	}
	return false
}

func (g *Game) findKing(color Color) (Position, bool) {
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			p := g.Board.grid[row][col]
			if p == nil {
				continue
			}
			if _, isKing := p.(*King); isKing && p.Color() == color {
				return Position{row, col}, true
			}
		}
	}
	return Position{}, false
}

// updateStatus — after a move, check if the NEW current player is in
// check, and if so, whether they have ANY legal move at all (checkmate)
// or just some (check, game continues).
func (g *Game) updateStatus() {
	inCheck := g.isKingInCheck(g.CurrentTurn)
	hasMove := g.hasAnyLegalMove(g.CurrentTurn)

	switch {
	case inCheck && !hasMove:
		g.Status = Checkmate
	case inCheck:
		g.Status = Check
	case !hasMove:
		g.Status = Stalemate
	default:
		g.Status = Active
	}
}

// hasAnyLegalMove — brute-force: try every piece's every candidate
// move, filtering out ones that would leave the king in check.
// O(pieces * moves_per_piece) — fine at chess-board scale (max ~16
// pieces, max ~27 moves for a queen).
func (g *Game) hasAnyLegalMove(color Color) bool {
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			p := g.Board.grid[row][col]
			if p == nil || p.Color() != color {
				continue
			}
			from := Position{row, col}
			for _, to := range p.ValidMoves(g.Board, from) {
				if !g.wouldLeaveKingInCheck(from, to, color) {
					return true
				}
			}
		}
	}
	return false
}

func contains(positions []Position, target Position) bool {
	for _, p := range positions {
		if p == target {
			return true
		}
	}
	return false
}
