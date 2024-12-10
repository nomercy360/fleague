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
		const date = new Date(match.match_date).toLocaleDateString()
		if (!acc[date]) acc[date] = []
		acc[date].push(match)
		return acc
	}, {})
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

export type Prediction = {
	match_id: number
	predicted_home_score: number | null
	predicted_away_score: number | null
	predicted_outcome: string | null
}

export const saveMatchPrediction = async (prediction: Prediction) => {
	await apiFetch({
		endpoint: `/predictions`,
		method: 'POST',
		body: prediction,
	})
}
