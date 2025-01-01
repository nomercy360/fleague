import { createEffect, createSignal, onCleanup, onMount, Show } from 'solid-js'
import { useNavigate } from '@solidjs/router'
import { useMainButton } from '~/lib/useMainButton'
import { useBackButton } from '~/lib/useBackButton'
import { setIsOnboardingComplete } from '~/pages/feed'


function OnboardingPage() {
	const [currentStep, setCurrentStep] = createSignal(0)
	const navigate = useNavigate()
	const mainButton = useMainButton()
	const backButton = useBackButton()

	const steps = [
		{
			title: 'Welcome to MatchPredict!',
			description: 'Make predictions on upcoming match scores or outcomes and earn points.',
			icon: 'sports_soccer',
		},
		{
			title: 'Earn Points',
			description: 'Guess the exact score to earn 7 points or predict the outcome for 3 points. Points are locked once the match starts.',
			icon: 'score',
		},
		{
			title: 'Bonus Streaks',
			description: 'Maintain a streak of correct predictions to earn bonus points. The longer your streak, the higher your bonus!',
			icon: 'emoji_events',
		},
		{
			title: 'Invite Friends',
			description: 'Invite friends to join MatchPredict and receive 10% of their prediction points. Grow your network and your rewards!',
			icon: 'people',
		},
		{
			title: 'Climb the Leaderboard',
			description: 'Compete with others and climb the leaderboard each season. Top players will receive exciting prizes!',
			icon: 'leaderboard',
		},
		{
			title: 'Get Started!',
			description: 'Let’s dive in and start predicting matches. Good luck!',
			icon: 'thumb_up',
		},
	]

	const nextStep = () => {
		if (currentStep() < steps.length - 1) {
			window.Telegram.WebApp.HapticFeedback.selectionChanged()
			setCurrentStep(currentStep() + 1)
		} else {
			navigate('/')
		}
	}

	const prevStep = () => {
		if (currentStep() > 0) {
			window.Telegram.WebApp.HapticFeedback.selectionChanged()
			setCurrentStep(currentStep() - 1)
		}
	}

	const navigateHome = async () => {
		setIsOnboardingComplete(true)
		await window.Telegram.WebApp.CloudStorage.setItem('onboarding_complete', 'true')
		navigate('/')
	}

	onMount(() => {
		mainButton.enable('Next')
	})

	const configureButtons = (step: number) => {
		if (step === steps.length - 1) {
			mainButton.enable('Close & Start')
			mainButton.onClick(navigateHome)
			mainButton.offClick(nextStep)
		} else {
			mainButton.enable('Next')
			mainButton.onClick(nextStep)
			mainButton.offClick(navigateHome)
		}

		if (step === 0) {
			backButton.setVisible()
			backButton.onClick(navigateHome)
			backButton.offClick(prevStep)
		} else {
			backButton.onClick(prevStep)
			backButton.offClick(navigateHome)
		}
	}

	createEffect(() => {
		configureButtons(currentStep())
	})

	onCleanup(() => {
		mainButton.hide()
		mainButton.offClick(nextStep)
		mainButton.offClick(navigateHome)
		backButton.offClick(prevStep)
	})

	return (
		<div class="min-h-screen bg-gradient-to-b from-blue-500 to-indigo-600 flex items-center justify-center p-4">
			<div
				class="h-[60vh] bg-card rounded-2xl shadow-lg w-full max-w-md p-6 flex flex-col justify-between items-center">
				<div class="flex flex-col items-center justify-center">
					<div class="flex items-center justify-center w-16 h-16 bg-blue-100 rounded-full mb-6">
          <span class="material-symbols-rounded text-4xl text-blue-500">
            {steps[currentStep()].icon}
          </span>
					</div>
					<h2 class="text-2xl font-bold text-center mb-4">{steps[currentStep()].title}</h2>
					<p class="text-center text-secondary-foreground mb-6">{steps[currentStep()].description}</p>
				</div>
				<div class="flex mt-4 space-x-2">
					{steps.map((_, index) => (
						<div
							class={`w-3 h-3 rounded-full ${
								index === currentStep()
									? 'bg-primary'
									: 'bg-muted'
							}`}
						></div>
					))}
				</div>
			</div>
		</div>
	)
}

export default OnboardingPage
