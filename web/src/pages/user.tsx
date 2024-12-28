import { Link } from '~/components/link'
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

function ProfileStat({ icon, value, label, color }: { icon: string; value: string; label: string, color: string }) {
	return (
		<div class="bg-secondary flex flex-col items-center text-center rounded-2xl p-2">
			<span class="material-symbols-rounded text-2xl mb-1"
						style={{ color }}>
				{icon}
			</span>
			<span class="font-semibold text-lg">{value}</span>
			<span class="text-xs text-muted-foreground">{label}</span>
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
					<p class="text-lg font-semibold mt-2">
						{userInfoQuery.data.user.first_name}
					</p>
					<Link href="/" class="text-muted-foreground flex flex-row items-center">
						<p class="text-sm">@{userInfoQuery.data.user.username}</p>
					</Link>
					<div class="grid grid-cols-3 gap-4 mt-6 w-full px-2">
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
						<Show when={userInfoQuery.data.user.current_win_streak > 3}>
							<ProfileStat
								icon="local_fire_department"
								value={userInfoQuery.data.user.current_win_streak}
								label="Current Streak"
								color="#E74C3C"
							/>
						</Show>
						<Show when={userInfoQuery.data.user.longest_win_streak > 3}>
							<ProfileStat
								icon="emoji_events"
								value={userInfoQuery.data.user.longest_win_streak}
								label="Max Streak"
								color="#FFC107"
							/>
						</Show>
						<Show when={userInfoQuery.data.user.favorite_team}>
							<div class="flex flex-col items-center justify-center text-center bg-secondary rounded-2xl p-2">
							<span class="material-symbols-rounded text-2xl mb-1"
										style={{ color: '#f33333' }}>
								favorite
							</span>
								<div class="flex flex-col items-center">
									<img
										src={userInfoQuery.data.user.favorite_team.crest_url}
										alt={`${userInfoQuery.data.user.favorite_team.name} crest`}
										class="size-5 object-cover rounded-full mb-1.5"
									/>
									<span
										class="font-semibold text-muted-foreground text-xs">{userInfoQuery.data.user.favorite_team.short_name}</span>
								</div>
							</div>
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
