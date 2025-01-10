import { createSignal, For, Show } from 'solid-js'
import { createQuery } from '@tanstack/solid-query'
import { fetchActiveSeason, fetchLeaderboard, fetchMatches } from '~/lib/api'

import { Drawer, DrawerTrigger } from '~/components/ui/drawer'
import { formatDate } from '~/lib/utils'
import { queryClient, setShowCommunityPopup, showCommunityPopup } from '~/App'
import MatchCard from '~/components/match-card'
import FootballScoreboard from '~/components/score-board'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '~/components/ui/tabs'
import { Link } from '~/components/link'
import { Button } from '~/components/ui/button'

import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '~/components/ui/dialog'

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

	const closePopup = () => {
		setShowCommunityPopup(false)
		window.Telegram.WebApp.CloudStorage.setItem('fb_community_popup', 'closed')
	}

	return (
		<>
			<div class="p-3 flex-col flex items-center justify-center space-y-3">
				<Show when={showCommunityPopup()}>
					<div class="w-full bg-secondary p-3 rounded-2xl flex flex-col items-center relative">
						<Button
							size="icon"
							variant="ghost"
							class="absolute top-1 right-1"
							onClick={closePopup}
						>
							<span class="material-symbols-rounded text-[20px] text-primary-foreground">
								close
							</span>
						</Button>
						<span class="material-symbols-rounded text-[40px] text-blue-400">
							people
						</span>
						<h1 class="tracking-wider text-lg uppercase font-bold">
							Join community
						</h1>
						<p class="text-sm text-secondary-foreground text-center">
							To discuss matches and get the latest updates
						</p>
						<Button
							class="w-full mt-3"
							onClick={() => {
								window.Telegram.WebApp.openTelegramLink(
									'https://t.me/match_predict_league',
								)
							}}
						>
							Open in Telegram
						</Button>
					</div>
				</Show>
				<Show when={seasonQuery.data && !showCommunityPopup()}
							fallback={<div class="w-full rounded-2xl bg-secondary" />}>
					<Dialog>
						<DialogContent>
							<DialogHeader>
								<DialogTitle>🏅 Monthly Seasons!</DialogTitle>
								<DialogDescription>
									<img
										class="mb-1 mt-4 rounded-xl h-[300px] w-full object-cover"
										src="/preview.jpg"
										alt="T-shirt Prize"
									/>
									<p class="mb-4 text-xs">
										Season {seasonQuery.data.name} prize - "Ural" FC T-Shirt
									</p>
									<p class="text-sm">
										Compete for the top spot each month! Points reset monthly, and the first-place winner gets a prize.
										🏆 Make your predictions count!
									</p>
								</DialogDescription>
							</DialogHeader>
						</DialogContent>
						<div class="relative w-full bg-secondary p-3 rounded-2xl flex items-center justify-start flex-col gap-1">
							<DialogTrigger class="size-8 absolute top-1 right-1">
								<span class="material-symbols-rounded text-[20px] text-secondary-foreground">
									info
								</span>
							</DialogTrigger>
							<span class="text-primary material-symbols-rounded text-[32px]">
								sports_soccer
							</span>
							<h1 class="text-primary-foreground text-2xl font-extrabold leading-none">Active
								Season {seasonQuery.data.name}</h1>
							<p class="text-sm text-muted-foreground text-center">
								Ends on <span class="text-primary">{formatDate(seasonQuery.data.end_date)}</span>
							</p>
						</div>
					</Dialog>
				</Show>
			</div>
			<Tabs defaultValue="preview" class="flex flex-col relative mr-auto w-full h-full">
				<TabsList class="flex-shrink-0 w-full justify-start rounded-none border-b bg-transparent p-0 h-[56px]">
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
				<TabsContent value="matches"
										 class="pt-4 px-3 space-y-2 w-full overflow-y-scroll pb-[300px]">
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
				<TabsContent value="leaderboard"
										 class="pt-4 px-3 space-y-2 w-full overflow-y-scroll pb-[300px]">
					<Show when={leaderboardQuery.data}>
						<For each={leaderboardQuery.data}>
							{(entry) => (
								<Link class="flex items-center justify-between h-12 px-3 bg-card rounded-2xl"
											href={`/users/${entry.user.username}`} state={{ from: '/matches' }}>
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
										<Show
											when={entry.user?.favorite_team}
										>
											<img
												src={entry.user?.favorite_team?.crest_url}
												alt={entry.user?.favorite_team?.short_name}
												class="size-4 ml-1"
											/>
										</Show>
										<Show
											when={entry.user?.current_win_streak >= 3}
										>
										<span class="text-xs text-orange-500 ml-1">
											{entry.user?.current_win_streak}
										</span>
											<span class="material-symbols-rounded text-[16px] text-orange-500">
											local_fire_department
										</span>
										</Show>
									</div>
									<div class="flex items-center">
										<p class="text-sm font-medium text-muted-foreground mr-2">{entry.points} DPS</p>
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
		<div class="w-full bg-secondary p-3 rounded-2xl flex items-start justify-start flex-row space-x-1">
			<span class="material-symbols-rounded text-[24px]">
				sports_soccer
			</span>
			<div class="flex flex-col items-start justify-start space-y-2">
				<h1 class="text-2xl font-bold leading-none">{title}</h1>
				<p class="text-sm text-secondary-foreground text-center">
					{text}
				</p>
			</div>
		</div>
	)
}
