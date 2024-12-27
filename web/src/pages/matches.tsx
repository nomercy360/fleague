import { createSignal, For, onCleanup, onMount, Show, Suspense } from 'solid-js'
import { createQuery } from '@tanstack/solid-query'
import { fetchActiveSeason, fetchLeaderboard, fetchMatches } from '~/lib/api'

import { Drawer, DrawerTrigger } from '~/components/ui/drawer'
import { formatDate } from '~/lib/utils'
import { queryClient } from '~/App'
import MatchCard from '~/components/match-card'
import FootballScoreboard from '~/components/score-board'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '~/components/ui/tabs'
import { Link } from '~/components/link'
import MatchStats from '~/components/match-stats'

export default function MatchesPage() {
	const [selectedMatch, setSelectedMatch] = createSignal({} as any)
	const [height, setHeight] = createSignal(0)
	const query = createQuery(() => ({
		queryKey: ['matches'],
		queryFn: () => fetchMatches(),
	}))

	const leaderboardQuery = createQuery(() => ({
		queryKey: ['leaderboard'],
		queryFn: () => fetchLeaderboard(),
	}))

	const seasonQuery = createQuery(() => ({
		queryKey: ['season'],
		queryFn: () => fetchActiveSeason(),
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
		setHeight(window.Telegram.WebApp.viewportHeight - 140 - 56)
	})

	onCleanup(() => {
		// enable scroll on body when drawer is closed
		document.body.style.overflow = 'auto'
	})

	function calculateDuration(date: string) {
		// until that date from now
		// format: 2d 3h 4m
		const now = new Date()
		const endDate = new Date(date)
		const diff = endDate.getTime() - now.getTime()

		const days = Math.floor(diff / (1000 * 60 * 60 * 24))
		const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
		const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))

		return `${days}d ${hours}h ${minutes}m`
	}

	return (
		<>
			<div class="p-3 flex-col flex items-center justify-center">
				<Show when={seasonQuery.data} fallback={<div class="w-full h-20 rounded-2xl bg-secondary" />}>
					<InfoCard title={`Active Season ${seasonQuery.data.name}`}
										text={`Ends in ${calculateDuration(seasonQuery.data.end_date)}`} />
				</Show>
			</div>
			<Tabs defaultValue="preview" class="relative mr-auto w-full">
				<TabsList class="w-full justify-start rounded-none border-b bg-transparent p-0 h-[56px]">
					<TabsTrigger
						value="matches"
						class="relative h-[56px] rounded-none border-b-2 border-b-transparent bg-transparent px-4 font-semibold text-muted-foreground shadow-none transition-none data-[selected]:border-b-primary data-[selected]:text-foreground data-[selected]:shadow-none"
					>
						Matches
					</TabsTrigger>
					<TabsTrigger
						value="leaderboard"
						class="relative h-[56px] rounded-none border-b-2 border-b-transparent bg-transparent px-4 font-semibold text-muted-foreground shadow-none transition-none data-[selected]:border-b-primary data-[selected]:text-foreground data-[selected]:shadow-none"
					>
						Leaderboard
					</TabsTrigger>
				</TabsList>
				<TabsContent value="matches" class="pt-4 px-3 space-y-2 w-full overflow-y-scroll pb-[120px]"
										 style={{ height: `${height()}px` }}>
					<Drawer>
						<Show when={!query.isLoading}>
							{Object.entries(query.data).map(([date, matches]) => (
								<>
									<p class="px-2 text-base font-semibold">
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
				<TabsContent value="leaderboard" class="pt-5 px-3 space-y-2 w-full overflow-y-scroll pb-[120px]"
										 style={{ height: `${height()}px` }}>
					<Show when={leaderboardQuery.data}>
						<For each={leaderboardQuery.data}>
							{(entry) => (
								<Link class="flex items-center justify-between h-12 px-3 bg-card rounded-2xl"
											href={`/users/${entry.user.username}`}>
									<div class="flex items-center">
										<span
											class="w-4 text-center text-base font-semibold text-secondary-foreground">{getUserPosition(entry.user_id)}</span>
										<img
											src={entry.user.avatar_url}
											alt="User avatar"
											class="ml-4 size-6 rounded-full object-cover"
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
		</>
	)
}

function InfoCard({ title, text }: { title: string; text: string }) {
	return (
		<div class="w-full bg-secondary p-3 rounded-2xl flex items-center justify-center flex-col">
			<span class="material-symbols-rounded text-[48px]">
				sports_soccer
			</span>
			<h1 class="text-xl font-bold text-center">{title}</h1>
			<p class="text-sm text-secondary-foreground text-center mt-2">
				{text}
			</p>
		</div>
	)
}
