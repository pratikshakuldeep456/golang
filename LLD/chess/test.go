package chess

import "testing"

func TestChess(t *testing.T) {
	b := &Board{}
	b.grid[4][4] = &Rook{color: White}
	b.grid[6][4] = &Pawn{color: White}
	b.grid[5][4] = &Pawn{color: Black}
	b.grid[3][3] = &Knight{color: White}
	b.grid[2][3] = &Pawn{color: White} // blocker, knight should jump it

	rookMoves := b.grid[4][4].ValidMoves(b, Position{4, 4})
	if !contains(rookMoves, Position{4, 7}) || contains(rookMoves, Position{5, 5}) {
		t.Error("rook: expected straight-line moves only")
	}

	pawnMoves := b.grid[6][4].ValidMoves(b, Position{6, 4})
	if contains(pawnMoves, Position{5, 4}) {
		t.Error("pawn: should not capture straight ahead")
	}

	knightMoves := b.grid[3][3].ValidMoves(b, Position{3, 3})
	if !contains(knightMoves, Position{1, 2}) {
		t.Error("knight: should jump over blocking piece")
	}

	g := &Game{Board: b, CurrentTurn: White}
	if err := g.MakeMove(Position{6, 4}, Position{6, 4}); err == nil {
		t.Error("expected illegal move (same square) to be rejected")
	}
	if err := g.MakeMove(Position{4, 4}, Position{4, 7}); err != nil {
		t.Fatalf("expected valid rook move, got %v", err)
	}
	if g.CurrentTurn != Black {
		t.Error("expected turn to flip to Black after valid move")
	}
}
