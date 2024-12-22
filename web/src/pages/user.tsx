import {
	IconChevronRight,
} from '~/components/icons'

import UserActivity from '~/components/prediction-card'
import { Link } from '~/components/link'
import { store } from '~/store'
import { createQuery } from '@tanstack/solid-query'
import { fetchLeaderboard, fetchUserInfo } from '~/lib/api'
import { useParams } from '@solidjs/router'
import { For, Show } from 'solid-js'
import MatchCard from '~/components/match-card'

export default function FeedPage() {
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
		<div class="bg-background text-foreground pb-24">
			<Show when={userInfoQuery.isLoading}>
				<div class="flex flex-col items-center justify-center h-full">
					<div class="loader" />
				</div>
			</Show>
			<Show when={userInfoQuery.data}>
				<div class="w-full bg-card rounded-b-[10%] px-4 pt-6 pb-8 mb-8 flex flex-col items-center">
					<img
						src={userInfoQuery.data.user.avatar_url}
						alt="User avatar"
						class="size-24 rounded-full object-cover"
					/>
					<p class="text-lg font-semibold mt-2">
						{userInfoQuery.data.user.first_name}
					</p>
					<Link href="/" class="text-muted-foreground flex flex-row items-center">
						<p class="text-sm">
							@{userInfoQuery.data.user.username}
						</p>
					</Link>
					<div class="flex flex-row items-center justify-center space-x-2 mt-4">
						<div class="flex flex-col rounded-2xl py-3 px-4 bg-secondary text-card-foreground w-[100px] self-stretch">
							<span class="text-2xl font-semibold">{userInfoQuery.data.user.correct_predictions}</span>
							<span class="text-xs text-muted-foreground">Correct predictions</span>
						</div>
						<div class="flex flex-col rounded-2xl py-3 px-4 bg-secondary text-card-foreground w-[100px] self-stretch">
							<span class="text-2xl font-semibold">
								#123
							</span>
							<span class="text-xs text-muted-foreground">Global ranking</span>
						</div>
						<div class="flex flex-col rounded-2xl py-3 px-4 bg-secondary text-card-foreground w-[100px] self-stretch">
							<span class="text-2xl font-semibold">{userInfoQuery.data.user.total_points}</span>
							<span class="text-xs text-muted-foreground">Points earned</span>
						</div>
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
