import { Button } from '~/components/ui/button'
import { store } from '~/store'
import { createQuery } from '@tanstack/solid-query'
import { fetchReferrals } from '~/lib/api'
import { createEffect, createSignal, For, Show } from 'solid-js'
import { cn } from '~/lib/utils'

export default function FriendsPage() {
	const [isCopied, setIsCopied] = createSignal(false)
	const [points, setPoints] = createSignal(0)

	const generateShareURL = () => {
		const botURL = import.meta.env.DEV
			? 'https://t.me/peatcher_testing_bot/peatch'
			: 'https://t.me/footbon_bot/app'

		return (
			'https://t.me/share/url?' +
			new URLSearchParams({
				url: `${botURL}?startapp=r_${store.user?.id}`,
				text: `Check out ${store.user?.first_name}'s profile`,
			}).toString()
		)
	}

	const shareProfileURL = () => {
		window.Telegram.WebApp.openTelegramLink(generateShareURL())
	}

	const copyProfileURL = () => {
		const profileURL = `https://t.me/footbon_bot/app?startapp=r_${store.user?.id}`
		navigator.clipboard
			.writeText(profileURL)
			.then(() => {
				setIsCopied(true)
				setTimeout(() => setIsCopied(false), 2000)
			})
			.catch(() => console.error('Failed to copy text'))
	}

	const referrals = createQuery(() => ({
		queryKey: ['referrals'],
		queryFn: fetchReferrals,
	}))

	createEffect(() => {
		if (referrals.data) {
			const totalPoints = referrals.data.reduce((acc) => acc + 10, 0)
			setPoints(totalPoints)
		}
	})

	return (
		<div class="p-3">
			<div class="bg-secondary p-3 rounded-2xl flex flex-col items-center justify-center">
				<span class="material-symbols-rounded text-[48px] text-primary-foreground">
					people
				</span>
				<h1 class="text-xl font-bold text-center">Invite Friends & Earn</h1>
				<p class="text-sm text-secondary-foreground text-center mt-2">
					Receive <span class="text-primary">10 points</span> for each referral and <span
					class="text-primary">5 points</span> for their referrals.
				</p>
			</div>
			<div class="mt-6">
				<div class="flex flex-row items-center justify-between w-full">
					<p class="text-lg font-semibold">Your Referrals</p>
					<p class="text-sm text-muted-foreground">+{points()} DPS</p>
				</div>

				<Show
					when={referrals.data?.length > 0}
					fallback={
						<div class="text-center mt-6">
							<p class="text-secondary-foreground">No friends invited yet!</p>
							<span class="material-symbols-rounded text-[48px] mt-4">
								person_off
							</span>
						</div>
					}
				>
					<For each={referrals.data}>
						{(referral) => (
							<div class="mt-4 flex items-center justify-between bg-card rounded-2xl p-3">
								<div class="flex items-center">
									<img
										class="size-10 rounded-full"
										src={referral.avatar_url}
										alt={referral.first_name}
									/>
									<span class="ml-4 font-medium">
										{referral.first_name} {referral.last_name}
									</span>
								</div>
								<span class="font-bold text-blue-500">+10 DPS</span>
							</div>
						)}
					</For>
				</Show>
			</div>
			<div class="mt-6 flex flex-row items-center space-x-2">
				<Button class="w-full" onClick={shareProfileURL}>
					Invite a Friend
				</Button>
				<Button
					class={cn('size-10', isCopied() && 'bg-green-500')}
					onClick={copyProfileURL}
				>
					<span class="shrink-0 material-symbols-rounded text-[24px]">
						content_copy
					</span>
				</Button>
			</div>
		</div>
	)
}
