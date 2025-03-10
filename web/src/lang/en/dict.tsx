export const dict = {
	'share': 'Share',
	'edit_profile': 'Edit profile',
	'correct': 'Correct',
	'rank': 'Monthly Rank',
	'points_earned': 'Points Earned',
	'max_streak': 'Max Streak',
	'check_out_profile': 'Check out {first_name}\'s profile',
	'first_name': 'First name',
	'last_name': 'Last name',
	'favorite_team': 'Your favorite team. Extra {points} points for correct prediction.',
	'select_favorite_team': 'Select your favorite team',
	'search_team': 'Search team',
	'save_and_close': 'Save & close',
	'close': 'Закрыть',
	'join_community': 'Join community',
	'join_community_description': 'To discuss matches and get the latest updates',
	'open_in_telegram': 'Open in Telegram',
	'monthly_seasons': '🏅 Monthly Seasons!',
	'monthly_seasons_description': 'Compete for the top spot each month! Points reset monthly, and the first-place winner gets a prize. 🏆 Make your predictions count!',
	'season_ends_on': (date: any) => <>Ends on <span class="text-primary">{date}</span></>,
	'big_season_ends_on': (date: any) => <>Big Season is the same as the football season, ends on <span
		class="text-primary">{date}</span></>,
	'matches': 'Matches',
	'leaderboard': 'Leaderboard',
	'user_position': 'Position',
	'points': 'Points',
	'favorite_team_icon': 'Favorite team',
	'match_duration': '{days}d {hours}h {minutes}m',
	'active_season': 'Active season {{ name }}',
	'show_more': 'Show more',
	'big_season': 'Big Season',
	'your_predictions': 'Your Predictions',
	'win_1': 'Team 1',
	'win_2': 'Team 2',
	'draw': 'Draw',
	'exact_score': 'Exact score',
	'save': 'Save & close',
	'onboarding': [
		{
			title: 'Welcome to MatchPredict!',
			description: 'Make predictions on upcoming match scores or outcomes and earn points.',
		},
		{
			title: 'Earn Points',
			description: <span>Guess the exact score to earn <span class="text-primary font-semibold">+7 points</span> or predict the outcome for <span
				class="text-primary font-semibold">+3 points.</span> Points are locked once the match starts.</span>,
		},
		{
			title: 'Bonus Streaks',
			description: 'Maintain a streak of correct predictions to earn bonus points. The longer your streak, the higher your bonus!',
		},
		{
			title: 'Invite Friends',
			description: <span>Invite friends to join MatchPredict and receive <span
				class="text-primary font-semibold">10%</span> of their prediction points. Grow your network and your rewards!</span>,
		},
		{
			title: 'Climb the Leaderboard',
			description: 'Compete with others and climb the leaderboard each season. Top players will receive exciting prizes!',
		},
		{
			title: 'Get Started!',
			description: 'Let’s dive in and start predicting matches. Good luck!',
		},
	],
	contest: {
		win_tshirt: 'Win a football t-shirt worth $50',
		results_announcement: 'Results will be announced on {{ date }} in the Telegram channel',
		join_channel: 'Join the channel',
	},
	'prediction_saved': 'Prediction saved successfully',
	'save_prediction': 'Save Prediction',
	'saving': 'Saving...',
	'match_outcome': 'Match Outcome',
	'win': 'Win',
	'your_prediction': 'Your Prediction',
	'home': 'Home',
	'away': 'Away',
	'vs': 'VS',
	'odds_home': 'Home:',
	'odds_draw': 'Draw:',
	'odds_away': 'Away:',
	'prediction': {
		'score': 'You will receive 5 points for correctly predicting the exact score',
		'outcome': 'You will receive 3 points for correctly predicting the match outcome',
	},
	'feature': {
		'title': 'New Feature!',
		'description': 'Friends, we want to add cash prizes to the game! Make a monthly contribution, compete in predictions, and win real money. What do you think?',
		'option_yes': 'Yes, I’m ready to participate!',
		'option_no': 'No, I’m not interested',
		'button_submit': 'Submit',
	},
}
