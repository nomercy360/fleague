import { store } from '~/store'
import { createQuery } from '@tanstack/solid-query'
import { fetchUserInfo } from '~/lib/api'
import { useParams } from '@solidjs/router'
import { For, Show } from 'solid-js'
import MatchCard from '~/components/match-card'
import { Button } from '~/components/ui/button'
import {
	IconShare,
} from '~/components/icons'
import { useTranslations } from '~/lib/locale-context'

type ProfileStatProps = {
	icon: string
	value?: any
	label: string
	color: string
}

export function ProfileStat(props: ProfileStatProps) {
	return (
		<div class="space-x-2 flex-grow bg-background flex flex-row items-start text-center rounded-2xl py-2 px-3">
			<span class="py-1 material-symbols-rounded text-[20px]" style={{ color: props.color }}>
				{props.icon}
			</span>
			<div class="flex flex-col items-start justify-start">
				<span class="font-extrabold text-lg">{props.value}</span>
				<span class="text-xs text-muted-foreground text-nowrap">
					{props.label}
				</span>
			</div>
		</div>
	)
}

//

export default function UserProfilePage() {
	const params = useParams()
	const username = params.username

	function shareProfileURL() {
		const url =
			'https://t.me/share/url?' +
			new URLSearchParams({
				url: 'https://t.me/footbon_bot/app?startapp=u_' + store.user?.username,
			}).toString() +
			`&text=Check out ${store.user?.first_name}'s profile`

		window.Telegram.WebApp.openTelegramLink(url)
	}

	const userInfoQuery = createQuery(() => ({
		queryKey: ['user', username],
		queryFn: () => fetchUserInfo(username),
	}))

	const { t } = useTranslations()

	return (
		<div class="bg-background text-foreground pb-24 h-screen overflow-y-scroll">
			<Show when={userInfoQuery.isLoading}>
				<div class="flex flex-col items-center justify-center h-full">
					<div class="loader" />
				</div>
			</Show>
			<Show when={userInfoQuery.data}>
				<div class="w-full bg-card rounded-b-[10%] px-4 pt-6 pb-8 mb-8 flex flex-col items-center">
					<div class="flex flex-row justify-between items-center w-full">
						<Button onClick={shareProfileURL} size="sm" variant="secondary">
							<IconShare class="size-6 mr-1" />
							{t('share')}
						</Button>
					</div>
					<img
						src={userInfoQuery.data.user.avatar_url}
						alt="User avatar"
						class="w-24 h-24 rounded-full object-cover mt-4"
					/>
					<div class="text-lg font-semibold mt-2 flex flex-row items-center">
						<span>{userInfoQuery.data.user.first_name}</span>
						<Show
							when={userInfoQuery.data.user.favorite_team}
						>
							<img
								src={userInfoQuery.data.user.favorite_team.crest_url}
								alt={userInfoQuery.data.user.favorite_team.name}
								class="size-4 ml-1"
							/>
						</Show>
						<Show
							when={userInfoQuery.data.user.current_win_streak >= 3}
						>
						<span class="text-xs text-orange-500 ml-1">
							{userInfoQuery.data.user.current_win_streak}
						</span>
							<span class="-ml-0.5 material-symbols-rounded text-[16px] text-orange-500">
							local_fire_department
						</span>
						</Show>
					</div>
					<p class="text-sm font-medium text-muted-foreground">@{userInfoQuery.data.user.username}</p>
					<Show when={userInfoQuery.data.user?.badges}>
						<div class="mt-3 flex flex-row flex-wrap gap-2 items-center justify-center">
							<For each={userInfoQuery.data.user.badges}>
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
							value={userInfoQuery.data.user.correct_predictions}
							color="#2ECC71"
							label={t('correct')}
						/>
						<ProfileStat
							icon="leaderboard"
							value={`#${userInfoQuery.data.user.ranks.find((r: any) => r.season_type === 'monthly')?.position}`}
							color="#3498DB"
							label={t('rank')}
						/>
						<ProfileStat
							icon="target"
							value={`${store.user?.prediction_accuracy}%`}
							color="#F1C40F"
							label={t('accuracy')}
						/>
						<Show when={userInfoQuery.data.user.longest_win_streak > 3}>
							<ProfileStat
								icon="emoji_events"
								value={userInfoQuery.data.user.longest_win_streak}
								color="#FFC107"
								label={t('max_streak')}
							/>
						</Show>
					</div>
				</div>
				<div class="px-3 space-y-2">
					<For each={userInfoQuery.data.predictions}>
						{(prediction) => (
							<MatchCard match={prediction.match} prediction={prediction} />
						)}
					</For>
				</div>
			</Show>
		</div>
	)
}
