import { Match, Show, Switch } from 'solid-js'
import { IconTrophy } from '~/components/icons'
import { cn, timeToLocaleString } from '~/lib/utils'
import { MatchResponse, PredictionResponse } from '~/lib/api'

type MatchCardProps = {
	match: MatchResponse
	prediction?: PredictionResponse
}

export default function MatchCard(props: MatchCardProps) {
	function predictionText(prediction: PredictionResponse, match: MatchResponse) {
		// if score is predicted show it like 2:0
		if (prediction.predicted_home_score !== null) {
			return `${prediction.predicted_home_score}:${prediction.predicted_away_score}`
		} else {
			// otherwise show predicted outcome like 1, X, 2
			switch (prediction.predicted_outcome) {
				case 'home':
					return match.home_team.short_name
				case 'draw':
					return 'Draw'
				case 'away':
					return match.away_team.short_name
				default:
					return ''
			}
		}
	}

	function predictionCompleted() {
		return props.prediction?.completed_at !== null
	}

	function predictionCorrect() {
		return props.prediction?.points_awarded || 0 > 0
	}

	function predictionLost() {
		return props.prediction?.completed_at && !predictionCorrect()
	}

	function predictionEditable() {
		return props.prediction?.completed_at === null
	}

	function resolveIcon() {
		if (props.prediction?.predicted_outcome !== null && props.prediction?.predicted_outcome !== 'draw') {
			return <IconTrophy class="size-4" />
		}
	}

	return (
		<div class="h-[120px] relative grid grid-cols-3 items-center rounded-2xl max-w-md mx-auto p-2.5 pt-4 bg-card">
			<Show when={props.prediction}>
				<span
					class={cn('text-xs flex flex-row justify-center space-x-1 font-semibold pt-1.5 z-0 absolute top-0 rounded-bl-full rounded-br-full left-1/2 transform -translate-x-1/2 w-32 h-8 bg-gradient-to-b to-transparent pointer-events-none', {
						'from-green-500': predictionCorrect(),
						'from-red-500': predictionLost(),
						'from-primary': predictionEditable(),
					})}
				>
					<Switch>
						<Match when={predictionCompleted()}>
							<Switch>
								<Match when={predictionCorrect()}>win +{props.prediction?.points_awarded}</Match>
								<Match when={predictionLost()}>lost</Match>
							</Switch>
						</Match>
						<Match when={predictionEditable()}>
							<Show when={predictionEditable()}>
								{resolveIcon()}
							</Show>
							<span>
								{predictionText(props.prediction!, props.match)}
							</span>
						</Match>
					</Switch>
				</span>
			</Show>
			<div class="z-10 flex flex-col items-center space-y-2 text-center">
				<img src={`/logos/${props.match.home_team.name}.png`} alt="" class="w-10" />
				<p class="max-w-20 text-xs text-foreground">{props.match.home_team.name}</p>
			</div>
			<Show when={props.match.status == 'scheduled'}>
				<div class="mb-3 flex flex-col items-center text-center justify-end self-stretch">
				<span class="text-2xl font-bold text-center">
					{timeToLocaleString(props.match.match_date)}
				</span>
					<p class="text-xs text-muted-foreground text-center">
						{props.match.tournament}
					</p>
				</div>
			</Show>
			<Show when={props.match.status == 'completed'}>
				<div class="mb-3 flex flex-col items-center text-center justify-end self-stretch">
					<span class="text-2xl font-bold text-center">
						{props.match.home_score} - {props.match.away_score}
					</span>
					<p class="text-xs text-muted-foreground text-center">
						{props.match.tournament}
					</p>
				</div>
			</Show>
			<div class="flex flex-col items-center space-y-2 text-center">
				<img src={`/logos/${props.match.away_team.name}.png`} alt="" class="w-10" />
				<p class="text-xs text-foreground">{props.match.away_team.short_name}</p>
			</div>
		</div>
	)
}
