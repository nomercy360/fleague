import {
	IconChevronRight,
} from '~/components/icons'

import UserActivity from '~/components/prediction-card'
import { Link } from '~/components/link'
import { store } from '~/store'
import { Button } from '~/components/ui/button'
import { createSignal, Match, onMount, Show, Switch } from 'solid-js'
import { useNavigate } from '@solidjs/router'

export const [isOnboardingComplete, setIsOnboardingComplete] = createSignal(false)

export default function FeedPage() {
	const navigate = useNavigate()

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
		if (!isComplete && !isOnboardingComplete()) {
			navigate('/onboarding')
		}
	}

	onMount(() => {
		// window.Telegram.WebApp.CloudStorage.removeItem('onboarding_complete')

		window.Telegram.WebApp.CloudStorage.getItem(
			'onboarding_complete',
			updateOnboardingComplete,
		)
	})

	return (
		<div class="h-full overflow-y-scroll bg-background text-foreground pb-[120px]">
			<div class="relative w-full bg-card rounded-b-[10%] px-4 pt-6 pb-8 mb-8 flex flex-col items-center">
				<div class="flex flex-row justify-between items-center w-full">
					<div class="flex flex-row items-center justify-start gap-1">
						<Button
							onClick={shareProfileURL}
							size="sm"
							variant="secondary"
						>
							<span class="material-symbols-rounded text-[16px] text-secondary-foreground">
								ios_share
							</span>
							Share
						</Button>
						<Button
							onClick={shareProfileURL}
							size="sm"
							variant="ghost"
							as={Link}
							href="/onboarding"
						>
							<span class="material-symbols-rounded text-[16px] text-secondary-foreground">
								info
							</span>
						</Button>
					</div>
					<Button
						href="/edit-profile"
						as={Link}
						size="sm">
						Edit profile
						<IconChevronRight class="size-6" />
					</Button>
				</div>
				<img
					src={store.user?.avatar_url}
					alt="User avatar"
					class="size-24 rounded-full object-cover"
				/>
				<div class="text-lg font-semibold mt-2 flex flex-row items-center">
					<span>{store.user?.first_name}</span>
					<Show
						when={store.user?.favorite_team}
					>
						<img
							src={store.user?.favorite_team?.crest_url}
							alt={store.user?.favorite_team?.short_name}
							class="size-4 ml-1"
						/>
					</Show>
					<Show
						when={store.user?.current_win_streak}
					>
						<span class="text-xs text-orange-500 ml-1">
							{store.user?.current_win_streak}
						</span>
						<span class="material-symbols-rounded text-[16px] text-[#E74C3C]">
							local_fire_department
						</span>
					</Show>
				</div>
				<Link href="/" class="text-muted-foreground flex flex-row items-center">
					<p class="text-sm">
						@{store.user?.username}
					</p>
				</Link>
				<div class="w-full flex flex-row items-center justify-center space-x-2 mt-4">
					<Show when={store.user?.global_rank && store.user?.correct_predictions > 0} fallback={<GoToMatchesLink />}>
						<div class="flex flex-col rounded-2xl py-3 px-4 bg-secondary text-card-foreground w-[100px] self-stretch">
							<span class="text-2xl font-semibold">{store.user?.correct_predictions}</span>
							<span class="text-xs text-muted-foreground">Correct predictions</span>
						</div>
						<div
							class="flex flex-col rounded-2xl py-3 px-4 bg-secondary text-card-foreground w-[100px] self-stretch">
							<span class="text-2xl font-semibold">
								#{store.user?.global_rank}
							</span>
							<span class="text-xs text-muted-foreground">Global ranking</span>
						</div>
						<div
							class="flex flex-col rounded-2xl py-3 px-4 bg-secondary text-card-foreground w-[100px] self-stretch">
							<span class="text-2xl font-semibold">{store.user?.total_points}</span>
							<span class="text-xs text-muted-foreground">Points earned</span>
						</div>
					</Show>
				</div>
			</div>
			<UserActivity />
		</div>
	)
}


const GoToMatchesLink = () => {
	return (
		<Link class="bg-secondary w-full flex flex-row h-14 justify-between items-center rounded-2xl p-3 space-x-6"
					href="/matches">
			<div>
				<p class="text-sm font-semibold">
					Make your first prediction
				</p>
				<p class="text-xs text-muted-foreground font-normal">12 matches available</p>
			</div>
			<IconChevronRight class="size-6" />
		</Link>
	)
}
