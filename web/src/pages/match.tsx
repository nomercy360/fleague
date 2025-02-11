import { createEffect, createSignal, For, Show } from 'solid-js'
import { createQuery } from '@tanstack/solid-query'
import { fetchActiveSeasons, fetchLeaderboard, fetchMatchByID, fetchMatches, Season } from '~/lib/api'

import { Drawer, DrawerTrigger } from '~/components/ui/drawer'
import { formatDate } from '~/lib/utils'
import { queryClient, setShowCommunityPopup, showCommunityPopup } from '~/App'
import MatchCard from '~/components/match-card'
import FootballScoreboard from '~/components/score-board'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '~/components/ui/tabs'
import { Link } from '~/components/link'
import { Button } from '~/components/ui/button'

import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '~/components/ui/dialog'
import { useTranslations } from '~/lib/locale-context'
import { store } from '~/store'
import { useNavigate, useParams, useSearchParams } from '@solidjs/router'

export default function MatchPage() {
	const params = useParams()

	const matchQuery = createQuery(() => ({
		queryKey: ['matches', params.id],
		queryFn: () => fetchMatchByID(params.id),
	}))

	return (
		<div>
			<Show when={matchQuery.isSuccess}>
				{JSON.stringify(matchQuery.data)}
			</Show>
		</div>
	)
}
