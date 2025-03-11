import { Match, Show, Switch } from 'solid-js'
import { IconTrophy } from '~/components/icons'
import { cn, formatDate, timeToLocaleString } from '~/lib/utils'
import { MatchResponse, PredictionResponse } from '~/lib/api'
import { store } from '~/store'

type MatchCardProps = {
	match: MatchResponse
	prediction?: PredictionResponse
}

export default function MatchCard(props: MatchCardProps) {
	const { match, prediction } = props
	const {
		home_team,
		away_team,
		match_date,
		status,
		home_score,
		away_score,
		home_odds,
		away_odds,
		draw_odds,
	} = match

	const {
		predicted_home_score,
		predicted_away_score,
		predicted_outcome,
		points_awarded,
		completed_at,
	} = prediction ?? {}

	const predictionCompleted = completed_at !== null
	const predictionCorrect = (points_awarded || 0) > 0
	const predictionLost = !!completed_at && !predictionCorrect
	const predictionEditable = !predictionCompleted

	function getPredictionText() {
		if (!prediction) return ''

		// If scores are predicted, show them as "2:0"
		if (predicted_home_score !== null && predicted_away_score !== null) {
			return `${predicted_home_score}:${predicted_away_score}`
		}

		// Otherwise show predicted outcome like "X", home short name, or away short name
		switch (predicted_outcome) {
			case 'home':
				return home_team.short_name
			case 'draw':
				return 'Draw'
			case 'away':
				return away_team.short_name
			default:
				return ''
		}
	}

	function getEditableOutcomeIcon() {
		if (predictionEditable && predicted_outcome && predicted_outcome !== 'draw') {
			return <span class="material-symbols-rounded text-primary-foreground text-[16px]">trophy</span>
		}
		return null
	}

	function getCompletedPredictionDetails() {
		if (!prediction) return null

		// If the prediction was score-based, show predicted scores
		if (predicted_home_score !== null && predicted_away_score !== null) {
			return `(${predicted_home_score}:${predicted_away_score})`
		}

		// Otherwise show outcome-based prediction in parentheses
		if (predicted_outcome === 'draw') {
			return '(X)'
		}

		if (predicted_outcome === 'home') {
			return `(${home_team.short_name})`
		}

		if (predicted_outcome === 'away') {
			return `(${away_team.short_name})`
		}

		return ''
	}

	const predictionClass = cn(
		'text-xs flex flex-row justify-center space-x-1 font-medium pt-1.5 z-0 absolute top-0 text-primary-foreground rounded-bl-full rounded-br-full left-1/2 transform -translate-x-1/2 w-36 h-8 bg-gradient-to-b to-transparent pointer-events-none',
		{
			'from-green-500': predictionCorrect,
			'from-red-500': predictionLost,
			'from-primary': predictionEditable,
			'from-gray-500': status === 'ongoing',
		},
	)

	return (
		<div class="h-[120px] relative grid grid-cols-3 items-center rounded-2xl max-w-md mx-auto p-2.5 pt-4 bg-card">
			<Show when={prediction}>
        <span class={predictionClass}>
          <Switch>
            <Match when={predictionCompleted}>
              <Switch>
                <Match when={predictionCorrect}>
                  <span>win +{points_awarded}</span>
                </Match>
                <Match when={predictionLost}>
                  <span>lost</span>
                </Match>
              </Switch>
              <span class="text-xs font-normal">
                {getCompletedPredictionDetails()}
              </span>
            </Match>
            <Match when={predictionEditable}>
              <Show when={predictionEditable}>
                {getEditableOutcomeIcon()}
              </Show>
              <span>
                {getPredictionText()}
              </span>
            </Match>
          </Switch>
        </span>
			</Show>

			<div class="z-10 flex flex-col items-center space-y-2 text-center">
				<img src={home_team.crest_url}
						 alt={home_team.name}
						 class="size-9 object-contain"
				/>
				<p class="max-w-20 text-xs text-foreground">{home_team.short_name}</p>
			</div>

			<Show when={status === 'scheduled'}>
				<div class="mb-1 flex flex-col items-center text-center justify-end self-stretch">
					<span class="leading-none text-2xl font-bold text-center">
						{timeToLocaleString(match_date, store.user?.language_code)}
					</span>
					<span class="text-xs text-center">
						{formatDate(match_date, false, store.user?.language_code)}
					</span>
					<p class="mt-1 text-xs text-muted-foreground">
						{match.tournament}
					</p>
				</div>
			</Show>

			<Show when={status === 'completed'}>
				<div class="mb-3 flex flex-col items-center text-center justify-end self-stretch">
          <span class="text-2xl font-bold text-center">
            {home_score} - {away_score}
          </span>
					<span class="text-xs text-center">
            {formatDate(match_date, false, store.user?.language_code)}
          </span>
				</div>
			</Show>

			<Show when={status === 'ongoing'}>
				<div class="mb-3 flex flex-col items-center text-center justify-end self-stretch">
					<span class="text-2xl font-bold text-center">
						{home_score} - {away_score}
					</span>
					<span class="flex items-center justify-center text-xs text-center">
						Live <span
						class="material-symbols-rounded icon-fill text-green-400 text-[16px] animate-pulse">fiber_manual_record</span>
					</span>
				</div>
			</Show>

			<div class="flex flex-col items-center space-y-2 text-center">
				<img
					src={away_team.crest_url}
					alt={away_team.name}
					class="size-9 object-contain"
				/>
				<p class="text-xs text-foreground">{away_team.short_name}</p>
			</div>
		</div>
	)
}
