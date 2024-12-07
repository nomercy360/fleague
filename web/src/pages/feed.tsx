import {
	createSignal, For,
	onMount, Show,
} from 'solid-js'
import { createQuery } from '@tanstack/solid-query'
import { fetchMatches } from '~/lib/api'
import { IconChevronDown, IconPlot, IconMinus, IconPlus, IconTrophy, IconSparkles } from '~/components/icons'
import { Button } from '~/components/ui/button'
import {
	Drawer, DrawerClose,
	DrawerContent,
	DrawerFooter,
	DrawerTrigger,
} from '~/components/ui/drawer'

import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '~/components/ui/collapsible'
import PredictionCard from '~/components/prediction-card'
import { formatDate } from '~/lib/utils'


export default function FeedPage() {
	const query = createQuery(() => ({
		queryKey: ['matches'],
		queryFn: () => fetchMatches(),
	}))

	onMount(() => {
		window.Telegram.WebApp.disableClosingConfirmation()
		window.Telegram.WebApp.disableVerticalSwipes()
		// window.Telegram.WebApp.CloudStorage.removeItem('profilePopup')
		// window.Telegram.WebApp.CloudStorage.removeItem('communityPopup')
		// window.Telegram.WebApp.CloudStorage.removeItem('rewardsPopup')
	})

	const dummyUsers = [
		{
			title: 'Сергей бестов',
			subtitle: 'Спортивный аналитик',
			image: '/avatars/sergey.png',
			score: 300,
		},
		{
			title: 'Гор Е',
			subtitle: 'Болельщик Челси',
			image: '/avatars/gor.jpg',
			score: 100,
		},
		{
			title: 'Максим К',
			subtitle: 'Болельщик ЦСКА',
			image: '/avatars/maksim.jpg',
			score: 50,
		},
		{
			title: 'John Doe',
			subtitle: 'Football fan',
			image: '/avatars/user.svg',
			score: 10,
		},
		{
			title: 'Максим К',
			subtitle: 'Болельщик ЦСКА',
			image: '/avatars/maksim.jpg',
			score: 50,
		},
	]

	const [homeTeam, setHomeTeam] = createSignal('')
	const [awayTeam, setAwayTeam] = createSignal('')

	return (
		<div class="bg-background text-foreground min-h-[110vh] pt-12">
			<div class="flex flex-row items-center space-x-4 border-b p-2.5 fixed top-0 left-0 right-0 bg-background z-50">
				<img
					src={window.Telegram.WebApp.initDataUnsafe.user.photo_url}
					alt=""
					class="w-10 h-10 rounded-full object-cover"
				/>
				<div class="space-x-1.5 flex flex-row flex-nowrap">
					<Button
						variant="secondary"
						size="sm"
					>
						<IconTrophy class="size-5" />
						Таблица лидеров
					</Button>
					<Button
						variant="secondary"
						size="sm"
					>
						<IconPlot class="size-5" />
						Матчи
					</Button>
					<Button
						variant="secondary"
						size="sm"
					>
						<IconSparkles class="size-5" />
						Прогнозы
					</Button>
				</div>
			</div>
			<Collapsible defaultOpen={true} class="p-1.5">
				<CollapsibleTrigger
					class="w-full mt-6 bg-blue-100 rounded-xl px-3 h-12 flex flex-row items-center justify-between">
					<div class="space-x-1 flex flex-row items-center">
						<IconTrophy class="size-5" />
						<p class="font-semibold text-sm">
							Мои прогнозы
						</p>
					</div>
					<IconChevronDown class="size-5 text-muted-foreground" />
				</CollapsibleTrigger>
				<CollapsibleContent class="space-y-2 mt-2 rounded-t-xl">
					<PredictionCard />
				</CollapsibleContent>
			</Collapsible>
			<Show when={query.data}>
				<Drawer>
					<Collapsible defaultOpen={true}>
						<div class="sticky top-14 z-50 px-1.5 pb-2 bg-background">
							<CollapsibleTrigger
								class="w-full mt-6 bg-blue-100 rounded-xl px-3 h-12 flex flex-row items-center justify-between">
								<div class="space-x-1 flex flex-row items-center">
									<IconPlot class="size-5" />
									<p class="font-semibold text-sm">
										Ближайшие матчи
									</p>
								</div>
								<IconChevronDown class="size-5 text-muted-foreground" />
							</CollapsibleTrigger>
						</div>
						<CollapsibleContent class="space-y-2 overflow-y-scroll h-screen rounded-t-xl p-1.5">
							<For
								each={query.data}
								fallback={<div>Loading...</div>}
							>
								{match => (
									<div class="rounded-xl max-w-md mx-auto p-3 bg-secondary flex flex-col justify-between">
										<div class="grid grid-cols-2 gap-6">
											<div class="space-y-1">
												<div class="flex items-center space-x-1">
													<img src={`/logos/uefa.png`} alt="" class="w-4" />
													<p class="text-xs">UEFA Champions League</p>
												</div>
												<div class="flex items-center space-x-1">
													<img src={`/logos/${match.home_team.name}.png`} alt="" class="w-4" />
													<p class="text-xs font-bold">{match.home_team.name}</p>
												</div>
												<div class="flex items-center space-x-1">
													<img src={`/logos/${match.away_team.name}.png`} alt="" class="w-4" />
													<p class="text-xs font-bold">{match.away_team.name}</p>
												</div>
												<p class="text-xs text-muted-foreground">{formatDate(match.match_date)}</p>
											</div>
											<div class="flex items-center justify-end space-x-1">
												<DrawerTrigger
													as={Button<'button'>}
													variant="default"
													size="sm"
													onClick={() => {
														setHomeTeam(match.home_team.name)
														setAwayTeam(match.away_team.name)
													}}
												>
													Прогноз
												</DrawerTrigger>
											</div>
										</div>
									</div>
								)}
							</For>
						</CollapsibleContent>
					</Collapsible>
					<DrawerContent>
						<div class="mx-auto w-full max-w-sm">
							<FootballScoreboard
								home_team={homeTeam()}
								away_team={awayTeam()}
							/>
							<DrawerFooter>
								<Button>
									Сохранить
								</Button>
								<DrawerClose as={Button<'button'>} variant="outline">
									Закрыть
								</DrawerClose>
							</DrawerFooter>
						</div>
					</DrawerContent>
				</Drawer>
			</Show>
		</div>
	)
}

