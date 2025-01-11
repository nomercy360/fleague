package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/invopop/jsonschema"
	"github.com/labstack/echo/v4"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/terrors"
	"net/http"
	"sort"
	"strings"
	"time"
)

func getUserID(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*contract.JWTClaims)
	return claims.UID
}

func (a API) ListMatches(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := a.storage.GetActiveMatches(ctx)
	uid := getUserID(c)

	if err != nil {
		return terrors.InternalServer(err, "failed to get active matches")
	}

	var matches []contract.MatchResponse
	for _, match := range res {
		homeTeam, err := a.storage.GetTeamByID(ctx, match.HomeTeamID)
		if err != nil && errors.Is(err, db.ErrNotFound) {
			return terrors.NotFound(err, fmt.Sprintf("team with id %s not found", match.HomeTeamID))
		} else if err != nil {
			return terrors.InternalServer(err, "failed to get home team")
		}

		awayTeam, err := a.storage.GetTeamByID(ctx, match.AwayTeamID)
		if err != nil && errors.Is(err, db.ErrNotFound) {
			return terrors.NotFound(err, fmt.Sprintf("team with id %s not found", match.AwayTeamID))
		} else if err != nil {
			return terrors.InternalServer(err, "failed to get away team")
		}

		prediction, err := a.storage.GetUserPredictionByMatchID(ctx, uid, match.ID)
		if err != nil && !errors.Is(err, db.ErrNotFound) {
			return terrors.InternalServer(err, "failed to get user prediction")
		}

		resp := toMatchResponse(match, homeTeam, awayTeam)

		if prediction.UserID != "" {
			resp.Prediction = &prediction
		}

		matches = append(matches, resp)
	}

	// sort matches by date
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].MatchDate.Before(matches[j].MatchDate)
	})

	return c.JSON(http.StatusOK, matches)
}

func toMatchResponse(match db.Match, homeTeam db.Team, awayTeam db.Team) contract.MatchResponse {
	return contract.MatchResponse{
		ID:         match.ID,
		Tournament: match.Tournament,
		MatchDate:  match.MatchDate,
		Status:     match.Status,
		HomeTeam:   homeTeam,
		AwayTeam:   awayTeam,
		HomeScore:  match.HomeScore,
		AwayScore:  match.AwayScore,
		HomeOdds:   match.HomeOdds,
		DrawOdds:   match.DrawOdds,
		AwayOdds:   match.AwayOdds,
	}
}

func (a API) predictionsByUserID(ctx context.Context, uid string, onlyCompleted bool) ([]contract.PredictionResponse, error) {
	predictions, err := a.storage.GetPredictionsByUserID(ctx, uid, onlyCompleted)
	if err != nil {
		return nil, err
	}

	var res []contract.PredictionResponse
	for _, prediction := range predictions {
		match, err := a.storage.GetMatchByID(ctx, prediction.MatchID)
		if err != nil && errors.Is(err, db.ErrNotFound) {
			return nil, terrors.NotFound(err, "match not found")
		} else if err != nil {
			return nil, err
		}

		homeTeam, err := a.storage.GetTeamByID(ctx, match.HomeTeamID)
		if err != nil {
			return nil, terrors.InternalServer(err, "failed to get home team")
		}

		awayTeam, err := a.storage.GetTeamByID(ctx, match.AwayTeamID)
		if err != nil {
			return nil, terrors.InternalServer(err, "failed to get away team")
		}

		res = append(res, contract.PredictionResponse{
			UserID:             prediction.UserID,
			MatchID:            prediction.MatchID,
			PredictedOutcome:   prediction.PredictedOutcome,
			PredictedHomeScore: prediction.PredictedHomeScore,
			PredictedAwayScore: prediction.PredictedAwayScore,
			PointsAwarded:      prediction.PointsAwarded,
			CreatedAt:          prediction.CreatedAt,
			CompletedAt:        prediction.CompletedAt,
			Match:              toMatchResponse(match, homeTeam, awayTeam),
		})
	}

	// sort predictions by match date
	sort.Slice(res, func(i, j int) bool {
		return res[i].Match.MatchDate.After(res[j].Match.MatchDate)
	})

	return res, nil
}

