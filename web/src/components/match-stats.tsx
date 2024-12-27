import { MatchResponse } from '~/lib/api'
import { DrawerContent } from '~/components/ui/drawer'
import { formatDate } from '~/lib/utils'

export default function MatchStats(props: { match: MatchResponse }) {
	return (
		<DrawerContent>
			<div class="mx-auto w-full p-4">
				<div class="flex flex-col items-center gap-4">
					<div class="w-full justify-between flex flex-col items-start gap-2">
						<div class="flex flex-row w-full justify-between items-center h-10">
							<div class="flex flex-row items-center space-x-2">
								<img src={props.match.home_team.crest_url} alt="" class="w-6" />
								<p class="mt-2 text-base mb-2">
									{props.match.home_team.short_name}
								</p>
							</div>
							<div class="text-lg font-bold text-foreground size-9 items-center flex justify-center">
								{props.match.home_score ?? '—'}
							</div>
						</div>
						<div class="flex flex-row w-full justify-between items-center h-10">
							<div class="flex flex-row items-center space-x-2">
								<img src={props.match.away_team.crest_url} alt="" class="w-6" />
								<p class="mt-2 text-base mb-2">
									{props.match.away_team.short_name}
								</p>
							</div>
							<div class="text-lg font-bold text-foreground size-9 items-center flex justify-center">
								{props.match.away_score ?? '—'}
							</div>
						</div>
						<p class="text-xs text-muted-foreground mt-2">
							{formatDate(props.match.match_date, true)}
						</p>
					</div>
				</div>
			</div>
		</DrawerContent>
	)
}
