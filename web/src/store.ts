import { createStore } from 'solid-js/store'
import { createSignal } from 'solid-js'


type User = {
	id: string
	first_name: string
	last_name: string
	username: string
	avatar_url: string
	chat_id: number
	language_code: 'en' | 'ru'
	created_at: string
	token: string
	total_points: number
	total_predictions: number
	correct_predictions: number
	referred_by: string
	global_rank: number
	favorite_team: {
		id: number
		name: string
		short_name: string
		crest_url: string
		country: string
	}
	current_win_streak: number
	longest_win_streak: number
	badges: {
		id: number
		name: string
		icon: string
		color: string
		awarded_at: string
	}[]
}

export const [store, setStore] = createStore<{
	user: User | null
	token: string
	following: number[]
}>({
	user: null as any,
	token: null as any,
	following: [],
})

export const setUser = (user: any) => setStore('user', user)

export const setToken = (token: string) => setStore('token', token)

export const setFollowing = (following: number[]) =>
	setStore('following', following)

export const [editUser, setEditUser] = createStore<any>({
	first_name: '',
	last_name: '',
	title: '',
	description: '',
	avatar_url: '',
	city: '',
	country: '',
	country_code: '',
	badge_ids: [],
	opportunity_ids: [],
})

export const [editCollaboration, setEditCollaboration] =
	createStore<any>({
		badge_ids: [],
		city: '',
		country: '',
		country_code: '',
		description: '',
		is_payable: false,
		opportunity_id: 0,
		title: '',
	})

export const [editPost, setEditPost] = createStore<{
	title: string
	description: string
	image_url: string | null
	country: string | null
	country_code: string | null
	city: string | null
}>({
	title: '',
	description: '',
	image_url: '',
	country: '',
	country_code: '',
	city: '',
})

export const [editCollaborationId, setEditCollaborationId] =
	createSignal<number>(0)

export const [editPostId, setEditPostId] = createSignal<number>(0)
