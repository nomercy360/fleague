import {lazy} from 'solid-js'
import type {RouteDefinition} from '@solidjs/router'

import FeedPage from '~/pages/feed'

export const routes: RouteDefinition[] = [
    {
        path: '/',
        component: FeedPage,
    },
    {
        path: '**',
        component: lazy(() => import('./pages/404')),
    },
]
