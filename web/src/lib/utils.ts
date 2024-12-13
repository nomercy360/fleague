import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs))
}

export function formatDate(dateString: string, dateTime = false) {
	const daysOfWeek = ['Воскресенье', 'Понедельник', 'Вторник', 'Среда', 'Четверг', 'Пятница', 'Суббота']
	const months = ['Янв', 'Фев', 'Мар', 'Апр', 'Май', 'Июн', 'Июл', 'Авг', 'Сен', 'Окт', 'Ноя', 'Дек']

	const now = new Date()
	const date = new Date(dateString)

	const isToday = now.toDateString() === date.toDateString()
	const isTomorrow = new Date(now.getTime() + 24 * 60 * 60 * 1000).toDateString() === date.toDateString()

	const timeOptions = { hour: '2-digit', minute: '2-digit' }
	const time = date.toLocaleTimeString('ru-RU', timeOptions as any)

	let result: string

	if (isToday) {
		result = `Сегодня`
	} else if (isTomorrow) {
		result = `Завтра`
	} else {
		console.log('Date:', dateString)
		const dayOfWeek = daysOfWeek[date.getDay()]
		const day = date.getDate()
		const month = months[date.getMonth()]
		result = `${dayOfWeek}, ${day} ${month}`
	}

	return dateTime ? `${result}, ${time}` : result
}

export function timeToLocaleString(dateString: string) {
	const date = new Date(dateString)
	const timeOptions = { hour: '2-digit', minute: '2-digit'}
	return date.toLocaleTimeString('ru-RU', timeOptions as any)
}
