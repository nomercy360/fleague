import { createEffect, createSignal, For, Show } from 'solid-js'
import { createQuery } from '@tanstack/solid-query'
import { fetchActiveSeasons, fetchLeaderboard, fetchMatches, Season } from '~/lib/api'

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
import { useTranslations } from '~/lib/locale-context'
import { store } from '~/store'
import { useNavigate, useSearchParams } from '@solidjs/router'

export default function MatchesPage() {
	const [selectedMatch, setSelectedMatch] = createSignal({} as any)

	const [activeSeason, setActiveSeason] = createSignal<Season | null>(null)

	const [searchParams, setSearchParams] = useSearchParams()

	console.log('Search params', searchParams.tab)
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

	const closePopup = () => {
		setShowCommunityPopup(false)
		window.Telegram.WebApp.CloudStorage.setItem('fb_community_popup', 'closed')
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

	return (
		<>
			<div class="min-h-[160px] p-3 flex-col flex items-center justify-center space-y-3">
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
							{t('join_community')}
						</h1>
						<p class="text-sm text-secondary-foreground text-center">
							{t('join_community_description')}
						</p>
						<Button
							class="w-full mt-3"
							onClick={() => {
								window.Telegram.WebApp.openTelegramLink('https://t.me/match_predict_league')
							}}
						>
							{t('open_in_telegram')}
						</Button>
					</div>
				</Show>
				<Show
					when={activeSeason() && !showCommunityPopup()}
					fallback={<div class="w-full rounded-2xl bg-secondary" />}
				>
					<SeasonCard season={activeSeason()!} type={activeSeason()?.type || 'monthly'} />
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

				<TabsContent value="leaderboard" class="pt-4 px-3 space-y-2 w-full overflow-y-scroll pb-[300px]">
					<Show when={leaderboardQuery.data}>
						<For each={leaderboardQuery.data.monthly}>
							{(entry) => (
								<LeaderBoardEntry entry={entry} position={getUserPosition(entry.user_id, 'monthly')}
																	tab="leaderboard" />)}
						</For>
					</Show>
				</TabsContent>

				<TabsContent value="big-season" class="pt-4 px-3 space-y-2 w-full overflow-y-scroll pb-[300px]">
					<Show when={leaderboardQuery.data}>
						<For each={leaderboardQuery.data.football}>
							{(entry) => (<LeaderBoardEntry entry={entry} position={getUserPosition(entry.user_id, 'football')}
																						 tab="big-season" />)}
						</For>
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
}

function LeaderBoardEntry(props: LeaderBoardEntryProps) {
	return (
		<Link
			class="flex items-center justify-between h-12 px-3 bg-card rounded-2xl"
			href={`/users/${props.entry.user.username}`}
			state={{ from: `/matches?tab=${props.tab}` }}
		>
			<div class="flex items-center">
				<span class="w-4 text-center text-base font-semibold text-secondary-foreground">
					{props.position}
				</span>
				<img
					src={props.entry.user.avatar_url}
					alt="User avatar"
					class="ml-4 size-6 rounded-full object-cover"
				/>
				<p class="text-base font-semibold ml-2">
					{props.entry.user?.first_name} {props.entry.user?.last_name}
				</p>
				<Show when={props.entry.user?.favorite_team}>
					<img
						src={props.entry.user?.favorite_team?.crest_url}
						alt={props.entry.user?.favorite_team?.short_name}
						class="size-4 ml-1"
					/>
				</Show>
				<Show when={props.entry.user?.current_win_streak >= 3}>
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

type SeasonCardProps = {
	season: Season
	type: string
}

function SeasonCard(props: SeasonCardProps) {
	const { t } = useTranslations()
	return (
		<Dialog>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>{props.type === 'monthly' ? 'Monthly Seasons' : 'Big Season'}</DialogTitle>
					<DialogDescription>
						<img
							class="mb-1 mt-4 rounded-xl h-[300px] w-full object-cover"
							src="/preview.jpg"
							alt="Season Prize"
						/>
						<p class="mb-4 text-xs">
							Season {props.season.name} prize - {props.type === 'monthly' ? '"Ural" FC T-Shirt' : 'Big Trophy'}
						</p>
						<p class="text-sm">
							{props.type === 'monthly'
								? 'Join our monthly challenges to win!'
								: 'Compete for the ultimate glory!'}
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
				<h1 class="text-foreground text-2xl font-extrabold leading-none">
					{props.type === 'monthly' ? t('active_season', { name: props.season.name }) : `${t('big_season')} ${props.season.name}`}
				</h1>
				<p class="text-sm text-secondary-foreground text-center">
					<Show when={props.type === 'monthly'}>
						{t('season_ends_on', formatDate(props.season.end_date, false, store.user?.language_code))}
					</Show>
					<Show when={props.type === 'football'}>
						{t('big_season_ends_on', formatDate(props.season.end_date, false, store.user?.language_code))}
					</Show>
				</p>
			</div>
		</Dialog>
	)
}
