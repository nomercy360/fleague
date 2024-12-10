import {
	IconChevronRight,
} from '~/components/icons'

import UserActivity from '~/components/prediction-card'
import { Link } from '~/components/link'
import { store } from '~/store'

export default function FeedPage() {
	return (
		<div class="bg-background text-foreground pb-24">
			<div class="w-full bg-card rounded-b-[10%] px-4 pt-6 pb-8 mb-8 flex flex-col items-center">
				<img
					src={window.Telegram.WebApp.initDataUnsafe.user.photo_url}
					alt=""
					class="size-24 rounded-full object-cover"
				/>
				<p class="text-lg font-semibold mt-2">
					{window.Telegram.WebApp.initDataUnsafe.user.first_name}
				</p>
				<Link href="/" class="text-muted-foreground flex flex-row items-center">
					<p class="text-sm">
						@{window.Telegram.WebApp.initDataUnsafe.user.username}
					</p>
				</Link>
				<div class="flex flex-row items-center justify-center space-x-2 mt-4">
					<div class="flex flex-col rounded-2xl py-3 px-4 bg-secondary text-card-foreground w-[100px] self-stretch">
						<span class="text-2xl font-semibold">{store.user?.correct_predictions}</span>
						<span class="text-xs text-muted-foreground">Correct predictions</span>
					</div>
					<div class="flex flex-col rounded-2xl py-3 px-4 bg-secondary text-card-foreground w-[100px] self-stretch">
						<span class="text-2xl font-semibold">#123</span>
						<span class="text-xs text-muted-foreground">Global ranking</span>
					</div>
					<div class="flex flex-col rounded-2xl py-3 px-4 bg-secondary text-card-foreground w-[100px] self-stretch">
						<span class="text-2xl font-semibold">{store.user?.total_points}{' '}DPS</span>
						<span class="text-xs text-muted-foreground">Points earned</span>
					</div>
				</div>
			</div>
			<div class="px-4 mb-6">
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