func (a API) GetUserInfo(c echo.Context) error {
	username := c.Param("username")
	ctx := c.Request().Context()

	user, err := a.storage.GetUserByUsername(username)

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err, "user not found")
	} else if err != nil {
		return terrors.InternalServer(err, "failed to get user")
	}

	userPredictions, err := a.predictionsByUserID(ctx, user.ID, false)

	if err != nil {
		return terrors.InternalServer(err, "failed to get user predictions")
	}

	rank, err := a.storage.GetUserRank(ctx, user.ID)
	if err != nil {
		return terrors.InternalServer(err, "failed to get user rank")
	}

	resp := &contract.UserInfoResponse{
		User: contract.UserProfile{
			ID:                 user.ID,
			FirstName:          user.FirstName,
			LastName:           user.LastName,
			Username:           user.Username,
			AvatarURL:          user.AvatarURL,
			TotalPoints:        user.TotalPoints,
			TotalPredictions:   user.TotalPredictions,
			CorrectPredictions: user.CorrectPredictions,
			GlobalRank:         rank,
			FavoriteTeam:       user.FavoriteTeam,
			CurrentWinStreak:   user.CurrentWinStreak,
			LongestWinStreak:   user.LongestWinStreak,
		},
		Predictions: userPredictions,
	}

	return c.JSON(http.StatusOK, resp)
}

func (a API) ListMyReferrals(c echo.Context) error {
	res, err := a.storage.ListUserReferrals(c.Request().Context(), getUserID(c))
	if err != nil {
		return terrors.InternalServer(err, "failed to get user referrals")
	}

	var users []contract.UserProfile
	for _, user := range res {
		users = append(users, contract.UserProfile{
			ID:                 user.ID,
			FirstName:          user.FirstName,
			LastName:           user.LastName,
			Username:           user.Username,
			AvatarURL:          user.AvatarURL,
			TotalPoints:        user.TotalPoints,
			TotalPredictions:   user.TotalPredictions,
			CorrectPredictions: user.CorrectPredictions,
		})
	}

	return c.JSON(http.StatusOK, users)
}

func (a API) UpdateUser(c echo.Context) error {
	var req contract.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return terrors.BadRequest(err, "failed to decode request")
	}

	if err := req.Validate(); err != nil {
		return terrors.BadRequest(err, "failed to validate request")
	}

	uid := getUserID(c)
	ctx := c.Request().Context()

	user, err := a.storage.GetUserByID(uid)

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err, "user not found")
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.FavoriteTeamID = req.FavoriteTeamID
	if req.LanguageCode != nil {
		user.LanguageCode = req.LanguageCode
	}

	if err := a.storage.UpdateUserInformation(ctx, user); err != nil {
		return terrors.InternalServer(err, "could not update user")
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "ok"})
}

func (a API) AutoPredictMatch(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	match, err := a.storage.GetMatchByID(ctx, id)
	if err != nil {
		return terrors.InternalServer(err, "failed to fetch matches")
	}

	homeTeam, err := a.storage.GetTeamByID(ctx, match.HomeTeamID)
	if err != nil {
		return terrors.InternalServer(err, "failed to get home team")
	}

	awayTeam, err := a.storage.GetTeamByID(ctx, match.AwayTeamID)
	if err != nil {
		return terrors.InternalServer(err, "failed to get away team")
	}

	homeLastMatches, err := a.storage.GetLastMatchesByTeamID(ctx, match.HomeTeamID, 5)
	if err != nil {
		return terrors.InternalServer(err, "failed to fetch home team last matches")
	}

	awayLastMatches, err := a.storage.GetLastMatchesByTeamID(ctx, match.AwayTeamID, 5)
	if err != nil {
		return terrors.InternalServer(err, "failed to fetch away team last matches")
	}

	formatMatches := func(matches []db.Match) string {
		var results []string
		for _, m := range matches {
			var result string
			if m.HomeTeamID == match.HomeTeamID {
				result = fmt.Sprintf("vs %s: %d-%d", m.AwayTeamID, *m.HomeScore, *m.AwayScore)
			} else {
				result = fmt.Sprintf("@ %s: %d-%d", m.HomeTeamID, *m.AwayScore, *m.HomeScore)
			}
			results = append(results, result)
		}
		return strings.Join(results, ", ")
	}

	homeStats := formatMatches(homeLastMatches)
	awayStats := formatMatches(awayLastMatches)

	prompt := fmt.Sprintf(`
			Predict the outcome of the following match:
			Home Team: %s, Away Team: %s.
			Odds: Home - %.2f, Draw - %.2f, Away - %.2f.
			Date: %s.
			Home Team Last Matches: %s.
			Away Team Last Matches: %s.
		`, homeTeam.Name, awayTeam.Name, *match.HomeOdds, *match.DrawOdds, *match.AwayOdds, match.MatchDate.Format(time.RFC3339), homeStats, awayStats)

	client := openai.NewClient(
		option.WithAPIKey(a.cfg.OpenAIKey),
	)

	prediction, err := callChatGPT(ctx, client, prompt)
	if err != nil {
		return terrors.InternalServer(err, "failed to get prediction")
	}

	return c.JSON(http.StatusOK, echo.Map{"prediction": prediction, "match": toMatchResponse(match, homeTeam, awayTeam)})
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
		Model: openai.F(openai.ChatModelGPT4o2024_08_06),
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
