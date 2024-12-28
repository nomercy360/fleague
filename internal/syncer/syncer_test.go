package syncer_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/syncer"
	"os"
	"testing"
	"time"
)

func setupTestDB(t *testing.T) (*db.Storage, func()) {
	// Create a temporary file for SQLite
	tempFile, err := os.CreateTemp("", "test.db")
	assert.NoError(t, err)

	// Ensure the file is removed after the test
	cleanup := func() {
		os.Remove(tempFile.Name())
	}

	// Initialize the storage with the temp DB
	storage, err := db.ConnectDB(tempFile.Name())
	assert.NoError(t, err)

	// Run migrations
	//err = db.RunMigrations(storage.DB())
	//assert.NoError(t, err)

	return storage, cleanup
}

type MockNotifier struct {
	mock.Mock
}

func (m *MockNotifier) SendTextNotification(params contract.SendNotificationParams) error {
	args := m.Called(params)
	return args.Error(0)
}

func TestSyncer_ProcessPredictions(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	mockNotifier := new(MockNotifier)

	sync := syncer.NewSyncer(storage, mockNotifier, "http://test.api.url", "test_api_key")

	user := db.User{
		ID:                 "user1",
		Username:           "testuser",
		ChatID:             123456789,
		CreatedAt:          time.Now(),
		TotalPoints:        0,
		TotalPredictions:   0,
		CorrectPredictions: 0,
		GlobalRank:         1,
		CurrentWinStreak:   0,
		LongestWinStreak:   0,
	}
	err := storage.CreateUser(user)
	assert.NoError(t, err)

	match := db.Match{
		ID:         "match1",
		Tournament: "Premier League",
		HomeTeamID: "team1",
		AwayTeamID: "team2",
		MatchDate:  time.Now().Add(-24 * time.Hour), // Yesterday
		Status:     db.MatchStatusCompleted,
		HomeScore:  intPtr(2),
		AwayScore:  intPtr(1),
	}
	err = storage.SaveMatch(context.Background(), match)
	assert.NoError(t, err)

	team1 := db.Team{
		ID:           "team1",
		Name:         "Team A",
		ShortName:    "TA",
		Abbreviation: "TMA",
		CrestURL:     "http://crest.url/ta.png",
		Country:      "ENG",
	}
	team2 := db.Team{
		ID:           "team2",
		Name:         "Team B",
		ShortName:    "TB",
		Abbreviation: "TMB",
		CrestURL:     "http://crest.url/tb.png",
		Country:      "ENG",
	}
	err = storage.SaveTeam(context.Background(), team1)
	assert.NoError(t, err)
	err = storage.SaveTeam(context.Background(), team2)
	assert.NoError(t, err)

	// Pre-populate predictions
	predictions := []db.Prediction{
		// Exact correct prediction
		{
			MatchID:            "match1",
			UserID:             "user1",
			PredictedHomeScore: intPtr(2),
			PredictedAwayScore: intPtr(1),
			PredictedOutcome:   stringPtr(db.MatchOutcomeHome),
			PointsAwarded:      0,
		},
		// Correct outcome but incorrect score
		{
			MatchID:            "match1",
			UserID:             "user1",
			PredictedHomeScore: intPtr(1),
			PredictedAwayScore: intPtr(0),
			PredictedOutcome:   stringPtr(db.MatchOutcomeHome),
			PointsAwarded:      0,
		},
		// Incorrect prediction
		{
			MatchID:            "match1",
			UserID:             "user1",
			PredictedHomeScore: intPtr(0),
			PredictedAwayScore: intPtr(2),
			PredictedOutcome:   stringPtr(db.MatchOutcomeAway),
			PointsAwarded:      0,
		},
	}

	for _, p := range predictions {
		err := storage.SavePrediction(context.Background(), p)
		assert.NoError(t, err)
	}

	// Mock notifier expectations
	// Expect notifications for streak achievements if applicable
	// For this test, assuming streaks are not long enough to trigger bonuses
	// If exact predictions influence streaks, adjust accordingly
	// For simplicity, not expecting notifications here
	mockNotifier.On("SendTextNotification", mock.Anything).Return(nil)

	// Run ProcessPredictions
	err = sync.ProcessPredictions(context.Background())
	assert.NoError(t, err)

	// Verify that predictions have updated ResultPoints
	updatedPredictions, err := storage.GetPredictionsForMatch(context.Background(), "match1")
	assert.NoError(t, err)
	assert.Len(t, updatedPredictions, 3)

	for _, p := range updatedPredictions {
		switch {
		case *p.PredictedHomeScore == 2 && *p.PredictedAwayScore == 1:
			// Exact correct prediction
			assert.Equal(t, 7, p.PointsAwarded)
		case *p.PredictedOutcome == db.MatchOutcomeHome:
			// Correct outcome but incorrect score
			assert.Equal(t, 3, p.PointsAwarded)
		default:
			// Incorrect prediction
			assert.Equal(t, 0, p.PointsAwarded)
		}
	}

	// Verify that user's points and streaks are updated
	updatedUser, err := storage.GetUserByID("user1")
	assert.NoError(t, err)
	// Total points: 7 + 3 + 0 = 10
	assert.Equal(t, 10, updatedUser.TotalPoints)
	// Total predictions: 3
	assert.Equal(t, 3, updatedUser.TotalPredictions)
	// Correct predictions: 2
	assert.Equal(t, 2, updatedUser.CorrectPredictions)
	// Current streak: 2 (assuming correct predictions are consecutive)
	assert.Equal(t, 2, updatedUser.CurrentWinStreak)
	// Longest streak: 2
	assert.Equal(t, 2, updatedUser.LongestWinStreak)

	// Verify that notifier was not called (since streak < 4)
	mockNotifier.AssertNotCalled(t, "SendTextNotification")
}

func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}
