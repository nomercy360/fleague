import { Link } from '~/components/link'
import { cn } from '~/lib/utils'
import { IconActivity, IconCalendar, IconUsers } from '~/components/icons'
import { useLocation } from '@solidjs/router'

export default function NavigationTabs() {
	const location = useLocation()

	return (
		<div
			class="flex flex-row items-center space-x-4 border-t px-2.5 h-20 fixed bottom-0 left-0 right-0 bg-background z-50">
			<div class="grid grid-cols-3 w-full">
				<Link
					href="/"
					class={cn('flex items-center flex-col h-full text-sm gap-1', {
						'text-primary': location.pathname === '/',
					})}
				>
					<IconActivity class="size-6" />
					Activity
				</Link>
				<Link
					href="/matches"
					class={cn('flex items-center flex-col h-full text-sm gap-1', {
						'text-primary': location.pathname === '/matches',
					})}
				>
					<IconCalendar class="size-6" />
					Matches
				</Link>
				<Link
					href="/"
					class={cn('flex items-center flex-col h-full text-sm gap-1', {
						'text-primary': location.pathname === '/friends',
					})}
				>
					<IconUsers class="size-6" />
					Friends
				</Link>
			</div>
		</div>
	)
}
