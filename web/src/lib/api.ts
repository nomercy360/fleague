import { store } from '~/store'

export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL as string
export const CDN_URL = 'https://assets.peatch.io'

export const apiFetch = async ({
																 endpoint,
																 method = 'GET',
																 body = null,
																 showProgress = true,
																 responseContentType = 'json' as 'json' | 'blob',
															 }: {
	endpoint: string
	method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
	body?: any
	showProgress?: boolean
	responseContentType?: string
}) => {
	const headers: { [key: string]: string } = {
		'Content-Type': 'application/json',
		Authorization: `Bearer ${store.token}`,
	}

	try {
		showProgress && window.Telegram.WebApp.MainButton.showProgress(false)

		const response = await fetch(`${API_BASE_URL}/v1${endpoint}`, {
			method,
			headers,
			body: body ? JSON.stringify(body) : undefined,
		})

		if (!response.ok) {
			const errorResponse = await response.json()
			throw { code: response.status, message: errorResponse.message }
		}

		switch (response.status) {
			case 204:
				return true
			default:
				return response[responseContentType as 'json' | 'blob']()
		}
	} finally {
		showProgress && window.Telegram.WebApp.MainButton.hideProgress()
	}
}

export const fetchMatches = async () => {
	const resp = await apiFetch({
		endpoint: '/matches',
	})

	return resp.reduce((acc: any, match: any) => {
		const date = new Date(match.match_date).toDateString()
		if (!acc[date]) acc[date] = []
		acc[date].push(match)
		return acc
	}, {})
}

export const fetchLeaderboard = async () => {
	return await apiFetch({
		endpoint: '/leaderboard',
	})
}

export const fetchActiveSeason = async () => {
	return await apiFetch({
		endpoint: '/seasons/active',
	})
}

export const fetchUserInfo = async (username: string) => {
	return await apiFetch({
		endpoint: '/users/' + username,
	})
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
	const { path, url } = await apiFetch({
		endpoint: `/presigned-url?filename=${file}`,
		showProgress: false,
	})

	return { path, url }
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
	await apiFetch({
		endpoint: `/predictions`,
		method: 'POST',
		body: prediction,
	})
}

export const fetchPredictions = async () => {
	return await apiFetch({
		endpoint: '/predictions',
	})
}

export const fetchReferrals = async () => {
	return await apiFetch({
		endpoint: '/referrals',
	})
}
