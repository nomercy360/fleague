import { createQuery } from '@tanstack/solid-query'
import { fetchPredictions, PredictionResponse } from '~/lib/api'
import { createSignal, For, Show } from 'solid-js'
import MatchCard from '~/components/match-card'
import { Drawer, DrawerTrigger } from '~/components/ui/drawer'
import FootballScoreboard from '~/components/score-board'
import { Link } from '~/components/link'
import { IconChevronRight } from '~/components/icons'


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
				Your Predictions
			</p>
			<div class="space-y-2">
				<Show when={query.data && !query.isLoading} fallback={<EmptyPage />}>
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

const EmptyPage = () => {
	return (
		<div class="flex flex-col items-center justify-start mt-4">
			<Link class="w-full flex flex-row h-14 justify-between items-center rounded-2xl p-3 bg-card space-x-6"
						href="/matches">
				<div>
					<p class="text-sm font-semibold">
						Make a prediction
					</p>
					<p class="text-xs text-muted-foreground font-normal">12 matches available</p>
				</div>
				<IconChevronRight class="size-6" />
			</Link>
		</div>
	)
}

export default UserActivity
