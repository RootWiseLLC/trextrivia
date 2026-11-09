package jeopardy

import (
	"testing"

	"github.com/rileythomp/jeopardy/be-jeopardy/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestSurnameForms(t *testing.T) {
	tests := []struct {
		answer string
		want   []string
	}{
		{"Abraham Lincoln", []string{"Lincoln"}},
		{"Vincent van Gogh", []string{"van Gogh", "Gogh"}},
		{"Martin Luther King Jr.", []string{"Luther King", "King"}},
		{"Ruth Bader Ginsburg", []string{"Bader Ginsburg", "Ginsburg"}},
		{"Edgar Allan Poe", []string{"Allan Poe", "Poe"}},
		{"(Robert) Frost", []string{"Frost"}},
		{"(William Randolph) Hearst", []string{"Hearst"}},
		{"Elizabeth I", nil}, // regnal number, not a surname
		{"Henry VIII", nil},
		{"Cher", nil},          // one word
		{"New York", nil},      // not a first name
		{"Pearl Harbor", nil},  // "Pearl" left out of the list on purpose
		{"Marshall Plan", nil}, // ends in a thing word
		{"George Washington Bridge", nil},
		{"Charles River", nil},
		{"Duke University", nil},
		{"Jack rabbit", nil}, // lowercase word, so not a name
		{"a centaur", nil},
	}
	for _, tt := range tests {
		t.Run(tt.answer, func(t *testing.T) {
			assert.Equal(t, tt.want, nilIfEmpty(surnameForms(tt.answer)))
		})
	}
}

func nilIfEmpty(s []string) []string {
	if len(s) == 0 {
		return nil
	}
	return s
}

func TestCheckAnswerAcceptsSurname(t *testing.T) {
	q := &Question{Question: db.Question{Alternatives: []string{"Vincent van Gogh"}}}
	for _, ans := range []string{"Vincent van Gogh", "van gogh", "Gogh", "Van Goh"} {
		assert.True(t, q.checkAnswer(ans), ans)
	}
	for _, ans := range []string{"Vincent", "Monet", "van"} {
		assert.False(t, q.checkAnswer(ans), ans)
	}

	place := &Question{Question: db.Question{Alternatives: []string{"George Washington Bridge"}}}
	assert.False(t, place.checkAnswer("Bridge"))
	assert.True(t, place.checkAnswer("George Washington Bridge"))
}
