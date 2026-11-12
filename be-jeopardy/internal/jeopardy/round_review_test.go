package jeopardy

import (
	"context"
	"testing"

	"github.com/rileythomp/jeopardy/be-jeopardy/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newReviewTestGame returns a full game with two players where every first-round
// clue except the last one has been played, and the last one is on the board.
func newReviewTestGame(t *testing.T) (*Game, *Player, *Player) {
	t.Helper()
	questions := fullRound(1, 200, "R1A", "R1B", "R1C", "R1D", "R1E", "R1F")
	questions = append(questions, fullRound(2, 400, "R2A", "R2B", "R2C", "R2D", "R2E", "R2F")...)
	questions = append(questions, db.Question{Round: 3, Category: "FINAL", Clue: "final clue", Answer: "final answer"})

	config, err := NewConfig(true, true, 0, 30, 15, 15, 30, nil, nil)
	require.NoError(t, err)
	g, err := NewGame(context.Background(), &fakeQuestionDB{questions: questions}, config)
	require.NoError(t, err)

	p1, p2 := NewPlayer("p1", "", ""), NewPlayer("p2", "", "")
	g.Players = []GamePlayer{p1, p2}
	for _, c := range g.FirstRound {
		for _, q := range c.Questions {
			q.CanChoose = false
		}
	}
	last := g.FirstRound[5].Questions[4]
	last.CanChoose = true
	last.Alternatives = []string{last.Answer}
	g.CurQuestion = last
	g.OfficialAnswer = last.Answer
	g.LastToPick = p1
	g.setState(RecvAns, p1)
	t.Cleanup(func() { g.cancelPickTimeout(); g.cancelDisputeTimeout() })
	return g, p1, p2
}

func TestRoundEndWaitsForReview(t *testing.T) {
	ctx := context.Background()

	t.Run("last clue answered correctly holds the board before the next round", func(t *testing.T) {
		g, p1, _ := newReviewTestGame(t)
		answer := g.CurQuestion.Answer

		assert.NoError(t, g.processAnswer(ctx, p1, answer))

		assert.Equal(t, RecvPick, g.State)
		assert.Equal(t, FirstRound, g.Round)
		assert.Equal(t, roundReviewTimeout, g.RoundReviewTimeout)
		assert.Equal(t, answer, g.OfficialAnswer)
		for _, p := range g.Players {
			assert.False(t, p.canPick(), "nobody picks during the review")
		}

		g.handleRoundEnd(ctx)
		assert.Equal(t, SecondRound, g.Round)
		assert.Equal(t, 0, g.RoundReviewTimeout)
	})

	t.Run("a wrong answer on the last clue can be disputed before moving on", func(t *testing.T) {
		g, p1, p2 := newReviewTestGame(t)

		assert.NoError(t, g.processAnswer(ctx, p1, "something else"))
		assert.Equal(t, RecvBuzz, g.State, "p2 can still buzz")
		assert.NoError(t, g.processBuzz(ctx, p2, true))

		assert.Equal(t, RecvPick, g.State)
		assert.Equal(t, FirstRound, g.Round)
		assert.Equal(t, roundReviewTimeout, g.RoundReviewTimeout)

		assert.NoError(t, g.initDispute(p1))
		assert.Equal(t, RecvDispute, g.State)
		assert.NoError(t, g.processDispute(ctx, p2, true))

		// Back to the review, with the dispute applied and still in the first round.
		assert.Equal(t, RecvPick, g.State)
		assert.Equal(t, FirstRound, g.Round)
		assert.Equal(t, roundReviewTimeout, g.RoundReviewTimeout)
		assert.True(t, g.CurQuestion.Answers[0].Overturned)
		assert.Equal(t, 1000, p1.score(), "the -1000 penalty is reversed and 1000 awarded")
	})

	t.Run("the last clue of the second round goes to Final after the review", func(t *testing.T) {
		g, p1, p2 := newReviewTestGame(t)
		p1.addToScore(1000)
		p2.addToScore(1000)
		g.Round = SecondRound
		g.SecondRound = g.FirstRound

		assert.NoError(t, g.processAnswer(ctx, p1, g.CurQuestion.Answer))
		assert.Equal(t, RecvPick, g.State)
		assert.Equal(t, SecondRound, g.Round)

		g.handleRoundEnd(ctx)
		assert.Equal(t, FinalRound, g.Round)
		assert.Equal(t, RecvWager, g.State)
	})
}
