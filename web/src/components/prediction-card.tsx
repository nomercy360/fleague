import { createQuery } from '@tanstack/solid-query'
import { fetchPredictions, PredictionResponse } from '~/lib/api'
import { createSignal, For, Show } from 'solid-js'
import MatchCard from '~/components/match-card'
import { Drawer, DrawerTrigger } from '~/components/ui/drawer'
import FootballScoreboard from '~/components/score-board'
import MatchStats from '~/components/match-stats'
import { useTranslations } from '~/lib/locale-context'

const UserActivity = () => {
	const query = createQuery(() => ({
		queryKey: ['predictions'],
		queryFn: () => fetchPredictions(),
	}))

	const onPredictionUpdate = () => {
		query.refetch()
	}

	const [selectedPrediction, setSelectedPrediction] = createSignal<PredictionResponse>({} as any)
	const [visibleCount, setVisibleCount] = createSignal(5)

	const { t } = useTranslations()

	return (
		<div class="px-3">
			<p class="mb-1 px-2 text-lg font-semibold">
				{t('your_predictions')}
			</p>
			<div class="space-y-2">
				<Show when={query.data && !query.isLoading}>
					<Drawer>
						<For each={query.data.slice(0, visibleCount())}>
							{(prediction: PredictionResponse) => (
								<DrawerTrigger class="w-full" onClick={() => {
									setSelectedPrediction(prediction)
								}}>
									<MatchCard match={prediction.match} prediction={prediction} />
								</DrawerTrigger>
							)}
						</For>
						<Show when={selectedPrediction()?.match?.status === 'scheduled'}>
							<FootballScoreboard
								match={selectedPrediction().match}
								onUpdate={onPredictionUpdate}
								prediction={selectedPrediction()}
							/>
						</Show>
						<Show when={selectedPrediction()?.match?.status === 'completed'}>
							<MatchStats match={selectedPrediction().match} />
						</Show>
					</Drawer>
					<Show when={query.data.length > visibleCount()}>
						<div class="text-center mt-2">
							<button
								class="px-4 text-sm font-medium h-10"
								onClick={() => setVisibleCount(visibleCount() + 5)}
							>
								{t('show_more')}
							</button>
						</div>
					</Show>
				</Show>
			</div>
		</div>
	)
}

export default UserActivity
