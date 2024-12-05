import {
	createSignal,
	onCleanup,
	onMount,
} from 'solid-js'
import useDebounce from '~/lib/useDebounce'
import { useMainButton } from '~/lib/useMainButton'
import { useNavigate } from '@solidjs/router'


export const [search, setSearch] = createSignal('')

export default function FeedPage() {
	const updateSearch = useDebounce(setSearch, 350)

	const mainButton = useMainButton()
	const navigate = useNavigate()

	const toCreateCollab = () => {
		navigate('/collaborations/edit')
	}

	const [dropDown, setDropDown] = createSignal(false)

	const closeDropDown = () => {
		setDropDown(false)
	}

	const openDropDown = () => {
		document.body.style.overflow = 'hidden'
		setDropDown(true)
	}

	onMount(() => {
		window.Telegram.WebApp.disableClosingConfirmation()
		// window.Telegram.WebApp.CloudStorage.removeItem('profilePopup')
		// window.Telegram.WebApp.CloudStorage.removeItem('communityPopup')
		// window.Telegram.WebApp.CloudStorage.removeItem('rewardsPopup')
	})


	onCleanup(() => {
		mainButton.hide()
		mainButton.offClick(toCreateCollab)
		mainButton.offClick(openDropDown)
		document.removeEventListener('click', closeDropDownOnOutsideClick)
		document.body.style.overflow = 'auto'
	})

	// if dropdown is open, every click outside of the dropdown will close
	const closeDropDownOnOutsideClick = (e: MouseEvent) => {
		if (
			dropDown() &&
			!e.composedPath().includes(document.getElementById('dropdown-menu')!)
		) {
			closeDropDown()
			document.body.style.overflow = 'auto'
		}
	}

	document.addEventListener('click', closeDropDownOnOutsideClick)

	return (
		<div class="min-h-screen bg-secondary pb-56 pt-[76px]">

			<h1>
				Hey there! This is the feed page. You can find the code for this page in src/pages/feed.tsx
			</h1>
		</div>
	)
}
