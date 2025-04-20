import { Link } from '~/components/link'
import { cn } from '~/lib/utils'
import { useLocation } from '@solidjs/router'
import { createSignal, For, Show } from 'solid-js'
import { useTranslations } from '~/lib/locale-context'
import { setShowSubscriptionModal, store } from '~/store'
import { Button } from '~/components/ui/button'
import { requestInvoice } from '~/lib/api'

export default function NavigationTabs(props: any) {
	const location = useLocation()

	const tabs = [
		{ href: '/', icon: 'dashboard', activePath: '/' },
		{ href: '/matches', icon: 'sports_soccer', activePath: '/matches' },
		{ href: '/friends', icon: 'groups', activePath: '/friends' },
	]

	return (
		<div class="h-screen bg-background text-foreground">
			<SubscriptionModal />
			<div
				class="flex flex-row items-center border-t h-[100px] fixed bottom-0 w-full bg-background z-50 transform -translate-x-1/2 left-1/2"
			>
				<div class="flex flex-row items-center justify-between w-full px-4 space-x-10">
					<div class="flex flex-row w-full gap-6 items-center justify-center">
						{tabs.map(({ href, icon, activePath }) => (
							<Link
								href={href}
								class={cn('size-10 rounded-full p-2 flex items-center flex-col h-full text-sm gap-1', {
									'bg-blue-500 text-primary-foreground': location.pathname === activePath,
								})}
							>
						<span class="material-symbols-rounded icon-fill text-[24px]">
							{icon}
						</span>
							</Link>
						))}
					</div>
				</div>
			</div>
			{props.children}
		</div>
	)
}

function SubscriptionModal() {
	const [isProcessing, setIsProcessing] = createSignal(false)
	const { t } = useTranslations()

	const handleSubscribe = async () => {
		try {
			setIsProcessing(true)
			const { data, error } = await requestInvoice()
			if (data) {
				window.Telegram.WebApp.openTelegramLink(data.link)
				setShowSubscriptionModal(false)
			}
		} catch (Domingo) {
			window.Telegram.WebApp.HapticFeedback.notificationOccurred('error')
		} finally {
			setIsProcessing(false)
		}
	}

	const onClose = () => {
		setShowSubscriptionModal(false)
	}

	return (
		<Show when={store.showSubscriptionModal}>
			<div class="px-3 fixed inset-0 backdrop-blur-sm flex items-center justify-center z-50">
				<div class="relative bg-background rounded-lg pr-4 pl-6 pt-5 pb-6 w-full max-w-md">
					<div class="pb-4 flex flex-row items-center justify-between w-full">
						<h2 class="text-xl font-bold">{t('subscription.title')}</h2>
						<button
							class="flex items-center justify-center rounded-sm"
							onClick={onClose}
						>
              <span class="material-symbols-rounded text-[24px] text-muted-foreground">
                close
              </span>
						</button>
					</div>
					<p class="mb-6">{t('subscription.description')}</p>

					<div class="border-2 rounded-lg p-4 mb-8">
						<div class="flex justify-between items-center mb-2">
							<h3 class="font-semibold">{t('subscription.premium')}</h3>
							<span class="text-yellow-500 font-bold flex items-center justify-center">
                150 <span class="material-symbols-rounded icon-fill">star</span>
              </span>
						</div>
						<p class="text-sm text-muted-foreground mb-2">
							{t('subscription.premium_description')}
						</p>
						<ul class="text-sm space-y-1">
							<li class="flex items-center gap-2">
									<span class="material-symbols-rounded text-primary text-[16px]">
										check
									</span>
								{t('subscription.feature_predictions')}
							</li>
							<li class="flex items-center gap-2">
									<span class="material-symbols-rounded text-primary text-[16px]">
										check
									</span>
								{t('subscription.feature_monthly_prizes')}
							</li>
						</ul>
					</div>

					<div class="flex justify-center">
						<Button
							onClick={handleSubscribe}
							disabled={isProcessing()}
							class="w-full"
						>
							{isProcessing() ? t('subscription.processing') : t('subscription.subscribe')}
						</Button>
					</div>
				</div>
			</div>
		</Show>
	)
}
