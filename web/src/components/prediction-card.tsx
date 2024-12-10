// Пример данных. Замените на реальные данные из API.
import { cn, formatDate } from '~/lib/utils'
import { Button } from '~/components/ui/button'
import { IconRefresh } from '~/components/icons'

const completedMatches = [
	{
		tournament: 'UEFA Champions League',
		homeTeam: 'GNK Dinamo Zagreb',
		awayTeam: 'Celtic FC',
		homeScore: 2,
		awayScore: 1,
		predictedWinner: 'Celtic FC',
		resultStatus: 'incorrect',
		pointsEarned: 0,
	},
	{
		tournament: 'UEFA Champions League',
		homeTeam: 'Real Madrid CF',
		awayTeam: 'Atalanta BC',
		homeScore: 3,
		awayScore: 0,
		predictedWinner: 'Real Madrid CF',
		resultStatus: 'correct',
		pointsEarned: 3,
	},
	{
		tournament: 'UEFA Champions League',
		homeTeam: 'FC Bayern München',
		awayTeam: 'FC Barcelona',
		homeScore: 2,
		awayScore: 1,
		predictedWinner: 'FC Bayern München',
		resultStatus: 'correct',
		pointsEarned: 3,
	},
]

const upcomingMatches = [
	{
		tournament: 'UEFA Champions League',
		homeTeam: 'Club Brugge KV',
		awayTeam: 'Liverpool FC',
		matchDate: '2024-12-11 17:45',
		predictedWinner: 'Liverpool FC',
		prediction: '2:0',
		isEditable: true,
	},
	{
		tournament: 'UEFA Champions League',
		homeTeam: 'FC Internazionale Milano',
		awayTeam: 'FC Red Bull Salzburg',
		matchDate: '2024-12-11 17:45',
		predictedWinner: 'FC Internazionale Milano',
		prediction: '3:0',
		isEditable: false,
	},
]

const PredictionCard = () => {
	return (
		<div class="px-4">
			<p class="mb-1 px-2 text-lg font-semibold">
				All Activity
			</p>
			<div class="space-y-2">
				{completedMatches.map((match) => (
					<div
						class={cn('h-[120px] relative grid grid-cols-3 items-center rounded-2xl max-w-md mx-auto p-2.5 pt-4 bg-card', {
							'border-l-4 border-green-500': match.resultStatus === 'correct',
							'border-l-4 border-red-500': match.resultStatus === 'incorrect',
						})}>
						<div
							class={cn('h-6 rounded-b-xl w-12 flex items-center justify-center text-white text-xs font-semibold absolute top-0 left-1/2 -translate-x-1/2 transform', {
								'bg-green-500': match.resultStatus === 'correct',
								'bg-red-500': match.resultStatus === 'incorrect',
							})}>
							{match.resultStatus === 'correct' ? 'win' : 'lose'}
						</div>
						<div class="flex flex-col items-center space-y-2 text-center">
							<img src={`/logos/${match.homeTeam}.png`} alt="" class="w-8" />
							<p class="text-sm font-bold">{match.homeTeam}</p>
						</div>
						<div class="flex flex-col items-center text-center">
							<p class="text-xs text-muted-foreground text-center">
								{match.tournament}
							</p>
							<span class="text-2xl font-bold text-center">{match.homeScore}:{match.awayScore}</span>
						</div>
						<div class="flex flex-col items-center space-y-2 text-center">
							<img src={`/logos/${match.awayTeam}.png`} alt="" class="w-8" />
							<p class="text-sm font-bold">{match.awayTeam}</p>
						</div>
					</div>
				))}
			</div>
			<p class="mt-6 mb-1 px-2 text-lg font-semibold">
				Upcoming
			</p>
			<div class="space-y-2">
				{upcomingMatches.map((match) => (
					<div
						class="h-[140px] border-l-4 border-blue-500 rounded-2xl p-3 bg-card flex flex-row items-start justify-between">
						<div class="space-y-2">
							<p class="text-xs text-muted-foreground">UEFA Champions League</p>
							<div class="grid gap-0.5">
								<div class="flex items-center space-x-1">
									<img src={`/logos/${match.homeTeam}.png`} alt="" class="w-6" />
									<p class="text-sm font-bold">{match.homeTeam}</p>
								</div>
								<div class="flex items-center space-x-1">
									<img src={`/logos/${match.awayTeam}.png`} alt="" class="w-6" />
									<p class="text-sm font-bold">{match.awayTeam}</p>
								</div>
							</div>
							<p class="text-xs text-muted-foreground">{formatDate(match.matchDate, true)}</p>
						</div>
						<Button variant="secondary"
										class={cn('gap-1', {
											'bg-primary text-primary-foreground': match.isEditable,
											'bg-secondary text-secondary-foreground cursor-not-allowed': !match.isEditable,
										})}
										disabled={!match.isEditable}
										size="sm">
							<span class="text-xs font-semibold">
								{match.prediction}
							</span>
							{match.isEditable && <IconRefresh class="w-4 h-4" />}
							{!match.isEditable && <span class="text-xs font-semibold">Predicted</span>}
						</Button>
					</div>
				))}
			</div>
		</div>
	)
}

export default PredictionCard
