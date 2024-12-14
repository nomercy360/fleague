import {
	IconChevronRight, IconShare,
} from '~/components/icons'

import UserActivity from '~/components/prediction-card'
import { Link } from '~/components/link'
import { store } from '~/store'
import { Button } from '~/components/ui/button'

export default function FeedPage() {
	function shareProfileURL() {
		const url =
			'https://t.me/share/url?' +
			new URLSearchParams({
				url: 'https://t.me/footbon_bot/app?startapp=u_' + store.user?.username,
			}).toString() +
			`&text=Check out ${store.user?.first_name}'s profile`

		window.Telegram.WebApp.openTelegramLink(url)
	}

	return (
		<div class="bg-background text-foreground pb-24">
			<div class="relativew-full bg-card rounded-b-[10%] px-4 pt-6 pb-8 mb-8 flex flex-col items-center">
				<Button
					class="absolute top-6 left-6"
					onClick={shareProfileURL}
					size="sm"
					variant="secondary"
				>
					<IconShare class="size-6" />
					Share
				</Button>
				<img
					src={store.user?.avatar_url}
					alt="User avatar"
					class="size-24 rounded-full object-cover"
				/>
				<p class="text-lg font-semibold mt-2">
					{store.user?.first_name}
				</p>
				<Link href="/" class="text-muted-foreground flex flex-row items-center">
					<p class="text-sm">
						@{store.user?.username}
					</p>
				</Link>
				<div class="flex flex-row items-center justify-center space-x-2 mt-4">
					<div class="flex flex-col rounded-2xl py-3 px-4 bg-secondary text-card-foreground w-[100px] self-stretch">
						<span class="text-2xl font-semibold">{store.user?.correct_predictions}</span>
						<span class="text-xs text-muted-foreground">Correct predictions</span>
					</div>
					<div class="flex flex-col rounded-2xl py-3 px-4 bg-secondary text-card-foreground w-[100px] self-stretch">
						<span class="text-2xl font-semibold">
							#{store.user?.global_rank}
						</span>
						<span class="text-xs text-muted-foreground">Global ranking</span>
					</div>
					<div class="flex flex-col rounded-2xl py-3 px-4 bg-secondary text-card-foreground w-[100px] self-stretch">
						<span class="text-2xl font-semibold">{store.user?.total_points}</span>
						<span class="text-xs text-muted-foreground">Points earned</span>
					</div>
				</div>
			</div>
			<div class="px-3 mb-6">
				<Link class="flex flex-row h-14 justify-between items-center rounded-2xl p-3 bg-secondary"
							href="/matches">
					<p class="text-sm font-semibold">
						Make a prediction{' '}
						<span class="text-muted-foreground font-normal">12 matches available</span>
					</p>
					<IconChevronRight class="size-6" />
				</Link>
			</div>
			<UserActivity />
		</div>
	)
}
