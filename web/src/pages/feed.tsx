import {
	createSignal,
	onMount,
} from 'solid-js'
import {
	IconUsers,
	IconActivity, IconCalendar, IconChevronRight,
} from '~/components/icons'

import PredictionCard from '~/components/prediction-card'
import { Link } from '~/components/link'
import { cn } from '~/lib/utils'
import { useLocation } from '@solidjs/router'


export default function FeedPage() {
	const dummyUsers = [
		{
			title: 'Сергей бестов',
			subtitle: 'Спортивный аналитик',
			image: '/avatars/sergey.png',
			score: 300,
		},
		{
			title: 'Гор Е',
			subtitle: 'Болельщик Челси',
			image: '/avatars/gor.jpg',
			score: 100,
		},
		{
			title: 'Максим К',
			subtitle: 'Болельщик ЦСКА',
			image: '/avatars/maksim.jpg',
			score: 50,
		},
		{
			title: 'John Doe',
			subtitle: 'Football fan',
			image: '/avatars/user.svg',
			score: 10,
		},
		{
			title: 'Максим К',
			subtitle: 'Болельщик ЦСКА',
			image: '/avatars/maksim.jpg',
			score: 50,
		},
	]

	const [homeTeam, setHomeTeam] = createSignal('')
	const [awayTeam, setAwayTeam] = createSignal('')

	const location = useLocation()

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
						<span class="text-2xl font-semibold">10</span>
						<span class="text-xs text-muted-foreground">Correct predictions</span>
					</div>
					<div class="flex flex-col rounded-2xl py-3 px-4 bg-secondary text-card-foreground w-[100px] self-stretch">
						<span class="text-2xl font-semibold">#123</span>
						<span class="text-xs text-muted-foreground">Global ranking</span>
					</div>
					<div class="flex flex-col rounded-2xl py-3 px-4 bg-secondary text-card-foreground w-[100px] self-stretch">
						<span class="text-2xl font-semibold">300</span>
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

			<PredictionCard />
		</div>
	)
}
