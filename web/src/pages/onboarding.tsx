import { createEffect, createSignal, onCleanup, onMount, Show } from 'solid-js'
import { useNavigate } from '@solidjs/router'
import { useMainButton } from '~/lib/useMainButton'
import { useBackButton } from '~/lib/useBackButton'
import { setIsOnboardingComplete } from '~/pages/feed'
import { useTranslations } from '~/lib/locale-context'


function OnboardingPage() {
	const [currentStep, setCurrentStep] = createSignal(0)
	const navigate = useNavigate()
	const mainButton = useMainButton()
	const backButton = useBackButton()

	const { t } = useTranslations()

	const steps = [
		{
			icon: 'sports_soccer',
			color: '#3498DB',
		},
		{
			icon: 'star',
			color: '#F1C40F',
		},
		{
			icon: 'local_fire_department',
			color: '#E74C3C',
		},
		{
			icon: 'people',
			color: '#2ECC71',
		},
		{
			icon: 'leaderboard',
			color: '#3498DB',
		},
		{
			icon: 'thumb_up',
			color: '#2ECC71',
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
		<div class="min-h-screen bg-gradient-to-b from-background to-primary flex items-center justify-center p-4">
			<div
				class="h-[60vh] bg-card rounded-2xl w-full max-w-md p-6 flex flex-col justify-between items-center">
				<div class="flex flex-col items-center justify-center">
					<div class="flex items-center justify-center w-16 h-16 bg-blue-100 rounded-full mb-6">
          <span class="icon-fill material-symbols-rounded text-4xl" style={{ color: steps[currentStep()].color }}>
            {steps[currentStep()].icon}
          </span>
					</div>
					<h2 class="text-2xl font-bold text-center mb-4">{t(`onboarding.${currentStep()}.title`)}</h2>
					<p class="text-center text-secondary-foreground mb-6">{t(`onboarding.${currentStep()}.description`)}</p>
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
