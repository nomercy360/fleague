import {
	createSignal, For,
	onMount, Show,
} from 'solid-js'
import { createQuery } from '@tanstack/solid-query'
import { fetchMatches } from '~/lib/api'
import { IconCrown, IconFootball, IconMinus, IconPlus } from '~/components/icons'
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '~/components/ui/accordion'
import { Button } from '~/components/ui/button'
import {
	Drawer, DrawerClose,
	DrawerContent,
	DrawerDescription, DrawerFooter,
	DrawerHeader,
	DrawerTitle,
	DrawerTrigger,
} from '~/components/ui/drawer'

function formatDate(dateString: string) {
	const date = new Date(dateString)
	const options = {
		weekday: 'short',
		hour: '2-digit',
		minute: '2-digit',
		day: 'numeric',
		month: 'long',
	}
	return date.toLocaleDateString('en-GB', options as any)
}

export default function FeedPage() {
	const query = createQuery(() => ({
		queryKey: ['matches'],
		queryFn: () => fetchMatches(),
	}))

	onMount(() => {
		window.Telegram.WebApp.disableClosingConfirmation()
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
		<div class="bg-secondary text-foreground min-h-screen pt-2 p-1.5">
			<div class="bg-background rounded-t-lg px-3 h-12 flex flex-row items-center space-x-1">
				<IconCrown class="size-6 text-primary" />
				<p class="font-bold text-sm">
					Таблица лидеров
				</p>
			</div>
			<div class="rounded-b-lg flex flex-col space-y-1 overflow-y-scroll h-64">
				<For each={dummyUsers}>
					{user => (
						<UserLeaderboardCard
							title={user.title}
							subtitle={user.subtitle}
							image={user.image}
							score={user.score}
						/>
					)}
				</For>
			</div>
			<Show when={query.isLoading}>
				<div>Loading...</div>
			</Show>
			<Show when={query.data}>
				<Drawer>
					<div class="mt-4 bg-background rounded-t-lg px-3 h-12 flex flex-row items-center space-x-1">
						<IconFootball class="size-7 text-primary" />
						<p class="font-bold text-sm">Upcoming matches</p>
					</div>
					<div class="w-full h-px bg-primary"></div>
					<Accordion multiple={false} collapsible class="w-full" defaultValue={['item-1']}>
						<AccordionItem value="item-1">
							<AccordionTrigger class="h-10 bg-background">
								<div class="flex items-center space-x-1 px-3">
									<img src={`/logos/uefa.png`} alt="" class="w-4" />
									<p class="text-xs">UEFA Champions League</p>
								</div>
							</AccordionTrigger>
							<AccordionContent>
								<div class="space-y-1 mt-1 overflow-y-scroll h-screen">
									<For
										each={query.data}
										fallback={<div>Loading...</div>}
									>
										{match => (
											<div class="max-w-md mx-auto p-3 bg-background flex flex-col justify-between">
												<div class="grid grid-cols-2 gap-6">
													<div class="space-y-1">
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
															class="h-full w-14"
															as={Button<'button'>}
															variant="secondary"
															onClick={() => {
																setHomeTeam(match.home_team.name)
																setAwayTeam(match.away_team.name)
															}}
														>
															П1
														</DrawerTrigger>
														<DrawerTrigger class="h-full w-14" as={Button<'button'>} variant="secondary">
															Х
														</DrawerTrigger>
														<DrawerTrigger class="h-full w-14" as={Button<'button'>} variant="secondary">
															П2
														</DrawerTrigger>
													</div>
												</div>
											</div>
										)}
									</For>
								</div>
							</AccordionContent>
						</AccordionItem>
					</Accordion>
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
		<div class="p-3 flex flex-col bg-background w-full">
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
