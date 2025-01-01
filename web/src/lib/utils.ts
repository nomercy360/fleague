import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs))
}

export function formatDate(dateString: string, dateTime = false) {
	const daysOfWeek = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday']
	const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']

	const now = new Date()
	const date = new Date(dateString)

	const isToday = now.toDateString() === date.toDateString()

	const tomorrow = new Date(now.getTime() + 24 * 60 * 60 * 1000)
	const isTomorrow = tomorrow.toDateString() === date.toDateString()

	const yesterday = new Date(now.getTime() - 24 * 60 * 60 * 1000)
	const isYesterday = yesterday.toDateString() === date.toDateString()

	// Prepare time string if needed
	const timeOptions = { hour: '2-digit', minute: '2-digit' } as const
	const time = date.toLocaleTimeString('en-US', timeOptions)

	let result: string

	if (isToday) {
		result = 'Today'
	} else if (isYesterday) {
		result = 'Yesterday'
	} else if (isTomorrow) {
		result = 'Tomorrow'
	} else {
		const day = date.getDate()
		const month = months[date.getMonth()]
		const dayOfWeek = daysOfWeek[date.getDay()]

		if (date > now) {
			// Future date beyond tomorrow: show day of week and date
			result = `${dayOfWeek}, ${day} ${month}`
		} else {
			// Past date beyond yesterday: just show date
			result = `${day} ${month}`
		}
	}

	return dateTime ? `${result}, ${time}` : result
}

export function timeToLocaleString(dateString: string) {
	const date = new Date(dateString)

	const timeOptions: Intl.DateTimeFormatOptions = {
		hour: '2-digit',
		minute: '2-digit',
		hour12: false,
	}

	return date.toLocaleTimeString('en-US', timeOptions)
}
