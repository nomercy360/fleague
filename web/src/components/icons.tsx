import { cn } from '~/lib/utils'
import { ComponentProps, splitProps } from 'solid-js'


type IconProps = ComponentProps<'svg'>


const Icon = (props: IconProps) => {
	const [, rest] = splitProps(props, ['class'])
	return (
		<svg
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
			class={cn('size-8', props.class)}
			{...rest}
		/>
	)
}


export function IconCrown(props: IconProps) {
	return (
		<Icon {...props} >
			<path
				d="M11.562 3.266a.5.5 0 0 1 .876 0L15.39 8.87a1 1 0 0 0 1.516.294L21.183 5.5a.5.5 0 0 1 .798.519l-2.834 10.246a1 1 0 0 1-.956.734H5.81a1 1 0 0 1-.957-.734L2.02 6.02a.5.5 0 0 1 .798-.519l4.276 3.664a1 1 0 0 0 1.516-.294z" />
			<path d="M5 21h14" />
		</Icon>
	)
}

export function IconPlot(props: IconProps) {
	return (
		<Icon {...props}>
			<path d="m12 8 6-3-6-3v10" />
			<path d="m8 11.99-5.5 3.14a1 1 0 0 0 0 1.74l8.5 4.86a2 2 0 0 0 2 0l8.5-4.86a1 1 0 0 0 0-1.74L16 12" />
			<path d="m6.49 12.85 11.02 6.3" />
			<path d="M17.51 12.85 6.5 19.15" />
		</Icon>
	)
}

export function IconStar(props: IconProps) {
	return (
		<Icon {...props}>
			<path
				d="M11.525 2.295a.53.53 0 0 1 .95 0l2.31 4.679a2.123 2.123 0 0 0 1.595 1.16l5.166.756a.53.53 0 0 1 .294.904l-3.736 3.638a2.123 2.123 0 0 0-.611 1.878l.882 5.14a.53.53 0 0 1-.771.56l-4.618-2.428a2.122 2.122 0 0 0-1.973 0L6.396 21.01a.53.53 0 0 1-.77-.56l.881-5.139a2.122 2.122 0 0 0-.611-1.879L2.16 9.795a.53.53 0 0 1 .294-.906l5.165-.755a2.122 2.122 0 0 0 1.597-1.16z" />
		</Icon>
	)
}

export function IconMinus(props: IconProps) {
	return (
		<Icon {...props}>
			<path d="M5 12l14 0" />
		</Icon>
	)
}


export function IconPlus(props: IconProps) {
	return (
		<Icon {...props}>
			<path d="M12 5l0 14" />
			<path d="M5 12l14 0" />
		</Icon>
	)
}

export function IconShare(props: IconProps) {
	return (
		<Icon {...props}>
			<path d="M4 12v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8" />
			<polyline points="16 6 12 2 8 6" />
			<line x1="12" x2="12" y1="2" y2="15" />
		</Icon>
	)
}

export function IconChevronDown(props: IconProps) {
	return (
		<Icon {...props}>
			<path d="m6 9 6 6 6-6" />
		</Icon>
	)
}

export function IconChevronRight(props: IconProps) {
	return (
		<Icon {...props}>
			<path d="M9 18 15 12 9 6" />
		</Icon>
	)
}

export function IconRefresh(props: IconProps) {
	return (
		<Icon {...props}>
			<path d="M21 12a9 9 0 0 0-9-9 9.75 9.75 0 0 0-6.74 2.74L3 8" />
			<path d="M3 3v5h5" />
			<path d="M3 12a9 9 0 0 0 9 9 9.75 9.75 0 0 0 6.74-2.74L21 16" />
			<path d="M16 16h5v5" />
		</Icon>
	)
}

export function IconTrophy(props: IconProps) {
	return (
		<Icon {...props}>
			<path d="M6 9H4.5a2.5 2.5 0 0 1 0-5H6" />
			<path d="M18 9h1.5a2.5 2.5 0 0 0 0-5H18" />
			<path d="M4 22h16" />
			<path d="M10 14.66V17c0 .55-.47.98-.97 1.21C7.85 18.75 7 20.24 7 22" />
			<path d="M14 14.66V17c0 .55.47.98.97 1.21C16.15 18.75 17 20.24 17 22" />
			<path d="M18 2H6v7a6 6 0 0 0 12 0V2Z" />
		</Icon>
	)
}

export function IconActivity(props: IconProps) {
	return (
		<Icon {...props}>
			<rect width="18" height="18" x="3" y="3" rx="2" />
			<path d="M17 12h-2l-2 5-2-10-2 5H7" />
		</Icon>
	)
}

export function IconCalendar(props: IconProps) {
	return (
		<Icon {...props}>
			<path d="M8 2v4" />
			<path d="M16 2v4" />
			<rect width="18" height="18" x="3" y="4" rx="2" />
			<path d="M3 10h18" />
			<path d="M8 14h.01" />
			<path d="M12 14h.01" />
			<path d="M16 14h.01" />
			<path d="M8 18h.01" />
			<path d="M12 18h.01" />
			<path d="M16 18h.01" />
		</Icon>
	)
}

export function IconUsers(props: IconProps) {
	return (
		<Icon {...props}>
			<path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
			<circle cx="9" cy="7" r="4" />
			<path d="M22 21v-2a4 4 0 0 0-3-3.87" />
			<path d="M16 3.13a4 4 0 0 1 0 7.75" />
		</Icon>
	)
}

export function IconSparkles(props: IconProps) {
	return (
		<Icon {...props}>
			<path
				d="M9.937 15.5A2 2 0 0 0 8.5 14.063l-6.135-1.582a.5.5 0 0 1 0-.962L8.5 9.936A2 2 0 0 0 9.937 8.5l1.582-6.135a.5.5 0 0 1 .963 0L14.063 8.5A2 2 0 0 0 15.5 9.937l6.135 1.581a.5.5 0 0 1 0 .964L15.5 14.063a2 2 0 0 0-1.437 1.437l-1.582 6.135a.5.5 0 0 1-.963 0z" />
			<path d="M20 3v4" />
			<path d="M22 5h-4" />
			<path d="M4 17v2" />
			<path d="M5 18H3" />
		</Icon>
	)
}
