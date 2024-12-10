// Пример данных. Замените на реальные данные из API.
import { cn, formatDate, timeToLocaleString } from '~/lib/utils'
import { Button } from '~/components/ui/button'
import { IconRefresh } from '~/components/icons'
import { createQuery } from '@tanstack/solid-query'
import { fetchPredictions, Prediction, PredictionResponse } from '~/lib/api'
import { For, Match, Show, Switch } from 'solid-js'


const upcomingMatches = [
	{
		tournament: 'UEFA Champions League',
		homeTeam: 'Club Brugge KV',
		awayTeam: 'Liverpool FC',
		matchDate: '2024-12-11 17:45',
		predictedWinner: 'Liverpool FC',
		prediction: '2:0',
		isEditable: true,
	},
	{
		tournament: 'UEFA Champions League',
		homeTeam: 'FC Internazionale Milano',
		awayTeam: 'FC Red Bull Salzburg',
		matchDate: '2024-12-11 17:45',
		predictedWinner: 'FC Internazionale Milano',
		prediction: '3:0',
		isEditable: false,
	},
]

const PredictionCard = () => {
	const query = createQuery(() => ({
		queryKey: ['predictions'],
		queryFn: () => fetchPredictions(),
	}))

	return (
		<div class="px-4">
			<p class="mb-1 px-2 text-lg font-semibold">
				My Predictions
			</p>
			<div class="space-y-2">
				<Show when={query.data}>
					<For each={query.data}>
						{(prediction: PredictionResponse) => (
							<div
								class={cn('h-[120px] relative grid grid-cols-3 items-center rounded-2xl max-w-md mx-auto p-2.5 pt-4 bg-card', {
									'border-l-4 border-green-500': prediction.points_awarded > 0,
									'border-l-4 border-red-500': prediction.points_awarded == 0 && prediction.completed_at,
									'border-l-4 border-gray-500': prediction.points_awarded == 0 && !prediction.completed_at && prediction.match.status == 'ongoing',
									'border-l-4 border-primary': prediction.points_awarded == 0 && !prediction.completed_at && prediction.match.status == 'scheduled',
								})}>
								<div
									class={cn('h-6 rounded-b-xl w-12 flex items-center justify-center text-white text-xs font-semibold absolute top-0 left-1/2 -translate-x-1/2 transform', {
										'bg-green-500': prediction.points_awarded > 0,
										'bg-red-500': prediction.points_awarded == 0 && prediction.completed_at,
									})}>
									{prediction.points_awarded > 0 ? `win +${prediction.points_awarded}` : 'lose'}
								</div>
								<div class="flex flex-col items-center space-y-2 text-center">
									<img src={`/logos/${prediction.match.home_team.name}.png`} alt="" class="w-8" />
									<p class="text-sm font-bold">{prediction.match.home_team.short_name}</p>
								</div>
								<div class="flex flex-col items-center text-center">
									<p class="text-xs text-muted-foreground text-center">
										{prediction.match.tournament}
									</p>
									<Switch>
										<Match when={prediction.completed_at}>
											<span class="text-2xl font-bold text-center">
												{prediction.match.home_score} - {prediction.match.away_score}
											</span>
										</Match>
										<Match when={!prediction.completed_at}>
											<span class="text-lg font-bold text-center">
												{timeToLocaleString(prediction.match.match_date)}
											</span>
											<span class="text-xs text-muted-foreground text-center">
												{formatDate(prediction.match.match_date)}
											</span>
										</Match>
									</Switch>
								</div>
								<div class="flex flex-col items-center space-y-2 text-center">
									<img src={`/logos/${prediction.match.away_team.name}.png`} alt="" class="w-8" />
									<p class="text-sm font-bold">{prediction.match.away_team.short_name}</p>
								</div>
							</div>
						)}
					</For>
				</Show>
			</div>
		</div>
	)
}

export default PredictionCard
