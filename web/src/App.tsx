import { createEffect, createSignal, Match, Switch } from 'solid-js'
import { setToken, setUser } from './store'
import { API_BASE_URL } from '~/lib/api'
import { NavigationProvider } from './lib/useNavigation'
import { useNavigate } from '@solidjs/router'
import { QueryClient, QueryClientProvider } from '@tanstack/solid-query'
import Toast from '~/components/toast'
import NavigationTabs from '~/components/navigation-tabs'

export const queryClient = new QueryClient({
	defaultOptions: {
		queries: {
			retry: 2,
			staleTime: 1000 * 60 * 5, // 5 minutes
			gcTime: 1000 * 60 * 5, // 5 minutes
		},
		mutations: {
			retry: 2,
		},
	},
})

function transformStartParam(startParam?: string): string | null {
	if (!startParam) return null

	// Check if the parameter starts with "redirect-to-"
	if (startParam.startsWith('u_')) {
		const path = startParam.slice('u_'.length)

		return '/users/' + path.replace(/-/g, '/')
	} else {
		return null
	}
}

export default function App(props: any) {
	const [isAuthenticated, setIsAuthenticated] = createSignal(false)
	const [isLoading, setIsLoading] = createSignal(true)

	const navigate = useNavigate()

	createEffect(async () => {
		const initData = window.Telegram.WebApp.initData

		console.log('WEBAPP:', window.Telegram)

		try {
			const resp = await fetch(`${API_BASE_URL}/auth/telegram?` + initData, {
				method: 'POST',
			})

			const data = await resp.json()

			setUser(data.user)
			setToken(data.token)

			window.Telegram.WebApp.ready()
			window.Telegram.WebApp.expand()
			window.Telegram.WebApp.disableClosingConfirmation()
			window.Telegram.WebApp.disableVerticalSwipes()

			setIsAuthenticated(true)
			setIsLoading(false)

			const startapp = window.Telegram.WebApp.initDataUnsafe.start_param

			const redirectUrl = transformStartParam(startapp)

			if (redirectUrl) {
				navigate(redirectUrl)
				return
			}
		} catch (e) {
			console.error('Failed to authenticate user:', e)
			setIsAuthenticated(false)
			setIsLoading(false)
		}
	})
	return (

		<NavigationProvider>
			<QueryClientProvider client={queryClient}>
				<Switch>
					<Match when={isAuthenticated()}>
						{props.children}
						<NavigationTabs />
					</Match>
					<Match when={!isAuthenticated() && isLoading()}>
						<div class="min-h-screen w-full flex-col items-start justify-center bg-main" />
					</Match>
					<Match when={!isAuthenticated() && !isLoading()}>
						<div
							class="min-h-screen w-full flex-col items-start justify-center bg-main text-3xl text-main">
							Something went wrong. Please try again later.
						</div>
					</Match>
				</Switch>
				<Toast />
			</QueryClientProvider>
		</NavigationProvider>
	)
}
