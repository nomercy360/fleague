import { Button } from '~/components/ui/button'
import { store } from '~/store'
import { createQuery } from '@tanstack/solid-query'
import { fetchReferrals } from '~/lib/api'
import { createEffect, createSignal, For, Show } from 'solid-js'
import { cn } from '~/lib/utils'
import { useTranslations } from '~/lib/locale-context'

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
			}).toString() +
			`&text=Compete with friends in predicting football matches!`
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
			const totalPoints = referrals.data.reduce((acc: any) => acc + 10, 0)
			setPoints(totalPoints)
		}
	})

	const { t } = useTranslations()

	return (
		<div class="h-full p-3 overflow-y-scroll pb-[180px]">
			<button class="relative bg-secondary p-3 rounded-2xl flex flex-col items-start text-start justify-center"
							onClick={shareProfileURL}>
				<span
					class="material-symbols-rounded text-[20px] absolute top-3 right-3 text-secondary-foreground">
					arrow_outward
				</span>
				<span class="material-symbols-rounded text-[32px]">
					people
				</span>
				<h1 class="mt-3 text-base font-bold">
					{t('invite_friends')}
				</h1>
				<p class="text-sm text-secondary-foreground mt-1">
					{t('invite_friends_description')}
				</p>
			</button>
			<div class="grid grid-cols-2 gap-2 mt-2">
				<div class="relative bg-secondary p-3 rounded-2xl flex flex-col items-start text-start justify-center">
					<span
						class="material-symbols-rounded text-[20px] absolute top-3 right-3 text-secondary-foreground">
						arrow_outward
					</span>
					<span class="material-symbols-rounded text-[32px]">
						star
					</span>
					<h1 class="mt-3 text-base font-bold">
						{t('buy_points')}
					</h1>
					<p class="text-sm text-secondary-foreground mt-1">
						{t('buy_points_description')}
					</p>
				</div>
				<div class="bg-secondary p-3 rounded-2xl flex flex-col items-start text-start justify-center">
					<span class="material-symbols-rounded text-[32px]">
						redeem
					</span>
					<h1 class="mt-3 text-base font-bold">
						{t('daily_bonus')}
					</h1>
					<p class="text-sm text-secondary-foreground mt-1">
						{t('daily_bonus_description')}
					</p>
				</div>
			</div>
			<div class="mt-6">
				<div class="flex flex-row items-center justify-between w-full">
					<p class="text-lg font-semibold">
						{t('friends_invited')}
					</p>
				</div>
				<Show when={referrals.data?.length > 0}>
					<For each={referrals.data}>
						{(referral) => (
							<div class="mt-1 flex items-center justify-between bg-card rounded-2xl p-3">
								<div class="flex items-center">
									<img
										class="size-6 rounded-full"
										src={referral.avatar_url}
										alt={referral.first_name}
									/>
									<span class="ml-4 font-medium">
										{referral.first_name} {referral.last_name}
									</span>
								</div>
								<div class="flex items-center">
									<p class="text-sm font-medium text-secondary-foreground mr-0.5">
										+10
									</p>
									<span class="text-[12px] material-symbols-rounded text-yellow-200">star</span>
								</div>
							</div>
						)}
					</For>
				</Show>
			</div>
		</div>
	)
}
