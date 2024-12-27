import { store } from '~/store'
import { addToast } from '~/components/toast'

export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL as string

export async function apiRequest(endpoint: string, options: RequestInit = {}) {
	try {
		const response = await fetch(`${API_BASE_URL}/v1${endpoint}`, {
			...options,
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${store.token}`,
				...(options.headers || {}),
			},
		})

		const data = await response.json()

		if (!response.ok) {
			const errorMessage =
				Array.isArray(data?.error)
					? data.error.join('\n')
					: typeof data?.error === 'string'
						? data.error
						: 'An error occurred'

			addToast(errorMessage)

			return { data: null }
		}

		return { data }
	} catch (err) {
		const errorMessage = 'A network error occurred. Please try again later.'

		addToast(errorMessage)

		return { data: null }
	}
}

export const fetchMatches = async () => {
	const { data } = await apiRequest('/matches', {
		method: 'GET',
	})

	return data.reduce((acc: any, match: any) => {
		const date = new Date(match.match_date).toDateString()
		if (!acc[date]) acc[date] = []
		acc[date].push(match)
		return acc
	}, {})
}

export const fetchLeaderboard = async () => {
	const { data } = await apiRequest('/leaderboard', {
		method: 'GET',
	})

	return data
}

export const fetchActiveSeason = async () => {
	const { data } = await apiRequest('/seasons/active', {
		method: 'GET',
	})

	return data
}

export const fetchUserInfo = async (username: string) => {
	const { data } = await apiRequest(`/users/${username}`, {
		method: 'GET',
	})

	return data
}

export const uploadToS3 = (
	url: string,
	file: File,
	onProgress: (e: ProgressEvent) => void,
	onFinished: () => void,
): Promise<void> => {
	return new Promise<void>((resolve, reject) => {
		const req = new XMLHttpRequest()
		req.onreadystatechange = () => {
			if (req.readyState === 4) {
				if (req.status === 200) {
					onFinished()
					resolve()
				} else {
					reject(new Error('Failed to upload file'))
				}
			}
		}
		req.upload.addEventListener('progress', onProgress)
		req.open('PUT', url)
		req.send(file)
	})
}

export const fetchPresignedUrl = async (file: string) => {
	const { data } = await apiRequest(`/presigned-url?filename=${file}`, {
		method: 'GET',
	})

	return data
}

export type PredictionRequest = {
	match_id: number
	predicted_home_score: number | null
	predicted_away_score: number | null
	predicted_outcome: string | null
}


export type MatchResponse = {
	id: number
	tournament: string
	home_team: {
		id: number
		name: string
		short_name: string
		crest_url: string
		country: string
		abbreviation: string
	}
	away_team: {
		id: number
		name: string
		short_name: string
		crest_url: string
		country: string
		abbreviation: string
	}
	match_date: string
	status: string
	away_score: any
	home_score: any
	prediction: any
}

export type PredictionResponse = {
	id: number
	user_id: number
	match_id: number
	predicted_outcome: any
	predicted_home_score: number
	predicted_away_score: number
	points_awarded: number
	created_at: string
	completed_at: string
	match: MatchResponse
}


export const saveMatchPrediction = async (prediction: PredictionRequest) => {
	return await apiRequest('/predictions', {
		method: 'POST',
		body: JSON.stringify(prediction),
	})
}

export const fetchPredictions = async () => {
	const { data } = await apiRequest('/predictions', {
		method: 'GET',
	})

	return data
}

export const fetchReferrals = async () => {
	const { data } = await apiRequest('/referrals', {
		method: 'GET',
	})

	return data
}

export const fetchTeams = async () => {
	const { data } = await apiRequest('/teams', {
		method: 'GET',
	})

	return data
}

export const fetchUpdateUser = async (user: any) => {
	return await apiRequest('/users', {
		method: 'PUT',
		body: JSON.stringify(user),
	})
}
