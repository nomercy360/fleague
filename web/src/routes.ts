import { lazy } from 'solid-js'
import type { RouteDefinition } from '@solidjs/router'

import FeedPage from '~/pages/feed'
import MatchesPage from '~/pages/matches'

export const routes: RouteDefinition[] = [
	{
		path: '/',
		component: FeedPage,
	},
	{
		path: '/matches',
		component: MatchesPage,
	},
	{
		'path': '/users/:username',
		'component': lazy(() => import('./pages/user')),
	},
	{
		path: '**',
		component: lazy(() => import('./pages/404')),
	},
]
