import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs))
}

export function formatDate(dateString: string, dateTime = false, locale: 'en' | 'ru' = 'en') {
	const daysOfWeek = {
		en: ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'],
		ru: ['вс.', 'пн.', 'вт.', 'ср.', 'чт.', 'пт.', 'сб.'],
	}

	const months = {
		en: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
		ru: ['Янв', 'Фев', 'Мар', 'Апр', 'Май', 'Июн', 'Июл', 'Авг', 'Сен', 'Окт', 'Ноя', 'Дек'],
	}

	const localizedLabels = {
		en: { today: 'Today', yesterday: 'Yesterday', tomorrow: 'Tomorrow' },
		ru: { today: 'Сегодня', yesterday: 'Вчера', tomorrow: 'Завтра' },
	}

	const now = new Date()
	const date = new Date(dateString)

	const formatedLocale = locale === 'en' ? 'en-US' : 'ru-RU'

	const isToday = now.toDateString() === date.toDateString()

	const tomorrow = new Date(now.getTime() + 24 * 60 * 60 * 1000)
	const isTomorrow = tomorrow.toDateString() === date.toDateString()

	const yesterday = new Date(now.getTime() - 24 * 60 * 60 * 1000)
	const isYesterday = yesterday.toDateString() === date.toDateString()

	// Prepare time string if needed
	const timeOptions = { hour: '2-digit', minute: '2-digit' } as const
	const time = date.toLocaleTimeString(formatedLocale, timeOptions)

	let result: string

	if (isToday) {
		result = localizedLabels[locale].today
	} else if (isYesterday) {
		result = localizedLabels[locale].yesterday
	} else if (isTomorrow) {
		result = localizedLabels[locale].tomorrow
	} else {
		const day = date.getDate()
		const month = months[locale][date.getMonth()]
		const dayOfWeek = daysOfWeek[locale][date.getDay()]

		if (date > now) {
			result = `${dayOfWeek}, ${day} ${month}`
		} else {
			result = `${day} ${month}`
		}
	}

	return dateTime ? `${result}, ${time}` : result
}


export function timeToLocaleString(dateString: string, locale = 'en') {
	const date = new Date(dateString)
	const formatedLocale = locale === 'en' ? 'en-US' : 'ru-RU'

	const timeOptions: Intl.DateTimeFormatOptions = {
		hour: '2-digit',
		minute: '2-digit',
		hour12: false,
	}

	return date.toLocaleTimeString(formatedLocale, timeOptions)
}
