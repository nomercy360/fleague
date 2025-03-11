import { useMainButton } from '~/lib/useMainButton'
import { createEffect, createSignal, Match, onCleanup, onMount, Show, Switch } from 'solid-js'
import { fetchMatchByID, MatchResponse, PredictionRequest, saveMatchPrediction } from '~/lib/api'
import { createQuery } from '@tanstack/solid-query'
import { useNavigate, useParams } from '@solidjs/router'
import { showToast } from '~/components/ui/toast'
import { useTranslations } from '~/lib/locale-context'
import { cn, formatDate } from '~/lib/utils'
import { store } from '~/store'
import { queryClient } from '~/App'

function MatchPage() {
	const [predictionType, setPredictionType] = createSignal('outcome')
	const [selectedOutcome, setSelectedOutcome] = createSignal('')
	const [homeScore, setHomeScore] = createSignal('')
	const [awayScore, setAwayScore] = createSignal('')
	const [isSubmitting, setIsSubmitting] = createSignal(false)

	const params = useParams()
	const { t } = useTranslations()
	const navigate = useNavigate()

	const matchQuery = createQuery<MatchResponse>(() => ({
		queryKey: ['matches', params.id],
		queryFn: () => fetchMatchByID(params.id),
	}))

	const handleSubmit = async () => {
		setIsSubmitting(true)
		const prediction: PredictionRequest = {
			match_id: matchQuery.data!.id,
			predicted_home_score: null,
			predicted_away_score: null,
			predicted_outcome: null,
		}

		if (predictionType() === 'score') {
			prediction.predicted_home_score = Number(homeScore())
			prediction.predicted_away_score = Number(awayScore())
		}

		if (predictionType() === 'outcome') {
			prediction.predicted_outcome = selectedOutcome()
		}

		const { error } = await saveMatchPrediction(prediction)
		if (!error) {
			setIsSubmitting(false)
			await matchQuery.refetch()
			await queryClient.invalidateQueries({ queryKey: ['predictions'] })
			showToast({ variant: 'success', title: t('prediction_saved'), duration: 3000 })
			navigate('/')
		} else {
			setIsSubmitting(false)
		}
	}

	const mainButton = useMainButton()

	createEffect(() => {
		if (matchQuery.data?.status !== 'scheduled') {
			mainButton.hide()
			return
		}

		if (predictionType() === 'outcome' && !selectedOutcome()) {
			mainButton.disable(t('save_prediction'))
		} else if (predictionType() === 'score' && (homeScore() === '' || awayScore() === '')) {
			mainButton.disable(t('save_prediction'))
		} else if (isSubmitting()) {
			mainButton.disable(t('saving'))
		} else {
			mainButton.enable(t('save_prediction'))
		}
	})

	createEffect(() => {
		if (matchQuery.data && matchQuery.data.prediction) {
			if (matchQuery.data.prediction.predicted_outcome) {
				setSelectedOutcome(matchQuery.data.prediction.predicted_outcome)
			} else if (matchQuery.data.prediction.predicted_home_score != null) {
				setHomeScore(matchQuery.data.prediction.predicted_home_score.toString())
				setAwayScore(matchQuery.data.prediction.predicted_away_score.toString())
				setPredictionType('score')
			}
		}
	})

	onMount(() => {
		mainButton.onClick(handleSubmit)
	})

	onCleanup(() => {
		mainButton.hide()
		mainButton.offClick(handleSubmit)
	})

	return (
		<div class="w-full min-h-screen flex flex-col items-center justify-start">
			<Show when={matchQuery.isSuccess}>
				<div class="w-full bg-secondary p-3 text-center">
					<h2 class="text-lg font-bold text-foreground">{matchQuery.data?.tournament}</h2>
					<p class="text-muted-foreground text-sm">
						{formatDate(matchQuery.data!.match_date, true, store.user?.language_code)}
					</p>
				</div>

				<div class="w-full p-4 flex items-center justify-between">
					<div class="flex flex-col items-center w-2/5">
						<img
							src={matchQuery.data?.home_team.crest_url || '/placeholder.svg'}
							alt={matchQuery.data?.home_team.name}
							class="size-14 object-contain"
						/>
						<p class="text-sm mt-2 text-center font-semibold">{matchQuery.data?.home_team.short_name}</p>
					</div>

					<div class="text-center text-secondary-foreground font-bold">
						<Show when={matchQuery.data?.status === 'scheduled'}>
							{t('vs')}
						</Show>
						<Show when={matchQuery.data?.status === 'completed'}>
							{matchQuery.data?.home_score} - {matchQuery.data?.away_score}
						</Show>
					</div>

					<div class="flex flex-col items-center w-2/5">
						<img
							src={matchQuery.data?.away_team.crest_url || '/placeholder.svg'}
							alt={matchQuery.data?.away_team.name}
							class="size-14 object-contain"
						/>
						<p class="text-sm mt-2 text-center font-semibold">{matchQuery.data?.away_team.short_name}</p>
					</div>
				</div>

				<div class="w-full px-4 pb-4">
					<div class="grid grid-cols-3 text-center justify-between text-sm text-muted-foreground mb-2">
						<span>{t('home')} {matchQuery.data?.home_odds}</span>
						<span>{t('draw')} {matchQuery.data?.draw_odds}</span>
						<span>{t('away')} {matchQuery.data?.away_odds}</span>
					</div>
				</div>

				<div class="px-1 py-4 border-t w-full border-secondary">
					<h3 class="px-2 text-lg font-semibold mb-2">{t('your_prediction')}</h3>

					<Show when={matchQuery.data?.status != 'scheduled' && matchQuery.data?.prediction?.completed_at}>
						<div
							class={cn(
								'flex flex-row justify-between items-center w-full mb-3 p-3 bg-secondary rounded-lg',
								matchQuery.data?.prediction.points_awarded ? 'shadow-green-500 bg-green-100' : 'shadow-red-400 bg-red-100',
							)}
						>
							<div class="flex flex-col items-start">
								<Switch>
									<Match when={matchQuery.data?.prediction.predicted_outcome === 'home'}>
										<p class="font-bold text-lg">{matchQuery.data?.home_team.abbreviation}</p>
										<p class="text-sm text-secondary-foreground">{t('win')}</p>
									</Match>
									<Match when={matchQuery.data?.prediction.predicted_outcome === 'draw'}>
										<p class="font-bold text-lg">X</p>
										<p class="text-sm text-secondary-foreground">{t('draw')}</p>
									</Match>
									<Match when={matchQuery.data?.prediction.predicted_outcome === 'away'}>
										<p class="font-bold text-lg">{matchQuery.data?.away_team.abbreviation}</p>
										<p class="text-sm text-secondary-foreground">{t('win')}</p>
									</Match>
									<Match when={matchQuery.data?.prediction.predicted_home_score != null}>
										<p class="font-bold text-lg">
											{matchQuery.data?.prediction.predicted_home_score} - {matchQuery.data?.prediction.predicted_away_score}
										</p>
										<p class="text-sm text-secondary-foreground">{t('exact_score')}</p>
									</Match>
								</Switch>
							</div>
							<Show when={matchQuery.data?.prediction.points_awarded}>
								<div class="flex flex-col items-end">
									<p class="font-bold text-lg">+{matchQuery.data?.prediction.points_awarded}</p>
									<p class="text-xs text-secondary-foreground">{t('points')}</p>
								</div>
							</Show>
						</div>
					</Show>

					<Show when={matchQuery.data?.status === 'scheduled'}>
						<div class="flex mb-4 rounded-lg p-1">
							<button
								type="button"
								onClick={() => setPredictionType('outcome')}
								class={`font-medium text-xs flex-1 py-3 h-9 flex items-center justify-center text-center rounded-xl transition ${
									predictionType() === 'outcome' ? 'bg-primary text-primary-foreground' : 'text-secondary-foreground'
								}`}
							>
								{t('match_outcome')}
							</button>
							<button
								type="button"
								onClick={() => setPredictionType('score')}
								class={`font-medium text-xs flex-1 py-3 h-9 flex items-center justify-center text-center rounded-xl transition ${
									predictionType() === 'score' ? 'bg-primary text-primary-foreground' : 'text-secondary-foreground'
								}`}
							>
								{t('exact_score')}
							</button>
						</div>
						<div class="mb-1 ml-2">
							<Show when={predictionType() === 'outcome'}>
								<span class="text-xs font-medium">
									{t('prediction_cost_and_reward', { cost: 20, points: 7 })}
								</span>
							</Show>
							<Show when={predictionType() === 'score'}>
								<span class="text-xs font-medium">
									{t('prediction_cost_and_reward', { cost: 30, points: 10 })}
								</span>
							</Show>
						</div>
						<Show when={predictionType() === 'outcome'}>
							<div class="grid grid-cols-3 gap-3 mb-4">
								<button
									type="button"
									onClick={() => setSelectedOutcome('home')}
									class={`p-3 rounded-lg border ${
										selectedOutcome() === 'home' ? 'bg-secondary' : 'border'
									}`}
								>
									<div class="text-center">
										<p class="font-semibold">{matchQuery.data?.home_team.abbreviation}</p>
										<p class="text-xs text-secondary-foreground">{t('win')}</p>
									</div>
								</button>

								<button
									type="button"
									onClick={() => setSelectedOutcome('draw')}
									class={`p-3 rounded-lg border ${
										selectedOutcome() === 'draw' ? 'bg-secondary' : 'border'
									}`}
								>
									<div class="text-center">
										<p class="font-semibold">X</p>
										<p class="text-xs text-secondary-foreground">{t('draw')}</p>
									</div>
								</button>

								<button
									type="button"
									onClick={() => setSelectedOutcome('away')}
									class={`p-3 rounded-lg border ${
										selectedOutcome() === 'away' ? 'bg-secondary' : 'border'
									}`}
								>
									<div class="text-center">
										<p class="font-semibold">{matchQuery.data?.away_team.abbreviation}</p>
										<p class="text-xs text-secondary-foreground">{t('win')}</p>
									</div>
								</button>
							</div>
						</Show>

						<Show when={predictionType() === 'score'}>
							<div class="flex items-center justify-center gap-4 mb-4">
								<div class="w-1/3">
									<label class="block text-sm text-center mb-1">{matchQuery.data?.home_team.abbreviation}</label>
									<input
										type="number"
										min="0"
										max="20"
										value={homeScore()}
										onInput={(e) => setHomeScore(e.target.value)}
										class="w-full p-3 text-center text-xl bg-secondary rounded-lg"
									/>
								</div>

								<div class="text-xl font-bold">:</div>

								<div class="w-1/3">
									<label class="block text-sm text-center mb-1">{matchQuery.data?.away_team.abbreviation}</label>
									<input
										type="number"
										min="0"
										max="20"
										value={awayScore()}
										onInput={(e) => setAwayScore(e.target.value)}
										class="w-full p-3 text-center text-xl bg-secondary rounded-lg"
									/>
								</div>
							</div>
						</Show>
					</Show>
				</div>
			</Show>
		</div>
	)
}

export default MatchPage