interface UserCardProps {
	title: string
	subtitle: string
	image: string
	score: number
}

function UserLeaderboardCard(props: UserCardProps) {
	return (
		<div class="p-3 flex flex-col bg-secondary w-full">
			<div class="flex flex-row items-center justify-between">
				<div class="flex flex-row items-center justify-center space-x-2">
					<img src={props.image} alt="" class="size-6 rounded-full object-cover" />
					<div>
						<p class="text-xs font-bold">{props.title}</p>
						<p class="text-xs">{props.subtitle}</p>
					</div>
				</div>
				<p class="text-accent-foreground text-sm font-bold">{props.score}</p>
			</div>
		</div>
	)
}

interface ScoreboardProps {
	home_team: string
	away_team: string
}

function FootballScoreboard(props: ScoreboardProps) {
	const [team1Score, setTeam1Score] = createSignal(0)
	const [team2Score, setTeam2Score] = createSignal(0)

	const increment = (setScore: (value: number) => void) => setScore((prev) => prev + 1)
	const decrement = (setScore: (value: number) => void) =>
		setScore((prev) => (prev > 0 ? prev - 1 : 0))
	const resetScores = () => {
		setTeam1Score(0)
		setTeam2Score(0)
	}

	return (
		<div class="flex flex-col items-center pt-4">
			<div class="flex flex-row items-center gap-6 mb-6">
				<div class="flex flex-col items-center w-24">
					<div class="flex flex-col items-center justify-center text-center">
						<img src={`/logos/${props.home_team}.png`} alt="" class="w-8" />
						<p class="text-xs font-semibold mb-2">
							{props.home_team}
						</p>
					</div>
					<div class="text-4xl font-bold text-blue-600 mb-4">{team1Score()}</div>
					<div class="flex space-x-2">
						<Button
							variant="outline"
							size="icon"
							class="size-8 shrink-0 rounded-full"
							onClick={() => decrement(setTeam1Score)}
							disabled={team1Score() === 0}
						>
							<IconMinus class="size-4" />
							<span class="sr-only">Decrease</span>
						</Button>
						<Button
							variant="outline"
							size="icon"
							class="size-8 shrink-0 rounded-full"
							onClick={() => increment(setTeam1Score)}
						>
							<IconPlus class="size-4" />
							<span class="sr-only">Increase</span>
						</Button>
					</div>
				</div>
				<div class="flex flex-col items-center text-center w-24">
					<p class="text-sm font-bold text-accent-foreground">
						завтра
					</p>
					<p class="text-xs text-muted-foreground">20:00</p>
				</div>
				<div class="flex flex-col items-center w-24">
					<div class="flex flex-col items-center justify-center text-center">
						<img src={`/logos/${props.away_team}.png`} alt="" class="w-8" />
						<p class="text-xs font-semibold mb-2">
							{props.away_team}
						</p>
					</div>
					<div class="text-4xl font-bold text-red-600 mb-4">{team2Score()}</div>
					<div class="flex space-x-2">
						<Button
							variant="outline"
							size="icon"
							class="size-8 shrink-0 rounded-full"
							onClick={() => decrement(setTeam2Score)}
							disabled={team2Score() === 0}
						>
							<IconMinus class="size-4" />
							<span class="sr-only">Decrease</span>
						</Button>
						<Button
							variant="outline"
							size="icon"
							class="size-8 shrink-0 rounded-full"
							onClick={() => increment(setTeam2Score)}
						>
							<IconPlus class="size-4" />
							<span class="sr-only">Increase</span>
						</Button>
					</div>
				</div>
			</div>
		</div>
	)
}
