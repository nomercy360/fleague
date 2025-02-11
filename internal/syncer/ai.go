package syncer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/terrors"
	"log"
	"strings"
	"time"
)

const (
	aiUserID = "ai1"
)

func strPtr(s string) *string {
	return &s
}

func (s *Syncer) PredictWeeklyMatches(ctx context.Context) error {
	user, err := s.storage.GetUserByID(aiUserID)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		user = db.User{
			ID:           aiUserID,
			FirstName:    strPtr("AI Predictor"),
			LastName:     strPtr("Bot"),
			Username:     aiUserID,
			AvatarURL:    strPtr("https://assets.peatch.io/ai.svg"),
			LanguageCode: strPtr("en"),
			ChatID:       123456789,
		}
		if err := s.storage.CreateUser(user); err != nil {
			return fmt.Errorf("failed to create AI user: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	matches, err := s.storage.GetActiveMatches(ctx, aiUserID)

	for _, match := range matches {
		// skip if match is already predicted
		if match.Prediction != nil {
			log.Printf("Match %s already predicted by AI", match.ID)
			continue
		}

		if err := s.predictMatch(ctx, match); err != nil {
			return fmt.Errorf("failed to predict match: %w", err)
		}
	}

	return nil
}

func (s *Syncer) predictMatch(ctx context.Context, match db.Match) error {
	homeLastMatches, err := s.storage.GetLastMatchesByTeamID(ctx, match.HomeTeamID, 5)
	if err != nil {
		return fmt.Errorf("failed to fetch last matches for home team: %w", err)
	}

	awayLastMatches, err := s.storage.GetLastMatchesByTeamID(ctx, match.AwayTeamID, 5)
	if err != nil {
		return fmt.Errorf("failed to fetch last matches for away team: %w", err)
	}

	formatMatches := func(matches []db.Match) string {
		var results []string
		for _, m := range matches {
			var result string
			if m.HomeTeamID == match.HomeTeamID {
				result = fmt.Sprintf("vs %s: %d-%d", m.AwayTeam.Name, *m.HomeScore, *m.AwayScore)
			} else {
				result = fmt.Sprintf("@ %s: %d-%d", m.HomeTeam.Name, *m.AwayScore, *m.HomeScore)
			}
			results = append(results, result)
		}
		return strings.Join(results, ", ")
	}

	homeStats := formatMatches(homeLastMatches)
	awayStats := formatMatches(awayLastMatches)

	// skip if there is no odds
	if match.HomeOdds == nil || match.DrawOdds == nil || match.AwayOdds == nil {
		log.Printf("Skipping match %s as odds are not available", match.ID)
		return nil
	}

	prompt := fmt.Sprintf(`
			Predict the outcome of the following match:
			Home Team: %s, Away Team: %s.
			Odds: Home - %.2f, Draw - %.2f, Away - %.2f.
			Date: %s.
			Home Team Last Matches: %s.
			Away Team Last Matches: %s.
		`, match.HomeTeam.Name, match.AwayTeam.Name, *match.HomeOdds, *match.DrawOdds, *match.AwayOdds, match.MatchDate.Format(time.RFC3339), homeStats, awayStats)

	client := openai.NewClient(
		option.WithAPIKey(s.cfg.OpenAIKey),
	)

	prediction, err := callChatGPT(ctx, client, prompt)
	if err != nil {
		return terrors.InternalServer(err, "failed to get prediction")
	}

	prediction, err = callChatGPT(ctx, client, prompt)
	if err != nil {
		return fmt.Errorf("failed to get prediction from ChatGPT: %w", err)
	}

	log.Printf("AI prediction for match %s: %+v", match.ID, prediction)

	return s.storage.SavePrediction(ctx, db.Prediction{
		UserID:             aiUserID,
		MatchID:            match.ID,
		PredictedHomeScore: intPtr(prediction.HomeScore),
		PredictedAwayScore: intPtr(prediction.AwayScore),
	})
}

func intPtr(i int) *int {
	return &i
}

type MatchPrediction struct {
	Outcome    string `json:"outcome" jsonschema_description:"The predicted outcome: 'home', 'away', or 'draw'"`
	HomeScore  int    `json:"home_score" jsonschema_description:"Predicted score for the home team"`
	AwayScore  int    `json:"away_score" jsonschema_description:"Predicted score for the away team"`
	Confidence string `json:"confidence" jsonschema_description:"Confidence level or short reasoning for the prediction"`
}

func GenerateSchema[T any]() interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

var MatchPredictionResponseSchema = GenerateSchema[MatchPrediction]()

func callChatGPT(ctx context.Context, client *openai.Client, prompt string) (MatchPrediction, error) {
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("match_prediction"),
		Description: openai.F("Predicted outcome of a football match"),
		Schema:      openai.F(MatchPredictionResponseSchema),
		Strict:      openai.Bool(true),
	}

	chat, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a sports prediction assistant. Predict match outcomes based on provided details."),
			openai.UserMessage(prompt),
		}),
		Model: openai.F(openai.ChatModelGPT4oMini2024_07_18),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
			openai.ResponseFormatJSONSchemaParam{
				Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: openai.F(schemaParam),
			},
		),
	})

	if err != nil {
		return MatchPrediction{}, err
	}

	// Extract the response into the MatchPrediction struct
	var prediction MatchPrediction
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &prediction)
	if err != nil {
		return MatchPrediction{}, err
	}

	return prediction, nil
}
