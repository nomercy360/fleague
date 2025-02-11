import { lazy } from 'solid-js'
import type { RouteDefinition } from '@solidjs/router'

import FeedPage from '~/pages/feed'
import MatchesPage from '~/pages/matches'
import NavigationTabs from '~/components/navigation-tabs'
import FriendsPage from '~/pages/friends'
import OnboardingPage from '~/pages/onboarding'
import MatchPage from '~/pages/match'

export const routes: RouteDefinition[] = [
	{
		path: '/',
		component: NavigationTabs,
		children: [
			{
				'path': '/',
				'component': FeedPage,
			},
			{
				'path': '/matches',
				'component': MatchesPage,
			},
			{
				'path': '/friends',
				'component': FriendsPage,
			},
		],
	},
	{
		'path': '/users/:username',
		'component': lazy(() => import('./pages/user')),
	},
	{
		'path': '/edit-profile',
		'component': lazy(() => import('./pages/edit-profile')),
	},
	{
		'path': '/onboarding',
		'component': OnboardingPage,
	},
	{
		'path': '/matches/:id',
		'component': MatchPage,
	},
	{
		path: '**',
		component: lazy(() => import('./pages/404')),
	},
]
