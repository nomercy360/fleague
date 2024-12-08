import {
	createSignal,
	onMount,
} from 'solid-js'
import { createQuery } from '@tanstack/solid-query'
import { fetchMatches } from '~/lib/api'
import {
	IconMinus,
	IconPlus,
	IconUsers,
	IconActivity, IconCalendar, IconChevronRight,
} from '~/components/icons'
import { Button } from '~/components/ui/button'

import PredictionCard from '~/components/prediction-card'
import { Link } from '~/components/link'
import { useMainButton } from '~/lib/useMainButton'
import { cn } from '~/lib/utils'
import { useLocation } from '@solidjs/router'


export default function FeedPage() {
	const query = createQuery(() => ({
		queryKey: ['matches'],
		queryFn: () => fetchMatches(),
	}))

	const mainButton = useMainButton()

	onMount(() => {
		window.Telegram.WebApp.disableClosingConfirmation()
		window.Telegram.WebApp.disableVerticalSwipes()
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

	const location = useLocation()

	return (
		<div class="bg-background text-foreground pb-24">
			<div class="w-full bg-card rounded-b-[10%] px-4 pt-6 pb-8 mb-8 flex flex-col items-center">
				<img
					src={window.Telegram.WebApp.initDataUnsafe.user.photo_url}
					alt=""
					class="size-24 rounded-full object-cover"
				/>
				<p class="text-lg font-semibold mt-2">
					{window.Telegram.WebApp.initDataUnsafe.user.first_name}
				</p>
				<Link href="/" class="text-muted-foreground flex flex-row items-center">
					<p class="text-sm">
						@{window.Telegram.WebApp.initDataUnsafe.user.username}
					</p>
				</Link>
				<div class="flex flex-row items-center justify-center space-x-2 mt-4">
					<div class="flex flex-col rounded-2xl py-3 px-4 bg-secondary text-white w-[100px] self-stretch">
						<span class="text-2xl font-semibold">10</span>
						<span class="text-xs">Correct predictions</span>
					</div>
					<div class="flex flex-col rounded-2xl py-3 px-4 bg-secondary text-white w-[100px] self-stretch">
						<span class="text-2xl font-semibold">#123</span>
						<span class="text-xs">Global ranking</span>
					</div>
					<div class="flex flex-col rounded-2xl py-3 px-4 bg-secondary text-white w-[100px] self-stretch">
						<span class="text-2xl font-semibold">300</span>
						<span class="text-xs">Points earned</span>
					</div>
				</div>
			</div>
			<div class="px-4 mb-6">
				<Link class="flex flex-row h-14 justify-between items-center rounded-2xl p-3 bg-secondary"
							href="/leaderboard">
					<p class="text-sm font-semibold">
						Make a prediction{' '}
						<span class="text-muted-foreground font-normal">12 matches available</span>
					</p>
					<IconChevronRight class="size-6" />
				</Link>
			</div>

			<PredictionCard />
			<div
				class="flex flex-row items-center space-x-4 border-t px-2.5 h-20 fixed bottom-0 left-0 right-0 bg-background z-50">
				<div class="grid grid-cols-3 w-full">
					<Link
						href="/"
						class={cn('flex items-center flex-col h-full text-sm gap-1', {
							'text-primary': location.pathname === '/',
						})}
					>
						<IconActivity class="size-6" />
						Activity
					</Link>
					<Link
						href="/"
						class={cn('flex items-center flex-col h-full text-sm gap-1', {
							'text-primary': location.pathname === '/matches',
						})}
					>
						<IconCalendar class="size-6" />
						Matches
					</Link>
					<Link
						href="/"
						class={cn('flex items-center flex-col h-full text-sm gap-1', {
							'text-primary': location.pathname === '/friends',
						})}
					>
						<IconUsers class="size-6" />
						Friends
					</Link>
				</div>
			</div>
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

// function Matches() {
// 	return (
// 		<Drawer>
// 			<Collapsible defaultOpen={true}>
// 				<div class="sticky top-14 z-50 px-1.5 pb-2 bg-background">
// 					<CollapsibleTrigger
// 						class="w-full mt-6 bg-blue-100 rounded-xl px-3 h-12 flex flex-row items-center justify-between">
// 						<div class="space-x-1 flex flex-row items-center">
// 							<IconPlot class="size-5" />
// 							<p class="font-semibold text-sm">
// 								Ближайшие матчи
// 							</p>
// 						</div>
// 						<IconChevronDown class="size-5 text-muted-foreground" />
// 					</CollapsibleTrigger>
// 				</div>
// 				<CollapsibleContent class="space-y-2 overflow-y-scroll h-screen rounded-t-xl p-1.5">
// 					<For
// 						each={query.data}
// 						fallback={<div>Loading...</div>}
// 					>
// 						{match => (
// 							<div class="rounded-xl max-w-md mx-auto p-3 bg-secondary flex flex-col justify-between">
// 								<div class="grid grid-cols-2 gap-6">
// 									<div class="space-y-1">
// 										<div class="flex items-center space-x-1">
// 											<img src={`/logos/uefa.png`} alt="" class="w-4" />
// 											<p class="text-xs">UEFA Champions League</p>
// 										</div>
// 										<div class="flex items-center space-x-1">
// 											<img src={`/logos/${match.home_team.name}.png`} alt="" class="w-4" />
// 											<p class="text-xs font-bold">{match.home_team.name}</p>
// 										</div>
// 										<div class="flex items-center space-x-1">
// 											<img src={`/logos/${match.away_team.name}.png`} alt="" class="w-4" />
// 											<p class="text-xs font-bold">{match.away_team.name}</p>
// 										</div>
// 										<p class="text-xs text-muted-foreground">{formatDate(match.match_date)}</p>
// 									</div>
// 									<div class="flex items-center justify-end space-x-1">
// 										<DrawerTrigger
// 											as={Button<'button'>}
// 											variant="default"
// 											size="sm"
// 											onClick={() => {
// 												setHomeTeam(match.home_team.name)
// 												setAwayTeam(match.away_team.name)
// 											}}
// 										>
// 											Прогноз
// 										</DrawerTrigger>
// 									</div>
// 								</div>
// 							</div>
// 						)}
// 					</For>
// 				</CollapsibleContent>
// 			</Collapsible>
// 			<DrawerContent>
// 				<div class="mx-auto w-full max-w-sm">
// 					<FootballScoreboard
// 						home_team={homeTeam()}
// 						away_team={awayTeam()}
// 					/>
// 					<DrawerFooter>
// 						<Button>
// 							Сохранить
// 						</Button>
// 						<DrawerClose as={Button<'button'>} variant="outline">
// 							Закрыть
// 						</DrawerClose>
// 					</DrawerFooter>
// 				</div>
// 			</DrawerContent>
// 		</Drawer>
// 	)
// }
