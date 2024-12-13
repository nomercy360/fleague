// Пример данных. Замените на реальные данные из API.
import { cn, formatDate, timeToLocaleString } from '~/lib/utils'
import { Button } from '~/components/ui/button'
import { IconRefresh } from '~/components/icons'
import { createQuery } from '@tanstack/solid-query'
import { fetchPredictions, PredictionRequest, PredictionResponse } from '~/lib/api'
import { createSignal, For, Match, Show, Switch } from 'solid-js'
import MatchCard from '~/components/match-card'
import { Drawer, DrawerTrigger } from '~/components/ui/drawer'
import FootballScoreboard from '~/components/score-board'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '~/components/ui/tabs'


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

const UserActivity = () => {
	const query = createQuery(() => ({
		queryKey: ['predictions'],
		queryFn: () => fetchPredictions(),
	}))

	const onPredictionUpdate = () => {
		query.refetch()
	}

	const [selectedPrediction, setSelectedPrediction] = createSignal<PredictionResponse>({} as any)

	return (
		<div class="px-3">
			<p class="mb-1 px-2 text-lg font-semibold">
				My Predictions
			</p>
			<div class="space-y-2">
				<Show when={query.data}>
					<Drawer>
						<For each={query.data}>
							{(prediction: PredictionResponse) => (
								<DrawerTrigger class="w-full" onClick={() => {
									setSelectedPrediction(prediction)
								}}>
									<MatchCard match={prediction.match} prediction={prediction} />
								</DrawerTrigger>
							)}
						</For>
						<FootballScoreboard
							match={selectedPrediction().match}
							onUpdate={onPredictionUpdate}
							prediction={selectedPrediction()}
						/>
					</Drawer>
				</Show>
			</div>
		</div>
	)
}

export default UserActivity
