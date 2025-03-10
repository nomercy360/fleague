import { Link } from '~/components/link'
import { cn } from '~/lib/utils'
import { useLocation } from '@solidjs/router'
import { createSignal, onMount, Show } from 'solid-js'
import { Button } from '~/components/ui/button'
import { sendFeedback } from '~/lib/api'

export default function NavigationTabs(props: any) {
	const location = useLocation()

	const tabs = [
		{ href: '/', icon: 'dashboard', activePath: '/' },
		{ href: '/matches', icon: 'sports_soccer', activePath: '/matches' },
		{ href: '/friends', icon: 'groups', activePath: '/friends' },
	]

	return (
		<div class="h-screen bg-background text-foreground">
			<PredictionDialog />
			<div
				class="flex flex-row items-start border-t h-[100px] fixed bottom-0 w-full bg-background z-50 transform -translate-x-1/2 left-1/2"
			>
				<div class="px-2.5 py-4 flex flex-row w-full gap-10 items-center justify-center">
					{tabs.map(({ href, icon, activePath }) => (
						<Link
							href={href}
							class={cn('size-10 rounded-full p-2 flex items-center flex-col h-full text-sm gap-1', {
								'bg-blue-500 text-primary-foreground': location.pathname === activePath,
							})}
						>
						<span class="material-symbols-rounded icon-fill text-[24px]">
							{icon}
						</span>
						</Link>
					))}
				</div>
			</div>
			{props.children}
		</div>
	)
}


const PredictionDialog = () => {
	const [selectedOption, setSelectedOption] = createSignal<string | null>(null)
	const [showSurvey, setShowSurvey] = createSignal(false)

	const updateSurveyComplete = (err: unknown, value: unknown) => {
		const isComplete = value === 'true'
		setShowSurvey(!isComplete)
	}

	onMount(() => {
		// window.Telegram.WebApp.CloudStorage.removeItem('fl_survey_complete')
		window.Telegram.WebApp.CloudStorage.getItem(
			'fl_survey_complete',
			updateSurveyComplete,
		)
	})

	const handleSubmit = async (e: any) => {
		e.preventDefault()
		if (!selectedOption()) return

		try {
			const { data, error } = await sendFeedback({
				feature: 'prediction_prizes',
				preference: selectedOption(),
			})

			if (data) {
				onClose()
			}
		} catch (error) {
			console.error('Error submitting feedback:', error)
		}
	}

	const onClose = () => {
		window.Telegram.WebApp.CloudStorage.setItem('fl_survey_complete', 'true')
		setShowSurvey(false)
	}

	return (
		<Show when={showSurvey()}>
			<div class="px-3 fixed inset-0 backdrop-blur-sm flex items-center justify-center z-50">
				<div class="relative bg-background rounded-lg pr-4 pl-6 pt-5 pb-6 w-full max-w-md">
					<div class="pb-4 flex flex-row items-center justify-between w-full">
						<h2 class="text-xl font-bold">Новая фича!</h2>
						<button class="flex items-center justify-center rounded-sm"
										onClick={() => onClose()}>
						<span
							class="material-symbols-rounded text-[24px] text-muted-foreground"
						>
							close
						</span>
						</button>
					</div>
					<p class="mb-6">
						Друзья, мы хотим добавить в игру денежные призы! Делайте взнос, соревнуйтесь в прогнозах и выигрывайте
						реальные деньги. Что думаете?
					</p>

					<div class="space-y-3 mb-8">
						<label class="flex items-center gap-3 cursor-pointer">
							<input
								type="radio"
								name="option"
								value="yes"
								checked={selectedOption() === 'yes'}
								onChange={(e) => setSelectedOption(e.target.value)}
								class="hidden peer"
							/>
							<div
								class="w-5 h-5 border-2 rounded flex items-center justify-center peer-checked:border-primary peer-checked:bg-primary">
								<svg
									class={selectedOption() === 'yes' ? 'block' : 'hidden'}
									width="12"
									height="12"
									viewBox="0 0 12 12"
									fill="none"
								>
									<path
										d="M2 6L5 9L10 3"
										stroke="white"
										stroke-width="2"
										stroke-linecap="round"
										stroke-linejoin="round"
									/>
								</svg>
							</div>
							<span>Да, готов участвовать!</span>
						</label>

						<label class="flex items-center gap-3 cursor-pointer">
							<input
								type="radio"
								name="option"
								value="no"
								checked={selectedOption() === 'no'}
								onChange={(e) => setSelectedOption(e.target.value)}
								class="hidden peer"
							/>
							<div
								class="w-5 h-5 border-2 rounded flex items-center justify-center peer-checked:border-primary peer-checked:bg-primary">
								<svg
									class={selectedOption() === 'no' ? 'block' : 'hidden'}
									width="12"
									height="12"
									viewBox="0 0 12 12"
									fill="none"
								>
									<path
										d="M2 6L5 9L10 3"
										stroke="white"
										stroke-width="2"
										stroke-linecap="round"
										stroke-linejoin="round"
									/>
								</svg>
							</div>
							<span>Нет, мне это не интересно</span>
						</label>
					</div>

					<div class="flex justify-center">
						<Button
							onClick={handleSubmit}
							disabled={!selectedOption()}
						>
							Отправить
						</Button>
					</div>
				</div>
			</div>
		</Show>
	)
}

