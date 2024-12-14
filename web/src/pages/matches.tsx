import { createSignal, For, onCleanup, onMount, Show, Suspense } from 'solid-js'
import { createQuery } from '@tanstack/solid-query'
import { fetchLeaderboard, fetchMatches } from '~/lib/api'

import { Drawer, DrawerTrigger } from '~/components/ui/drawer'
import { formatDate } from '~/lib/utils'
import { queryClient } from '~/App'
import MatchCard from '~/components/match-card'
import FootballScoreboard from '~/components/score-board'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '~/components/ui/tabs'
import { Link } from '~/components/link'

export default function MatchesPage() {
	const [selectedMatch, setSelectedMatch] = createSignal({} as any)

	const query = createQuery(() => ({
		queryKey: ['matches'],
		queryFn: () => fetchMatches(),
	}))

	const leaderboardQuery = createQuery(() => ({
		queryKey: ['leaderboard'],
		queryFn: () => fetchLeaderboard(),
	}))

	const onPredictionUpdate = () => {
		query.refetch()
		queryClient.invalidateQueries({ queryKey: ['predictions'] })
	}

	const getUserPosition = (id: number) => {
		const idx = leaderboardQuery.data.findIndex((u: any) => u.user_id === id)

		if (idx === 0) return '🥇'
		if (idx === 1) return '🥈'
		if (idx === 2) return '🥉'
		return idx + 1
	}

	onMount(() => {
		// disable scroll on body when drawer is open
		document.body.style.overflow = 'hidden'
	})

	onCleanup(() => {
		// enable scroll on body when drawer is closed
		document.body.style.overflow = 'auto'
	})

	return (
		<div>
			<div class="px-5 pt-6">
				<div class="rounded-2xl bg-secondary w-full p-4">
					<p class="text-base font-semibold text-card-foreground uppercase tracking-widest">
						Upcoming Matches
					</p>
					<p class="mt-1 text-sm text-muted-foreground">
						Make predictions to earn points and appear on the leaderboard
					</p>
				</div>
			</div>
			<Tabs defaultValue="preview" class="mt-6 relative mr-auto w-full">
				<div class="flex items-center justify-between">
					<TabsList class="w-full justify-start rounded-none border-b bg-transparent p-0 h-14">
						<TabsTrigger
							value="matches"
							class="relative h-14 rounded-none border-b-2 border-b-transparent bg-transparent px-4 font-semibold text-muted-foreground shadow-none transition-none data-[selected]:border-b-primary data-[selected]:text-foreground data-[selected]:shadow-none"
						>
							Matches
						</TabsTrigger>
						<TabsTrigger
							value="leaderboard"
							class="relative h-14 rounded-none border-b-2 border-b-transparent bg-transparent px-4 font-semibold text-muted-foreground shadow-none transition-none data-[selected]:border-b-primary data-[selected]:text-foreground data-[selected]:shadow-none"
						>
							Leaderboard
						</TabsTrigger>
					</TabsList>
				</div>
				<TabsContent value="matches" class="pb-12 pt-2 px-3 space-y-2 w-full overflow-y-scroll h-[400px]">
					<Drawer>
						<Show when={!query.isLoading}>
							{Object.entries(query.data).map(([date, matches]) => (
								<>
									<p class="mb-1 mt-5 px-2 text-base font-semibold">
										{formatDate(date)}
									</p>
									<For each={matches as any}>
										{match => (
											<DrawerTrigger class="w-full" onClick={() => {
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
				</TabsContent>
				<TabsContent value="leaderboard" class="pb-12 pt-2 px-3 space-y-2 w-full overflow-y-scroll h-[400px]">
					<Show when={leaderboardQuery.data}>
						<For each={leaderboardQuery.data}>
							{(entry) => (
								<Link class="flex items-center justify-between h-12 px-3 bg-card rounded-2xl"
											href={`/users/${entry.user.username}`}>
									<div class="flex items-center">
										<span
											class="w-4 text-base font-semibold text-secondary-foreground">{getUserPosition(entry.user_id)}</span>
										<img
											src={entry.user.avatar_url}
											alt="User avatar"
											class="ml-3 size-6 rounded-full object-cover"
										/>
										<p class="text-base font-semibold ml-2">
											{entry.user?.first_name}{' '}{entry.user?.last_name}
										</p>
									</div>
									<div class="flex items-center">
										<p class="text-base font-semibold mr-2">{entry.points} DPS</p>
									</div>
								</Link>
							)}
						</For>
					</Show>
				</TabsContent>
			</Tabs>
		</div>
	)
}
