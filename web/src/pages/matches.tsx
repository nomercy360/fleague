import { createEffect, createSignal, For, onMount, Show } from 'solid-js'
import { createQuery } from '@tanstack/solid-query'
import { fetchMatches, Match, Prediction, saveMatchPrediction } from '~/lib/api'
import { IconMinus, IconPlus, IconStar } from '~/components/icons'
import { Button } from '~/components/ui/button'

import { Drawer, DrawerClose, DrawerContent, DrawerFooter, DrawerTrigger } from '~/components/ui/drawer'
import { cn, formatDate, timeToLocaleString } from '~/lib/utils'
import { Switch, SwitchControl, SwitchLabel, SwitchThumb } from '~/components/ui/switch'
import { queryClient } from '~/App'


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
									<DrawerTrigger
										class="h-[120px] relative grid grid-cols-3 items-center rounded-2xl max-w-md mx-auto p-2.5 pt-4 bg-card"
										onClick={() => {
											setSelectedMatch(match)
										}}
									>
										<Show when={match.prediction}>
											<div
												class="z-0 rounded-l-2xl absolute inset-y-0 left-0 w-1/2 bg-gradient-to-r from-primary to-transparent pointer-events-none"></div>
											<IconStar class="text-white absolute left-4 top-1/2 transform -translate-y-1/2 size-4"
																fill="currentColor" />
										</Show>
										<div class="z-10 flex flex-col items-center space-y-2 text-center">
											<img src={`/logos/${match.home_team.name}.png`} alt="" class="w-10" />
											<p class="max-w-20 text-xs text-foreground">{match.home_team.name}</p>
										</div>
										<div class="flex flex-col items-center text-center">
											<p class="text-xs text-muted-foreground text-center">
												{match.tournament}
											</p>
											<span class="text-lg font-bold text-center">
									{timeToLocaleString(match.match_date)}
								</span>
										</div>
										<div class="flex flex-col items-center space-y-2 text-center">
											<img src={`/logos/${match.away_team.name}.png`} alt="" class="w-10" />
											<p class="text-xs text-foreground">{match.away_team.name}</p>
										</div>
									</DrawerTrigger>)}
							</For>
						</>
					))}
				</Show>
				<FootballScoreboard match={selectedMatch()} onUpdate={onPredictionUpdate} />
			</Drawer>
		</div>
	)
}


interface ScoreboardProps {
	match: Match
	onUpdate: () => void
}


