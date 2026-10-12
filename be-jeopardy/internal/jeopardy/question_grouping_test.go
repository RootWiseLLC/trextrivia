package jeopardy

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/rileythomp/jeopardy/be-jeopardy/internal/db"
	"github.com/stretchr/testify/assert"
)

// fakeQuestionDB returns a canned question list so setQuestions can be exercised
// without a database.
type fakeQuestionDB struct {
	questions []db.Question
}

func (f *fakeQuestionDB) GetQuestions(_ context.Context, _, _ int) ([]db.Question, error) {
	return f.questions, nil
}
func (f *fakeQuestionDB) GetCategoryQuestions(_ context.Context, _ db.Category) ([]db.Question, error) {
	return nil, nil
}
func (f *fakeQuestionDB) AddAlternative(_ context.Context, _, _ string) error { return nil }
func (f *fakeQuestionDB) AddIncorrect(_ context.Context, _, _ string) error   { return nil }
func (f *fakeQuestionDB) SaveGameAnalytics(_ context.Context, _ uuid.UUID, _ int64, _ db.AnalyticsRound, _ db.AnalyticsRound) error {
	return nil
}
func (f *fakeQuestionDB) IncrementPlayerGames(_ context.Context, _ string, _, _, _, _ int) error {
	return nil
}
func (f *fakeQuestionDB) Close() {}

func category(round int, title string, values ...int) []db.Question {
	qs := []db.Question{}
	for _, v := range values {
		qs = append(qs, db.Question{
			Round:    round,
			Value:    v,
			Category: title,
			Clue:     fmt.Sprintf("%s clue %d", title, v),
			Answer:   fmt.Sprintf("%s answer %d", title, v),
		})
	}
	return qs
}

func fullRound(round int, base int, titles ...string) []db.Question {
	qs := []db.Question{}
	for _, title := range titles {
		qs = append(qs, category(round, title, base, base*2, base*3, base*4, base*5)...)
	}
	return qs
}

func TestSetQuestionsGroupsByCategory(t *testing.T) {
	t.Run("every category title matches its own clues", func(t *testing.T) {
		questions := fullRound(1, 200, "R1A", "R1B", "R1C", "R1D", "R1E", "R1F")
		questions = append(questions, fullRound(2, 400, "R2A", "R2B", "R2C", "R2D", "R2E", "R2F")...)
		questions = append(questions, db.Question{Round: 3, Category: "FINAL", Clue: "final clue"})

		g := Game{jeopardyDB: &fakeQuestionDB{questions: questions}, GameConfig: GameConfig{FullGame: true}}
		assert.NoError(t, g.setQuestions(context.Background()))

		assert.Len(t, g.FirstRound, numCategories)
		assert.Len(t, g.SecondRound, numCategories)
		for _, round := range [][]Category{g.FirstRound, g.SecondRound} {
			for _, c := range round {
				assert.Len(t, c.Questions, numQuestions)
				for _, q := range c.Questions {
					assert.Equal(t, c.Title, q.Category)
				}
			}
		}
		assert.Equal(t, "final clue", g.FinalQuestion.Clue)
	})

	t.Run("a short category does not shift the titles of later ones", func(t *testing.T) {
		// R1B comes back with 4 clues instead of 5, which is what used to push
		// every later category one clue out of step with its header.
		questions := category(1, "R1A", 200, 400, 600, 800, 1000)
		questions = append(questions, category(1, "R1B", 200, 400, 600, 800)...)
		questions = append(questions, fullRound(1, 200, "R1C", "R1D", "R1E", "R1F", "R1G")...)

		g := Game{jeopardyDB: &fakeQuestionDB{questions: questions}}
		assert.NoError(t, g.setQuestions(context.Background()))

		assert.Len(t, g.FirstRound, numCategories)
		for _, c := range g.FirstRound {
			assert.Len(t, c.Questions, numQuestions)
			for _, q := range c.Questions {
				assert.Equal(t, c.Title, q.Category)
			}
		}
		// The incomplete category is dropped rather than merged into its neighbours.
		for _, c := range g.FirstRound {
			assert.NotEqual(t, "R1B", c.Title)
		}
	})

	t.Run("errors instead of serving a short board", func(t *testing.T) {
		questions := fullRound(1, 200, "R1A", "R1B", "R1C")
		g := Game{jeopardyDB: &fakeQuestionDB{questions: questions}}
		assert.Error(t, g.setQuestions(context.Background()))
	})
}
