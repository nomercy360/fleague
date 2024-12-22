import { createQuery } from '@tanstack/solid-query'
import { fetchPredictions, PredictionResponse } from '~/lib/api'
import { createSignal, For, Show } from 'solid-js'
import MatchCard from '~/components/match-card'
import { Drawer, DrawerTrigger } from '~/components/ui/drawer'
import FootballScoreboard from '~/components/score-board'


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
