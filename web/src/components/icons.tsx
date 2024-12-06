import { cn } from '~/lib/utils'
import { ComponentProps, splitProps } from 'solid-js'


type IconProps = ComponentProps<'svg'>


const Icon = (props: IconProps) => {
	const [, rest] = splitProps(props, ['class'])
	return (
		<svg
			viewBox="0 0 24 24"
			fill="none" stroke="currentColor"
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

export function IconFootball(props: IconProps) {
	return (
		<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 32 32" fill="none" {...props}>
			<path
				d="M16 27.25C22.2132 27.25 27.25 22.2132 27.25 16C27.25 9.78679 22.2132 4.74998 16 4.74998C9.78676 4.74998 4.75 9.78679 4.75 16C4.75 22.2132 9.78676 27.25 16 27.25Z"
				stroke="currentColor" stroke-width="1.70382" stroke-linecap="round" stroke-linejoin="round" />
			<path d="M16.0002 9.3996L9.79004 13.9139L12.1627 21.2149H19.8378L22.2105 13.9139L16.0002 9.3996Z"
						stroke="currentColor"
						stroke-width="1.70382" stroke-linecap="round" stroke-linejoin="round" />
			<path d="M16 4.75V9.3997" stroke="currentColor" stroke-width="1.70382" stroke-linecap="round"
						stroke-linejoin="round" />
			<path d="M22.209 13.9141L26.6835 12.481" stroke="currentColor" stroke-width="1.70382" stroke-linecap="round"
						stroke-linejoin="round" />
			<path d="M22.6325 25.0845L19.8379 21.2151" stroke="currentColor" stroke-width="1.70382" stroke-linecap="round"
						stroke-linejoin="round" />
			<path d="M12.1615 21.2151L9.35889 25.0845" stroke="currentColor" stroke-width="1.70382" stroke-linecap="round"
						stroke-linejoin="round" />
			<path d="M5.30762 12.481L9.79011 13.9141" stroke="currentColor" stroke-width="1.70382" stroke-linecap="round"
						stroke-linejoin="round" />
			<path d="M10.3384 23.7307L10.633 25.8963" stroke="currentColor" stroke-width="1.70382" stroke-linecap="round"
						stroke-linejoin="round" />
			<path d="M10.3396 23.7307L8.18994 24.1287" stroke="currentColor" stroke-width="1.70382" stroke-linecap="round"
						stroke-linejoin="round" />
			<path d="M4.93311 13.9217L6.89965 12.9822L5.87259 11.0555" stroke="currentColor" stroke-width="1.70382"
						stroke-linecap="round" stroke-linejoin="round" />
			<path d="M14.4785 4.84558L15.9913 6.41405L17.4881 4.82965" stroke="currentColor" stroke-width="1.70382"
						stroke-linecap="round" stroke-linejoin="round" />
			<path d="M26.1104 11.0555L25.0913 12.9902L27.0658 13.9138" stroke="currentColor" stroke-width="1.70382"
						stroke-linecap="round" stroke-linejoin="round" />
			<path d="M23.811 24.0969L21.6613 23.7307L21.3906 25.8884" stroke="currentColor" stroke-width="1.70382"
						stroke-linecap="round" stroke-linejoin="round" />
		</svg>
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