function FootballScoreboard(props: ScoreboardProps) {
	const [team1Score, setTeam1Score] = createSignal<number | null>(null)
	const [team2Score, setTeam2Score] = createSignal<number | null>(null)

	const [outcome, setOutcome] = createSignal<'home' | 'draw' | 'away' | null>(null)

	const [isExactScore, setIsExactScore] = createSignal(false)

	const increment = (setScore: (value: number) => void) => {
		setScore((prev) => prev + 1)
		setOutcome(null)
	}

	const decrement = (setScore: (value: number) => void) => {
		setScore((prev) => (prev > 0 ? prev - 1 : 0))
		setOutcome(null)
	}

	const updateOutcome = (outcome: 'home' | 'draw' | 'away') => {
		setOutcome(outcome)
		setTeam1Score(null)
		setTeam2Score(null)
	}

	const onPredictionSave = async () => {
		const prediction: Prediction = {
			match_id: props.match.id,
			predicted_home_score: null,
			predicted_away_score: null,
			predicted_outcome: null,
		}

		if (isExactScore()) {
			prediction.predicted_home_score = team1Score()
			prediction.predicted_away_score = team2Score()
		}

		if (outcome()) {
			prediction.predicted_outcome = outcome()
		}

		try {
			await saveMatchPrediction(prediction)
			props.onUpdate()
		} catch (e) {
			console.error('Failed to save prediction:', e)
		}
	}

	createEffect(() => {
		if (props.match.prediction) {
			setTeam1Score(props.match.prediction.predicted_home_score)
			setTeam2Score(props.match.prediction.predicted_away_score)
			setOutcome(props.match.prediction.predicted_outcome)
		} else {
			setTeam1Score(null)
			setTeam2Score(null)
			setOutcome(null)
		}

		setIsExactScore(!!props.match.prediction?.predicted_home_score)
	})

	return (
		<DrawerContent>
			<div class="mx-auto w-full max-w-sm">
				<div class="flex flex-col items-center gap-4">
					<Switch class="rounded-2xl p-2 flex w-full items-center justify-between space-x-6"
									checked={isExactScore()}
									onChange={setIsExactScore}>
						<SwitchLabel class="text-sm text-muted-foreground font-normal">
							Точный счёт
						</SwitchLabel>
						<SwitchControl>
							<SwitchThumb />
						</SwitchControl>
					</Switch>
					<div class="w-full justify-between flex flex-col items-start gap-2">
						<div class="flex flex-row w-full justify-between items-center h-10">
							<div class="flex flex-row items-center space-x-2">
								<img src={`/logos/${props.match.home_team.name}.png`} alt="" class="w-6" />
								<p class="mt-2 text-xs text-muted-foreground mb-2">
									{props.match.home_team.name}
								</p>
							</div>
							<Show when={isExactScore()}>
								<div class="space-x-2 flex items-center">
									<Button
										variant="outline"
										size="icon"
										onClick={() => decrement(setTeam1Score)}
										disabled={team1Score() === 0}
									>
										<IconMinus class="size-4" />
										<span class="sr-only">Decrease</span>
									</Button>
									<div class="text-lg font-bold text-foreground size-9 items-center flex justify-center">
										{team1Score() ?? '—'}
									</div>
									<Button
										variant="outline"
										size="icon"
										onClick={() => increment(setTeam1Score)}
									>
										<IconPlus class="size-4" />
										<span class="sr-only">Increase</span>
									</Button>
								</div>
							</Show>
						</div>
						<div class="flex flex-row w-full justify-between items-center h-10">
							<div class="flex flex-row items-center space-x-2">
								<img src={`/logos/${props.match.away_team.name}.png`} alt="" class="w-6" />
								<p class="mt-2 text-xs text-muted-foreground mb-2">
									{props.match.away_team.name}
								</p>
							</div>
							<Show when={isExactScore()}>
								<div class="space-x-2 flex items-center">
									<Button
										variant="outline"
										size="icon"
										onClick={() => decrement(setTeam2Score)}
										disabled={team2Score() === 0}
									>
										<IconMinus class="size-4" />
										<span class="sr-only">Decrease</span>
									</Button>
									<div class="text-lg font-bold text-foreground size-9 items-center flex justify-center">
										{team2Score() ?? '—'}
									</div>
									<Button
										variant="outline"
										size="icon"
										onClick={() => increment(setTeam2Score)}
									>
										<IconPlus class="size-4" />
										<span class="sr-only">Increase</span>
									</Button>
								</div>
							</Show>
						</div>
						<p class="text-xs text-muted-foreground mt-2">
							{formatDate(props.match.match_date, true)}
						</p>
					</div>
					<div class="h-12 w-full">
						<Show when={!isExactScore()}>
							<div class="grid grid-cols-3 w-full gap-2">
								<Button
									size="sm"
									variant="outline"
									class={cn(outcome() === 'home' && 'bg-primary text-primary-foreground')}
									onClick={() => updateOutcome('home')}
								>
									П1
								</Button>
								<Button
									size="sm"
									variant="outline"
									class={cn(outcome() === 'draw' && 'bg-primary text-primary-foreground')}
									onClick={() => updateOutcome('draw')}
								>
									Ничья
								</Button>
								<Button
									size="sm"
									variant="outline"
									class={cn(outcome() === 'away' && 'bg-primary text-primary-foreground')}
									onClick={() => updateOutcome('away')}
								>
									П2
								</Button>
							</div>
						</Show>
					</div>
				</div>
			</div>
			<DrawerFooter>
				<DrawerClose>
					<Button
						size="default"
						class="w-full"
						onClick={onPredictionSave}
					>
						Сохранить
					</Button>
				</DrawerClose>
			</DrawerFooter>
		</DrawerContent>
	)
}
