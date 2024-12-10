import { createSignal, For, Show } from 'solid-js'
import { createQuery } from '@tanstack/solid-query'
import { fetchMatches } from '~/lib/api'

import { Drawer, DrawerTrigger } from '~/components/ui/drawer'
import { formatDate } from '~/lib/utils'
import { queryClient } from '~/App'
import MatchCard from '~/components/match-card'
import FootballScoreboard from '~/components/score-board'


export default function MatchesPage() {
	const [selectedMatch, setSelectedMatch] = createSignal({} as any)

	const query = createQuery(() => ({
		queryKey: ['matches'],
		queryFn: () => fetchMatches(),
	}))

	const onPredictionUpdate = () => {
		query.refetch()
		queryClient.invalidateQueries({ queryKey: ['predictions'] })
	}

	return (
		<div class="p-3 space-y-2">
			<Drawer>
				<Show when={!query.isLoading}>
					{Object.entries(query.data).map(([date, matches]) => (
						<>
							<p class="mt-6 mb-1 px-2 text-lg font-semibold">
								{formatDate(date)}
							</p>
							<For each={matches as any}>
								{match => (
									<DrawerTrigger onClick={() => {
										setSelectedMatch(match)
									}}>
										<MatchCard match={match} prediction={match.prediction} />
									</DrawerTrigger>)
								}
							</For>
						</>
					))}
				</Show>
				<FootballScoreboard
					match={selectedMatch()}
					onUpdate={onPredictionUpdate}
					prediction={selectedMatch().prediction}
				/>
			</Drawer>
		</div>
	)
}
