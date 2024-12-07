// Пример данных. Замените на реальные данные из API.
import { cn, formatDate } from '~/lib/utils'
import { Button } from '~/components/ui/button'
import { IconPencil } from '~/components/icons'

const completedMatches = [
	{
		tournament: 'UEFA Champions League',
		homeTeam: 'GNK Dinamo Zagreb',
		awayTeam: 'Celtic FC',
		homeScore: 2,
		awayScore: 1,
		predictedWinner: 'Celtic FC',
		resultStatus: 'incorrect',
	},
	{
		tournament: 'UEFA Champions League',
		homeTeam: 'Real Madrid CF',
		awayTeam: 'Atalanta BC',
		homeScore: 3,
		awayScore: 0,
		predictedWinner: 'Real Madrid CF',
		resultStatus: 'correct',
	},
]

const upcomingMatches = [
	{
		tournament: 'UEFA Champions League',
		homeTeam: 'Club Brugge KV',
		awayTeam: 'Liverpool FC',
		matchDate: '2024-12-11 17:45',
		predictedWinner: 'Liverpool FC',
	},
]

const PredictionCard = () => {
	return (
		<div class="space-y-2">
			{/* Завершённые матчи */}
			<div>
				<div class="space-y-2">
					{completedMatches.map((match) => (
						<div class={cn('rounded-xl max-w-md mx-auto p-2.5 bg-secondary flex flex-col justify-between', {
							'bg-green-50': match.resultStatus === 'correct',
							'bg-red-50': match.resultStatus === 'incorrect',
						})}>
							<div class="flex items-center space-x-1 mb-2.5">
								<img src={`/logos/uefa.png`} alt="" class="w-4" />
								<p class="text-xs">UEFA Champions League</p>
							</div>
							<div class="grid grid-cols-2 gap-0.5">
								<div class="flex items-center space-x-1">
									<img src={`/logos/${match.homeTeam}.png`} alt="" class="w-4" />
									<p class="text-xs font-bold">{match.homeTeam}</p>
								</div>
								<span class="text-xs font-bold">{match.homeScore}</span>
								<div class="flex items-center space-x-1">
									<img src={`/logos/${match.awayTeam}.png`} alt="" class="w-4" />
									<p class="text-xs font-bold">{match.awayTeam}</p>
								</div>
								<span class="text-xs font-bold">{match.awayScore}</span>
							</div>
						</div>
					))}
				</div>
			</div>
			<div>
				<div class="space-y-4">
					{upcomingMatches.map((match) => (
						<div class="rounded-xl max-w-md mx-auto p-2.5 bg-secondary flex flex-row items-center justify-between">
							<div class="space-y-2">
								<div class="flex items-center space-x-1">
									<img src={`/logos/uefa.png`} alt="" class="w-4" />
									<p class="text-xs">UEFA Champions League</p>
								</div>
								<div class="grid gap-0.5">
									<div class="flex items-center space-x-1">
										<img src={`/logos/${match.homeTeam}.png`} alt="" class="w-4" />
										<p class="text-xs font-bold">{match.homeTeam}</p>
									</div>
									<div class="flex items-center space-x-1">
										<img src={`/logos/${match.awayTeam}.png`} alt="" class="w-4" />
										<p class="text-xs font-bold">{match.awayTeam}</p>
									</div>
								</div>
								<p class="text-xs text-muted-foreground">{formatDate(match.matchDate)}</p>
							</div>
							<Button variant="default" size="sm">
								<IconPencil />
								{match.predictedWinner} 3:1
							</Button>
						</div>
					))}
				</div>
			</div>
		</div>
	)
}

export default PredictionCard
