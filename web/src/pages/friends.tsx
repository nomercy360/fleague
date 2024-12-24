import { Button } from '~/components/ui/button'
import { store } from '~/store'

export default function FriendsPage() {
	function shareProfileURL() {
		const url =
			'https://t.me/share/url?' +
			new URLSearchParams({
				url: 'https://t.me/footbon_bot/app?startapp=r_' + store.user?.referral_code,
			}).toString() +
			`&text=Check out ${store.user?.first_name}'s profile`

		window.Telegram.WebApp.openTelegramLink(url)
	}


	return (
		<div class="text-white min-h-screen p-4">
			<h1 class="text-xl font-bold text-center">Invite Friends & Earn</h1>
			<p class="text-sm text-secondary-foreground text-center mt-2">
				Receive a <span class="text-primary">10% bonus</span> from your referrals and <span
				class="text-primary">5%</span> more from their referrals.
			</p>

			<div class="bg-card rounded-lg p-4 mt-4">
				<h2 class="text-sm text-secondary-foreground text-center">Available to claim</h2>
				<p class="text-3xl font-bold text-center">0.00 DPS</p>
				<Button
					class="mt-6 w-full"
					disabled
				>
					Claim
				</Button>
			</div>

			<div class="mt-6 text-center">
				<p class="text-secondary-foreground">Total <span class="font-bold">1,054 DPS</span> rewards from referrals</p>
			</div>

			<div class="mt-6">
				<h2 class="text-lg font-semibold">Your Referrals</h2>
				<p class="text-sm text-secondary-foreground">Total 9 friends</p>
				<div class="mt-4 flex items-center justify-between bg-card rounded-lg p-4">
					<div class="flex items-center">
						<img
							class="w-10 h-10 rounded-full"
							src="https://via.placeholder.com/40"
							alt="Maksim"
						/>
						<span class="ml-4 font-medium">Maksim</span>
					</div>
					<span class="font-bold text-blue-500">+223 DPS</span>
				</div>
			</div>

			<Button class="mt-6 w-full"
							onClick={shareProfileURL}>
				Invite a Friend
			</Button>

			<div class="flex justify-between mt-6">
				<button class="text-sm text-secondary-foreground hover:text-white">Earn</button>
				<button class="text-sm text-secondary-foreground hover:text-white">Tasks</button>
				<button class="text-sm text-blue-500 font-semibold">Friends</button>
			</div>
		</div>
	)
}
