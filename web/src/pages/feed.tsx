import {
	IconChevronRight,
} from '~/components/icons'

import UserActivity from '~/components/prediction-card'
import { Link } from '~/components/link'
import { store } from '~/store'
import { Button } from '~/components/ui/button'
import { createSignal, For, onMount, Show } from 'solid-js'
import { useNavigate } from '@solidjs/router'
import { ProfileStat } from '~/pages/user'
import { useLocale } from '@kobalte/core'
import { useTranslations } from '~/lib/locale-context'

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

	const { t } = useTranslations()

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
						size="sm">
						{t('edit_profile')}
						<span
							class="material-symbols-rounded text-[20px] text-primary-foreground"
						>
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
					<ProfileStat
						icon="leaderboard"
						value={`#${store.user?.ranks.find((r) => r.season_type === 'monthly')?.position}`}
						color="#3498DB"
						label={t('rank')}
					/>
					<ProfileStat
						icon="star"
						value={store.user?.total_points}
						color="#F1C40F"
						label={t('points_earned')}
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
		</div>
	)
}
