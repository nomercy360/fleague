import { useMainButton } from '~/lib/useMainButton'

import { createEffect, createSignal, onCleanup, onMount, Show } from 'solid-js'
import { fetchMatchByID, PredictionRequest, saveMatchPrediction } from '~/lib/api'
import { createQuery } from '@tanstack/solid-query'
import { useNavigate, useParams } from '@solidjs/router'
import { showToast } from '~/components/ui/toast'
import { useTranslations } from '~/lib/locale-context'
import { formatDate } from '~/lib/utils'
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

	const matchQuery = createQuery(() => ({
		queryKey: ['matches', params.id],
		queryFn: () => fetchMatchByID(params.id),
	}))

	const handleSubmit = async () => {
		setIsSubmitting(true)
		const prediction: PredictionRequest = {
			match_id: matchQuery.data.id,
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
		}
	}

	const mainButton = useMainButton()

	createEffect(() => {
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
					<h2 class="text-lg font-bold text-foreground">{matchQuery.data.tournament}</h2>
					<p
						class="text-muted-foreground text-sm">{formatDate(matchQuery.data.match_date, true, store.user?.language_code)}</p>
				</div>

				<div class="w-full p-4 flex items-center justify-between">
					<div class="flex flex-col items-center w-2/5">
						<img
							src={matchQuery.data.home_team.crest_url || '/placeholder.svg'}
							alt={matchQuery.data.home_team.name}
							class="size-14 object-contain"
						/>
						<p class="text-sm mt-2 text-center font-semibold">{matchQuery.data.home_team.short_name}</p>
					</div>

					<div class="text-center text-secondary-foreground font-bold">
						{t('vs')}
					</div>

					<div class="flex flex-col items-center w-2/5">
						<img
							src={matchQuery.data.away_team.crest_url || '/placeholder.svg'}
							alt={matchQuery.data.away_team.name}
							class="size-14 object-contain"
						/>
						<p class="text-sm mt-2 text-center font-semibold">{matchQuery.data.away_team.short_name}</p>
					</div>
				</div>

				<div class="w-full px-4 pb-4">
					<div class="grid grid-cols-3 text-center justify-between text-sm text-muted-foreground mb-2">
						<span>{t('home')} {matchQuery.data.home_odds}</span>
						<span>{t('draw')} {matchQuery.data.draw_odds}</span>
						<span>{t('away')} {matchQuery.data.away_odds}</span>
					</div>
				</div>

				<div class="p-4 border-t w-full border-secondary">
					<h3 class="text-lg font-semibold mb-3">
						{t('your_prediction')}
					</h3>

					<div class="flex mb-4 rounded-lg p-1">
						<button
							type="button"
							onClick={() => setPredictionType('outcome')}
							class={`font-medium text-xs flex-1 py-3 h-9 flex items-center justify-center text-center rounded-xl transition ${
								predictionType() === 'outcome' ? 'bg-primary text-primary-foreground' : 'text-muted-foreground'
							}`}
						>
							{t('match_outcome')}
						</button>
						<button
							type="button"
							onClick={() => setPredictionType('score')}
							class={`font-medium text-xs flex-1 py-3 h-9 flex items-center justify-center text-center rounded-xl transition ${
								predictionType() === 'score' ? 'bg-primary text-primary-foreground' : 'text-muted-foreground'
							}`}
						>
							{t('exact_score')}
						</button>
					</div>

					<Show when={predictionType() === 'outcome'}>
						<div class="grid grid-cols-3 gap-3 mb-4">
							<button
								type="button"
								onClick={() => setSelectedOutcome('home')}
								class={`p-3 rounded-lg border ${
									selectedOutcome() === 'home' ? 'bg-blue-100 border-blue-500' : 'border-gray-300'
								}`}
							>
								<div class="text-center">
									<p class="font-semibold">{matchQuery.data.home_team.abbreviation}</p>
									<p class="text-sm">{t('win')}</p>
								</div>
							</button>

							<button
								type="button"
								onClick={() => setSelectedOutcome('draw')}
								class={`p-3 rounded-lg border ${
									selectedOutcome() === 'draw' ? 'bg-blue-100 border-blue-500' : 'border-gray-300'
								}`}
							>
								<div class="text-center">
									<p class="font-semibold">X</p>
									<p class="text-sm">{t('draw')}</p>
								</div>
							</button>

							<button
								type="button"
								onClick={() => setSelectedOutcome('away')}
								class={`p-3 rounded-lg border ${
									selectedOutcome() === 'away' ? 'bg-blue-100 border-blue-500' : 'border-gray-300'
								}`}
							>
								<div class="text-center">
									<p class="font-semibold">{matchQuery.data.away_team.abbreviation}</p>
									<p class="text-sm">{t('win')}</p>
								</div>
							</button>
						</div>
					</Show>

					<Show when={predictionType() === 'score'}>
						<div class="flex items-center justify-center gap-4 mb-4">
							<div class="w-1/3">
								<label class="block text-sm text-center mb-1">{matchQuery.data.home_team.abbreviation}</label>
								<input
									type="number"
									min="0"
									max="20"
									value={homeScore()}
									onInput={(e) => setHomeScore(e.target.value)}
									class="w-full p-3 text-center text-xl border border-gray-300 rounded-lg"
								/>
							</div>

							<div class="text-xl font-bold">:</div>

							<div class="w-1/3">
								<label class="block text-sm text-center mb-1">{matchQuery.data.away_team.abbreviation}</label>
								<input
									type="number"
									min="0"
									max="20"
									value={awayScore()}
									onInput={(e) => setAwayScore(e.target.value)}
									class="w-full p-3 text-center text-xl border border-gray-300 rounded-lg"
								/>
							</div>
						</div>
					</Show>
				</div>
			</Show>
		</div>
	)
}

export default MatchPage

