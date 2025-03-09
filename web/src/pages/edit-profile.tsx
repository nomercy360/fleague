import { setUser, store } from '~/store'
import { createEffect, onCleanup, onMount, createSignal, Show, Switch, Match } from 'solid-js'
import { TextField, TextFieldInput } from '~/components/ui/text-field'
import { createStore } from 'solid-js/store'
import { useMainButton } from '~/lib/useMainButton'
import { Button } from '~/components/ui/button'
import { IconChevronRight } from '~/components/icons'
import { createQuery } from '@tanstack/solid-query'
import { fetchPresignedUrl, fetchTeams, fetchUpdateUser, uploadToS3 } from '~/lib/api'
import { cn } from '~/lib/utils'
import { useNavigate } from '@solidjs/router'
import { useTranslations } from '~/lib/locale-context'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '~/components/ui/select'

type Team = {
	id: number
	name: string
	short_name: string
	crest_url: string
	country: string
}

export default function EditUserPage() {
	const mainButton = useMainButton()

	const [editUser, setEditUser] = createStore({
		first_name: '',
		last_name: '',
		favorite_team_id: '',
		language_code: '',
		avatar_url: null,
	})

	const [showTeamSelector, setShowTeamSelector] = createSignal(false)

	const [searchTerm, setSearchTerm] = createSignal('')
	const [selectedTeam, setSelectedTeam] = createSignal({} as Team)

	const [imgFile, setImgFile] = createSignal<File | null>(null)
	const [imgUploadProgress, setImgUploadProgress] = createSignal(0)

	const [previewUrl, setPreviewUrl] = createSignal('')

	const teams = createQuery<Team[]>(() => ({
		queryKey: ['teams'],
		queryFn: () => fetchTeams(),
	}))

	const filteredTeams = () =>
		teams.data?.filter((team) =>
			team.name.toLowerCase().includes(searchTerm().toLowerCase()),
		) || []

	const navigate = useNavigate()

	createEffect(() => {
		if (store.user?.username) {
			setEditUser({
				first_name: store.user.first_name,
				last_name: store.user.last_name,
				language_code: store.user.language_code,
			})

			if (store.user.favorite_team) {
				setSelectedTeam(store.user.favorite_team)
			}

			if (store.user.avatar_url) {
				setPreviewUrl(store.user.avatar_url)
			}
		}
	})

	async function updateUser() {
		if (imgFile() && imgFile() !== null) {
			mainButton.showProgress(true)
			try {
				const { cdn_url, url } = await fetchPresignedUrl(imgFile()!.name)
				await uploadToS3(
					url,
					imgFile()!,
					e => {
						setImgUploadProgress(Math.round((e.loaded / e.total) * 100))
					},
					() => {
						setImgUploadProgress(0)
					},
				)
				setEditUser('avatar_url', cdn_url)
			} catch (e) {
				console.error(e)
			}
		}

		mainButton.hideProgress()

		const { error } = await fetchUpdateUser({
			...editUser,
			favorite_team_id: selectedTeam()?.id,
		})

		if (error) {
			return
		}

		setUser({
			...store.user,
			...editUser,
			favorite_team: selectedTeam(),
		})

		navigate('/')
	}

	onMount(() => {
		mainButton.enable('Save & close')
		mainButton.onClick(updateUser)
	})

	onCleanup(() => {
		mainButton.offClick(updateUser)
		mainButton.hide()
	})

	const { t } = useTranslations()

	const handleFileChange = (event: any) => {
		const file = event.target.files[0]
		if (file) {
			const maxSize = 1024 * 1024 * 5 // 7MB

			if (file.size > maxSize) {
				window.Telegram.WebApp.showAlert('File size is too big')
				return
			}

			setImgFile(file)
			setPreviewUrl('')

			const reader = new FileReader()
			reader.onload = e => {
				setPreviewUrl(e.target?.result as string)
			}
			reader.readAsDataURL(file)
		}
	}

	return (
		<div class="flex flex-col items-center justify-center bg-background text-foreground px-2 py-3 gap-3">
			<Show when={!showTeamSelector()}>
				<ImageBox imgURL={previewUrl()} onFileChange={handleFileChange} />
				<TextField>
					<TextFieldInput
						placeholder={t('first_name')}
						value={editUser.first_name}
						onInput={(e) => setEditUser('first_name', e.currentTarget.value)}
					/>
				</TextField>
				<TextField>
					<TextFieldInput
						placeholder={t('last_name')}
						value={editUser.last_name}
						onInput={(e) => setEditUser('last_name', e.currentTarget.value)}
					/>
				</TextField>
				<div class="flex-col w-full">
					<p class="px-2 py-1 text-sm text-muted-foreground">
						App & Notifications Language
					</p>
					<Select
						value={editUser.language_code}
						onChange={(value) => setEditUser('language_code', value as string)}
						options={['ru', 'en']}
						placeholder="App language"
						itemComponent={(props) => <SelectItem item={props.item}>{props.item.rawValue}</SelectItem>}
					>
						<SelectTrigger class="w-full">
							<SelectValue<string>>{(state) => state.selectedOption()}</SelectValue>
						</SelectTrigger>
						<SelectContent />
					</Select>

				</div>

				<div class="mt-2 w-full">
					<p class="px-2 text-sm text-muted-foreground">
						{t('favorite_team', { points: 3 })}
					</p>
					<Button
						size="sm"
						class="mt-1 h-12 w-full justify-between"
						variant="secondary"
						onClick={() => setShowTeamSelector(true)}
					>
          <span class="flex flex-row items-center gap-2">
            <Show
							when={selectedTeam().id}
							fallback={
								<span class="text-muted-foreground">
									{t('select_favorite_team')}
								</span>
							}
						>
							<>
								<img
									src={selectedTeam().crest_url}
									alt={selectedTeam().short_name}
									class="size-6"
								/>
								{selectedTeam().short_name}
							</>
						</Show>
          </span>
						<IconChevronRight class="size-6" />
					</Button>
				</div>
			</Show>

			<Show when={showTeamSelector()}>
				<div class="h-screen flex-col flex items-center justify-start w-full">
					<div class="w-full flex items-center relative">
						<TextField class="flex-grow">
							<TextFieldInput
								placeholder={t('search_team')}
								value={searchTerm()}
								onInput={(e) => setSearchTerm(e.currentTarget.value)}
							/>
							{searchTerm() && (
								<button
									class="z-10 text-muted-foreground absolute right-3 top-3"
									onClick={() => setSearchTerm('')}
								>
									<span class="material-symbols-rounded text-[24px]">close</span>
								</button>
							)}
						</TextField>
					</div>
					<div class="mt-4 grid grid-cols-3 gap-2 w-full overflow-y-scroll pb-[40px]">
						{filteredTeams().map((team) => (
							<button
								class={cn('border flex flex-col items-center p-3 rounded-2xl bg-secondary', selectedTeam()?.id === team.id && 'border-primary')}
								onClick={() => {
									setSelectedTeam(team)
									setShowTeamSelector(false)
								}}
							>
								<img src={team.crest_url} alt={team.name} class="size-12 mb-4" />
								<span class="text-xs text-secondary-foreground">{team.short_name}</span>
							</button>
						))}
					</div>
				</div>
			</Show>
		</div>
	)
}

function ImageBox(props: { imgURL: string; onFileChange: any }) {
	return (
		<div class="mt-5 flex h-full items-center justify-center relative">
			<div class="relative flex size-24 flex-col items-center justify-center gap-2">
				<img
					src={props.imgURL}
					alt="Uploaded image preview"
					class="size-24 rounded-full rounded-full object-cover"
				/>
				<input
					class="absolute size-full cursor-pointer rounded-full opacity-0"
					type="file"
					accept="image/*"
					onChange={props.onFileChange}
				/>
			</div>
			<div
				class="flex items-center justify-center rounded-full size-7 bg-primary text-primary-foreground absolute bottom-0 right-0 pointer-events-none">
				<span class="text-[16px] material-symbols-rounded">
					edit
				</span>
			</div>
		</div>
	)
}

