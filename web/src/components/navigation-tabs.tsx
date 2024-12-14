import { Link } from '~/components/link'
import { cn } from '~/lib/utils'
import { IconActivity, IconCalendar, IconUsers } from '~/components/icons'
import { useLocation } from '@solidjs/router'

export default function NavigationTabs() {
	const location = useLocation()

	return (
		<div
			class="flex flex-row items-center space-x-4 border shadow-sm px-2.5 h-14 rounded-[28px] fixed bottom-4 w-[240px] bg-background z-50 transform -translate-x-1/2 left-1/2">
			<div class="grid grid-cols-3 w-full">
				<Link
					href="/"
					class={cn('flex items-center flex-col h-full text-sm gap-1', {
						'text-primary': location.pathname === '/',
					})}
				>
					<IconActivity class="size-6" />
				</Link>
				<Link
					href="/matches"
					class={cn('flex items-center flex-col h-full text-sm gap-1', {
						'text-primary': location.pathname === '/matches',
					})}
				>
					<IconCalendar class="size-6" />
				</Link>
				<Link
					href="/"
					class={cn('flex items-center flex-col h-full text-sm gap-1', {
						'text-primary': location.pathname === '/friends',
					})}
				>
					<IconUsers class="size-6" />
				</Link>
			</div>
		</div>
	)
}
