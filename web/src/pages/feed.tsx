import UserActivity from '~/components/prediction-card'
import { Link } from '~/components/link'
import { store } from '~/store'
import { Button } from '~/components/ui/button'
import { createEffect, createSignal, For, onCleanup, onMount, Show } from 'solid-js'
import { useNavigate } from '@solidjs/router'
import { ProfileStat } from '~/pages/user'
import { useTranslations } from '~/lib/locale-context'
import { useMainButton } from '~/lib/useMainButton'

export default function FeedPage() {
	const navigate = useNavigate()
	const [isOnboardingComplete, setIsOnboardingComplete] = createSignal(true)
	const { t } = useTranslations()

	function shareProfileURL() {
		const url =
			'https://t.me/share/url?' +
			new URLSearchParams({
				url: 'https://t.me/footbon_bot/app?startapp=u_' + store.user?.username,
			}).toString() +
			`&text=Check out ${store.user?.first_name}'s profile`

		window.Telegram.WebApp.openTelegramLink(url)
	}

	const updateOnboardingComplete = (err: unknown, value: unknown) => {
		const isComplete = value === 'true'
		setIsOnboardingComplete(isComplete)
	}

	onMount(() => {
		// window.Telegram.WebApp.CloudStorage.removeItem('onboarding_complete')
		window.Telegram.WebApp.CloudStorage.getItem(
			'onboarding_complete',
			updateOnboardingComplete,
		)
	})

	const onboardingComplete = () => {
		window.Telegram.WebApp.CloudStorage.setItem('onboarding_complete', 'true')
		setIsOnboardingComplete(true)
	}

	const buySubscription = () => {
		// Открываем ссылку на покупку подписки через Telegram
		navigate('/subscribe') // Предполагаемый маршрут для покупки подписки
	}

	return (
		<div class="h-full overflow-y-scroll bg-background text-foreground pb-[120px]">
			<div class="relative w-full bg-card rounded-b-[10%] px-4 pt-6 pb-8 mb-8 flex flex-col items-center">
				<div class="flex flex-row justify-between items-center w-full">
					<Button
						onClick={shareProfileURL}
						size="sm"
						variant="secondary"
					>
						<span class="material-symbols-rounded text-[16px] text-secondary-foreground">
							ios_share
						</span>
						{t('share')}
					</Button>
					<Button
						href="/edit-profile"
						as={Link}
						class="gap-0"
						size="sm"
					>
						{t('edit_profile')}
						<span class="material-symbols-rounded text-[20px] text-primary-foreground">
							chevron_right
						</span>
					</Button>
				</div>
				<img
					src={store.user?.avatar_url}
					alt="User avatar"
					class="size-24 rounded-full object-cover"
				/>
				<div class="text-lg font-semibold mt-2 flex flex-row items-center">
					<span>{store.user?.first_name}</span>
					<Show when={store.user?.favorite_team}>
						<img
							src={store.user?.favorite_team?.crest_url}
							alt={store.user?.favorite_team?.short_name}
							class="size-4 ml-1"
						/>
					</Show>
					<Show when={store.user?.current_win_streak}>
						<span class="text-xs text-orange-500 ml-1">
							{store.user?.current_win_streak}
						</span>
						<span class="-ml-0.5 material-symbols-rounded text-[16px] text-orange-500">
							local_fire_department
						</span>
					</Show>
				</div>
				<p class="text-sm font-medium text-muted-foreground">@{store.user?.username}</p>
				<Show when={store.user?.badges}>
					<div class="mt-3 flex flex-row flex-wrap gap-2 items-center justify-center">
						<For each={store.user?.badges}>
							{(badge) => (
								<div class="bg-secondary rounded-2xl h-7 px-2 flex items-center gap-1">
									<span style={{ color: badge.color }}
												class="material-symbols-rounded text-[16px] text-primary-foreground">
										{badge.icon}
									</span>
									<span class="text-xs text-muted-foreground">{badge.name}</span>
								</div>
							)}
						</For>
					</div>
				</Show>
				<div class="grid grid-cols-2 gap-2 mt-6 w-full px-2">
					<ProfileStat
						icon="check_circle"
						value={store.user?.correct_predictions}
						color="#2ECC71"
						label={t('correct')}
					/>
					<Show when={store.user?.ranks.find((r) => r.season_type === 'monthly')?.position || 0 > 0}>
						<ProfileStat
							icon="leaderboard"
							value={`#${store.user?.ranks.find((r) => r.season_type === 'monthly')?.position}`}
							color="#3498DB"
							label={t('rank')}
						/>
					</Show>
					<ProfileStat
						icon="target"
						value={`${store.user?.prediction_accuracy}%`}
						color="#F1C40F"
						label={t('accuracy')}
					/>
					<Show when={store.user?.longest_win_streak || 0 > 3}>
						<ProfileStat
							icon="emoji_events"
							value={store.user?.longest_win_streak}
							color="#FFC107"
							label={t('max_streak')}
						/>
					</Show>
				</div>
			</div>
			<UserActivity />
			<Show when={!isOnboardingComplete()}>
				<OnboardingPage onComplete={() => onboardingComplete()} />
			</Show>
		</div>
	)
}

