import { createEffect, createSignal, For, Show } from 'solid-js'
import { fetchMatchStats, MatchResponse, PredictionRequest, PredictionResponse, saveMatchPrediction } from '~/lib/api'
import { DrawerClose, DrawerContent, DrawerFooter } from '~/components/ui/drawer'
import { Button } from '~/components/ui/button'
import { IconMinus, IconPlus } from '~/components/icons'
import { cn, formatDate } from '~/lib/utils'
import { useTranslations } from '~/lib/locale-context'
import { createQuery } from '@tanstack/solid-query'


interface ScoreboardProps {
	match: MatchResponse
	prediction?: PredictionResponse
	onUpdate: () => void
}

export default function FootballScoreboard(props: ScoreboardProps) {
	const [team1Score, setTeam1Score] = createSignal<number | null>(null)
	const [team2Score, setTeam2Score] = createSignal<number | null>(null)

	const [outcome, setOutcome] = createSignal<'home' | 'draw' | 'away' | null>(null)

	const [isExactScore, setIsExactScore] = createSignal(false)

	const increment = (setScore: (value: number) => void) => {
		window.Telegram.WebApp.HapticFeedback.selectionChanged()
		setScore((prev) => prev + 1)
		setOutcome(null)
	}

	const decrement = (setScore: (value: number) => void) => {
		window.Telegram.WebApp.HapticFeedback.selectionChanged()
		setScore((prev) => (prev > 0 ? prev - 1 : 0))
		setOutcome(null)
	}

	const updateOutcome = (newValue: 'home' | 'draw' | 'away') => {
		if (outcome() === newValue) {
			setOutcome(null)
			return
		}
		window.Telegram.WebApp.HapticFeedback.selectionChanged()
		setOutcome(newValue)
		setTeam1Score(null)
		setTeam2Score(null)
		setIsExactScore(false)
	}

	const statsQuery = createQuery<MatchResponse>(() => ({
		queryKey: ['match_stats', props.match.id],
		queryFn: () => fetchMatchStats(props.match.id),
		enabled: !!props.match.id,
	}))

	const onPredictionSave = async () => {
		const prediction: PredictionRequest = {
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

		const { error } = await saveMatchPrediction(prediction)
		if (!error) {
			props.onUpdate()
		}
	}

	createEffect(() => {
		if (props.prediction) {
			setTeam1Score(props.prediction.predicted_home_score)
			setTeam2Score(props.prediction.predicted_away_score)
			setOutcome(props.prediction.predicted_outcome)
		} else {
			setTeam1Score(null)
			setTeam2Score(null)
			setOutcome(null)
		}

		setIsExactScore(!!props.prediction?.predicted_home_score)
	})

	const updateSwitch = (value: boolean) => {
		setIsExactScore(value)
		setTeam1Score(null)
		setTeam2Score(null)
		setOutcome(null)
	}

	const { t } = useTranslations()

	return (
		<DrawerContent class="pb-3 bg-card">
			<div class="mx-auto w-full px-4">
				<div class="flex flex-col items-center gap-4">
					<div class="w-full justify-between flex flex-col items-start gap-2">
						<div class="flex flex-row w-full justify-between items-center h-10">
							<div class="flex flex-row w-full justify-between items-center">
								<div class="flex flex-row items-center space-x-2">
									<img src={props.match.home_team.crest_url} alt="" class="w-6" />
									<p class="mt-2 text-base mb-2">
										{props.match.home_team.short_name}
									</p>
								</div>
								{/*<WinLoseStats results={statsQuery.data?.home_team_results ?? []} />*/}
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
							<div class="flex flex-row w-full justify-between items-center">
								<div class="flex flex-row items-center space-x-2">
									<img src={props.match.away_team.crest_url} alt="" class="w-6" />
									<p class="mt-2 text-base mb-2">
										{props.match.away_team.short_name}
									</p>
								</div>
								{/*<WinLoseStats results={statsQuery.data?.away_team_results ?? []} />*/}
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
					<div class="w-full gap-2 flex flex-col mb-4">
						<div class="w-full">
							<div class="grid grid-cols-3 w-full gap-2">
								<Button
									size="sm"
									variant="outline"
									class={cn(outcome() === 'home' && 'bg-primary text-primary-foreground')}
									onClick={() => updateOutcome('home')}
								>
									<span>{t('win_1')}</span>{props.match.home_odds}
								</Button>
								<Button
									size="sm"
									variant="outline"
									class={cn(outcome() === 'draw' && 'bg-primary text-primary-foreground')}
									onClick={() => updateOutcome('draw')}
								>
									<span>{t('draw')}</span>{props.match.draw_odds}
								</Button>
								<Button
									size="sm"
									variant="outline"
									class={cn(outcome() === 'away' && 'bg-primary text-primary-foreground')}
									onClick={() => updateOutcome('away')}
								>
									<span>{t('win_2')}</span> {props.match.away_odds}
								</Button>
							</div>
						</div>
						<Button
							size="sm"
							class={cn(isExactScore() && 'bg-primary text-primary-foreground')}
							variant="outline"
							onClick={() => updateSwitch(!isExactScore())}
						>
							Exact Score
						</Button>
					</div>
				</div>
			</div>
			<DrawerFooter>
				<DrawerClose>
					<Button
						size="default"
						class="w-full bg-accent text-accent-foreground"
						disabled={(team1Score() == null || team2Score() == null) && outcome() == null}
						onClick={onPredictionSave}
					>
						Save & Close
					</Button>
				</DrawerClose>
			</DrawerFooter>
		</DrawerContent>
	)
}

function WinLoseStats(props: { results: string[] }) {
	return (
		<div class="flex flex-row items-center justify-start gap-1">
			<For each={props.results}>
				{(result) => (
					<div
						class={cn('size-6 flex items-center justify-center rounded bg-primary-foreground text-primary-foreground text-xs font-medium', {
							'bg-green-500': result === 'W',
							'bg-red-500': result === 'L',
						})}
					>
						{result}
					</div>
				)}
			</For>
		</div>
	)
}
