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

export function ProfileStat({ icon, value, label, color }: {
	icon: string;
	value: string;
	label: string,
	color: string
}) {
	return (
		<div class="space-x-2 flex-grow bg-background border flex flex-row items-start text-center rounded-2xl py-2 px-3">
			<span class="py-1 material-symbols-rounded text-[20px]" style={{ color }}>
					{icon}
				</span>
			<div class="flex flex-col items-start justify-start">
				<span class="font-extrabold text-lg">{value}</span>
				<span class="text-xs text-muted-foreground text-nowrap">{label}</span>
			</div>

		</div>
	)
}

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
							Share
						</Button>
					</div>
					<img
						src={userInfoQuery.data.user.avatar_url}
						alt="User avatar"
						class="w-24 h-24 rounded-full object-cover mt-4"
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
					<p class="text-sm font-medium text-muted-foreground">@{userInfoQuery.data.user.username}</p>
					<div class="grid grid-cols-2 gap-2 mt-6 w-full px-2">
						<ProfileStat
							icon="check_circle"
							value={userInfoQuery.data.user.correct_predictions}
							label="Correct"
							color="#2ECC71"
						/>
						<ProfileStat
							icon="leaderboard"
							value={`#${userInfoQuery.data.user.global_rank}`}
							label="Rank"
							color="#3498DB"
						/>
						<ProfileStat
							icon="star"
							value={userInfoQuery.data.user.total_points}
							label="Points Earned"
							color="#F1C40F"
						/>
						<Show when={userInfoQuery.data.user.longest_win_streak > 3}>
							<ProfileStat
								icon="emoji_events"
								value={userInfoQuery.data.user.longest_win_streak}
								label="Max Streak"
								color="#FFC107"
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