type OnboardingPageProps = {
	onComplete: () => void
}

function OnboardingPage(props: OnboardingPageProps) {
	const [step, setStep] = createSignal(0)
	const { t } = useTranslations()

	const steps = [
		{
			title: t('welcome_title'), // "Добро пожаловать в игру!" / "Welcome to the game!"
			description: t('welcome_description'), // "Угадывай результаты футбольных матчей и соревнуйся с другими за призы." / "Predict football match results and compete for prizes."
		},
		{
			title: t('how_to_predict_title'), // "Как делать прогнозы?" / "How to make predictions?"
			description: t('how_to_predict_description'), // "С подпиской делай прогнозы на точный счёт или исход матча и получай очки!" / "With a subscription, predict exact scores or match outcomes and earn points!"
		},
		{
			title: t('points_title'), // "Очки лидерборда" / "Leaderboard points"
			description: t('points_description'), // "Чем больше очков за верные прогнозы, тем выше ты в рейтинге!" / "The more points you earn from correct predictions, the higher your rank!"
		},
		{
			title: t('subscription_title'), // "Получи подписку" / "Get a subscription"
			description: t('subscription_description'), // "Покупай подписку за Telegram Stars и участвуй в игре каждый месяц." / "Buy a subscription with Telegram Stars and join the game every month."
		},
		{
			title: t('win_prizes_title'), // "Выигрывай призы" / "Win prizes"
			description: t('win_prizes_description'), // "В конце месяца лидер получает футбольную форму — отправим почтой!" / "At the end of the month, the leader wins a football jersey — shipped by mail!"
		},
	]

	const nextStep = () => {
		if (step() < steps.length - 1) setStep(step() + 1)
		else window.Telegram.WebApp.close() // Закрываем онбординг
	}

	const mainButton = useMainButton()

	onMount(() => {
		mainButton.enable(t('next')) // "Далее" / "Next"
	})

	createEffect(() => {
		if (step() === steps.length - 1) {
			mainButton.offClick(nextStep)
			mainButton.enable(t('start')).onClick(props.onComplete) // "Начать" / "Start"
		} else {
			mainButton.offClick(props.onComplete)
			mainButton.enable(t('next')).onClick(nextStep) // "Далее" / "Next"
		}
	})

	onCleanup(() => {
		mainButton.hide()
		mainButton.offClick(nextStep)
		mainButton.offClick(props.onComplete)
	})

	return (
		<div class="absolute top-0 z-50 backdrop-blur-lg h-screen w-full mx-auto flex flex-col justify-between p-4">
			<div class="flex-1 flex flex-col items-center justify-center text-center">
				<h1 class="text-2xl font-bold mb-4">{steps[step()].title}</h1>
				<p class="text-base px-2">{steps[step()].description}</p>
			</div>
			<div class="flex justify-center mb-6">
				{steps.map((_, index) => (
					<div
						class={`w-2 h-2 rounded-full mx-1 ${
							index === step() ? 'bg-blue-500' : 'bg-gray-300'
						}`}
					/>
				))}
			</div>
		</div>
	)
}
