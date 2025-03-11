import { createEffect, createMemo, createSignal, For, Show } from 'solid-js'
import { createQuery } from '@tanstack/solid-query'
import { fetchActiveSeasons, fetchLeaderboard, fetchMatches, Season } from '~/lib/api'

import { Drawer, DrawerTrigger } from '~/components/ui/drawer'
import { cn, formatDate } from '~/lib/utils'
import { queryClient } from '~/App'
import MatchCard from '~/components/match-card'
import FootballScoreboard from '~/components/score-board'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '~/components/ui/tabs'
import { Link } from '~/components/link'
import { Button } from '~/components/ui/button'

import { useTranslations } from '~/lib/locale-context'
import { getUserLeaderboardPoints, getUserLeaderboardPosition, store } from '~/store'
import { useNavigate, useSearchParams } from '@solidjs/router'

export default function MatchesPage() {
	const [selectedMatch, setSelectedMatch] = createSignal({} as any)

	const [activeSeason, setActiveSeason] = createSignal<Season | null>(null)

	const [searchParams, setSearchParams] = useSearchParams()

	const navigate = useNavigate()

	const [activeTab, setActiveTab] = createSignal(searchParams.tab || 'matches')

	const query = createQuery(() => ({
		queryKey: ['matches'],
		queryFn: () => fetchMatches(),
	}))

	const leaderboardQuery = createQuery(() => ({
		queryKey: ['leaderboard'],
		queryFn: () => fetchLeaderboard(),
	}))

	const seasonQuery = createQuery<Season[]>(() => ({
		queryKey: ['season'],
		queryFn: () => fetchActiveSeasons(),
	}))

	const onPredictionUpdate = () => {
		query.refetch()
		queryClient.invalidateQueries({ queryKey: ['predictions'] })
	}

	const getUserPosition = (id: number, type: 'monthly' | 'football') => {
		const data = type === 'monthly' ? leaderboardQuery.data.monthly : leaderboardQuery.data.football
		const idx = data.findIndex((entry: any) => entry.user_id === id)

		if (idx === 0) return '🥇'
		if (idx === 1) return '🥈'
		if (idx === 2) return '🥉'
		return idx + 1
	}

	const { t } = useTranslations()

	createEffect(() => {
		if (seasonQuery.data) {
			if (activeTab() == 'big-season') {
				const season = seasonQuery.data.find((season) => season.type === 'football')
				setActiveSeason(season!)
			} else {
				const season = seasonQuery.data.find((season) => season.type === 'monthly')
				setActiveSeason(season!)
			}
		}
	})

	const handleTabChange = (tab: string) => {
		setActiveTab(tab)
		setSearchParams({ tab }) // Update query parameter
	}

	const isUserInLeaderboard = (userId: string, type: 'monthly' | 'football') => {
		const leaderboard = type === 'monthly' ? leaderboardQuery.data.monthly : leaderboardQuery.data.football
		return leaderboard.some((entry: any) => entry.user_id === userId)
	}

	return (
		<>
			<div class="min-h-[160px] p-3 flex-col flex items-center justify-center space-y-3">
				<Show
					when={activeSeason()}
					fallback={<div class="w-full rounded-2xl bg-secondary" />}
				>
					<div class="shadow-md flex flex-row items-start justify-between w-full bg-secondary rounded-xl p-3">
						<div class="flex flex-col gap-1 mr-6">
							<p class="text-base font-bold">
								{t('contest.win_tshirt')}
							</p>
							<p class="text-xs text-secondary-foreground">
								{t('contest.results_announcement', { date: formatDate(activeSeason()!.end_date, false, store.user?.language_code) })}
							</p>
							<Button
								class="mt-2 h-9 text-sm"
								onClick={() => {
									window.Telegram.WebApp.openTelegramLink('https://t.me/mpl_footbal_analyst')
								}}
							>
								{t('contest.join_channel')}
							</Button>
						</div>
						<img src="/football-tshirt.png" class="shrink-0 w-24 rounded-lg" />
					</div>
				</Show>
			</div>
			<Tabs
				defaultValue={activeTab() as any}
				onChange={handleTabChange}
				class="flex flex-col relative mr-auto w-full h-full"
			>
				<TabsList class="flex-shrink-0 w-full justify-start rounded-none border-b bg-transparent p-0 h-[56px]">
					<TabsTrigger
						value="matches"
						class="relative h-[56px] rounded-none border-b-2 border-b-transparent bg-transparent px-4 font-semibold text-muted-foreground shadow-none transition-none data-[selected]:border-b-primary data-[selected]:text-foreground data-[selected]:shadow-none"
					>
						{t('matches')}
					</TabsTrigger>
					<TabsTrigger
						value="leaderboard"
						class="relative h-[56px] rounded-none border-b-2 border-b-transparent bg-transparent px-4 font-semibold text-muted-foreground shadow-none transition-none data-[selected]:border-b-primary data-[selected]:text-foreground data-[selected]:shadow-none"
					>
						{t('leaderboard')}
					</TabsTrigger>
					<TabsTrigger
						value="big-season"
						class="relative h-[56px] rounded-none border-b-2 border-b-transparent bg-transparent px-4 font-semibold text-muted-foreground shadow-none transition-none data-[selected]:border-b-primary data-[selected]:text-foreground data-[selected]:shadow-none"
					>
						{t('big_season')}
					</TabsTrigger>
				</TabsList>

				<TabsContent value="matches" class="pt-4 px-3 space-y-2 w-full overflow-y-scroll pb-[300px]">
					<Drawer>
						<Show when={!query.isLoading}>
							{Object.entries(query.data).map(([date, matches]) => (
								<>
									<p class="px-2 text-base font-semibold">
										{formatDate(date)}
									</p>
									<For each={matches as any}>
										{match => (<MatchCard match={match} prediction={match.prediction} />)}
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

				<TabsContent value="leaderboard" class="pt-4 px-3 space-y-2 w-full overflow-y-scroll pb-[300px]">
					<Show when={leaderboardQuery.data}>
						<For each={leaderboardQuery.data.monthly}>
							{(entry) => (
								<LeaderBoardEntry entry={entry} position={getUserPosition(entry.user_id, 'monthly')}
																	tab="leaderboard" />)}
						</For>
						<Show when={!isUserInLeaderboard(store.user.id, 'monthly')}>
							<LeaderBoardEntry entry={{ user: store.user, points: getUserLeaderboardPoints('monthly') }}
																position={getUserLeaderboardPosition('monthly')}
																tab="leaderboard"
																type="virtual" />

						</Show>
					</Show>
				</TabsContent>

				<TabsContent value="big-season" class="pt-4 px-3 space-y-2 w-full overflow-y-scroll pb-[300px]">
					<Show when={leaderboardQuery.data}>
						<For each={leaderboardQuery.data.football}>
							{(entry) => (<LeaderBoardEntry entry={entry} position={getUserPosition(entry.user_id, 'football')}
																						 tab="big-season" />)}
						</For>
						<Show when={!isUserInLeaderboard(store.user.id, 'football')}>
							<LeaderBoardEntry entry={{ user: store.user, points: getUserLeaderboardPoints('football') }}
																position={getUserLeaderboardPosition('football')}
																tab="big-season"
																type="virtual" />
						</Show>
					</Show>
				</TabsContent>
			</Tabs>
		</>
	)
}


type LeaderBoardEntryProps = {
	entry: any
	position: number
	tab: string
	type?: string
}

function LeaderBoardEntry(props: LeaderBoardEntryProps) {
	return (
		<Link
			class="flex items-center justify-between h-12 px-3 bg-card rounded-2xl"
			href={`/users/${props.entry.user.username}`}
			state={{ from: `/matches?tab=${props.tab}` }}
		>
			<div class="flex items-center">
				<span
					class={cn('w-4 text-center text-base font-semibold', props.type === 'virtual' ? 'text-primary-foreground' : 'text-muted-foreground')}>
					{props.position}
				</span>
				<img
					src={props.entry.user.avatar_url}
					alt="User avatar"
					class={cn('ml-3 rounded-full object-cover', props.type === 'virtual' ? 'size-7' : 'size-6')}
				/>
				<p class="text-base font-semibold ml-2">
					{props.entry.user?.first_name} {props.entry.user?.last_name}
				</p>
				<Show when={props.entry.user?.badges?.length > 0}>
					<div class="ml-1 flex items-center justify-center size-6 rounded-full bg-secondary"
							 style={{ color: props.entry.user?.badges[0].color }}>
						<span class="material-symbols-rounded text-[16px]">
							{props.entry.user?.badges[0].icon}
						</span>
					</div>
				</Show>
				<Show when={!props.entry.user?.badges?.length && props.entry.user?.favorite_team}>
					<img
						src={props.entry.user?.favorite_team?.crest_url}
						alt={props.entry.user?.favorite_team?.short_name}
						class="size-4 ml-1"
					/>
				</Show>
				<Show when={!props.entry.user?.badges?.length && props.entry.user?.current_win_streak >= 3}>
					<span class="text-xs text-orange-500 ml-1">
						{props.entry.user?.current_win_streak}
					</span>
					<span class="material-symbols-rounded text-[16px] text-orange-500">
						local_fire_department
					</span>
				</Show>
			</div>
			<div class="flex items-center">
				<p class="text-sm font-medium text-secondary-foreground mr-0.5">
					{props.entry.points}
				</p>
				<span class="text-[12px] material-symbols-rounded text-yellow-200 icon-fill">star</span>
			</div>
		</Link>
	)
}

