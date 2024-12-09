import { createEffect, createSignal, Match, Switch } from 'solid-js'
import { setToken, setUser } from './store'
import { API_BASE_URL } from '~/lib/api'
import { NavigationProvider } from './lib/useNavigation'
import { useNavigate } from '@solidjs/router'
import { QueryClient, QueryClientProvider } from '@tanstack/solid-query'
import Toast from '~/components/toast'
import { ColorModeProvider, ColorModeScript } from '@kobalte/core'
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

			const user = await resp.json()

			setUser(user)
			setToken(user.token)

			window.Telegram.WebApp.ready()
			window.Telegram.WebApp.expand()
			window.Telegram.WebApp.disableClosingConfirmation()
			window.Telegram.WebApp.disableVerticalSwipes()

			setIsAuthenticated(true)
			setIsLoading(false)

			// if there is a redirect url, redirect to it
			// ?startapp=redirect-to=/users/

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
						<div>
							{props.children}
							<NavigationTabs />
						</div>
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
