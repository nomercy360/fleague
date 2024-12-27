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
		<div class="h-full overflow-y-scroll bg-background text-foreground pb-[120px]">
			<div class="relative w-full bg-card rounded-b-[10%] px-4 pt-6 pb-8 mb-8 flex flex-col items-center">
				<div class="flex flex-row justify-between items-center w-full">
					<Button
						onClick={shareProfileURL}
						size="sm"
						variant="secondary"
					>
						<IconShare class="size-6" />
						Share
					</Button>
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
					<img
						src={store.user?.favorite_team?.crest_url}
						alt={store.user?.favorite_team?.short_name}
						class="size-4 ml-1"
					/>
				</div>
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
			<UserActivity />
		</div>
	)
}
