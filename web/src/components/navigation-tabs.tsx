import { Link } from '~/components/link'
import { cn } from '~/lib/utils'
import { useLocation } from '@solidjs/router'

export default function NavigationTabs(props: any) {
	const location = useLocation()

	const tabs = [
		{ href: '/', icon: 'dashboard', activePath: '/' },
		{ href: '/matches', icon: 'sports_soccer', activePath: '/matches' },
		{ href: '/friends', icon: 'groups', activePath: '/friends' },
	]

	return (
		<div class="h-screen bg-background text-foreground">
			<div
				class="flex flex-row items-start border-t h-[100px] fixed bottom-0 w-full bg-background z-50 transform -translate-x-1/2 left-1/2"
			>
				<div class="px-2.5 py-4 flex flex-row w-full gap-10 items-center justify-center">
					{tabs.map(({ href, icon, activePath }) => (
						<Link
							href={href}
							class={cn('size-10 rounded-full p-2 flex items-center flex-col h-full text-sm gap-1', {
								'bg-blue-500 text-primary-foreground': location.pathname === activePath,
							})}
						>
						<span class="material-symbols-rounded icon-fill text-[24px]">
							{icon}
						</span>
						</Link>
					))}
				</div>
			</div>
			{props.children}
		</div>
	)
}
