import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs))
}

export function formatDate(dateString: string) {
	const daysOfWeek = ['Воскресенье', 'Понедельник', 'Вторник', 'Среда', 'Четверг', 'Пятница', 'Суббота']
	const months = ['Янв', 'Фев', 'Мар', 'Апр', 'Май', 'Июн', 'Июл', 'Авг', 'Сен', 'Окт', 'Ноя', 'Дек']

	const now = new Date()
	const date = new Date(dateString)

	const isToday = now.toDateString() === date.toDateString()
	const isTomorrow = new Date(now.getTime() + 24 * 60 * 60 * 1000).toDateString() === date.toDateString()

	const timeOptions = { hour: '2-digit', minute: '2-digit', timeZone: 'UTC' }
	const time = date.toLocaleTimeString('ru-RU', timeOptions as any)

	if (isToday) {
		return `Сегодня ${time}`
	} else if (isTomorrow) {
		return `Завтра ${time}`
	} else {
		const dayOfWeek = daysOfWeek[date.getDay()]
		const day = date.getDate()
		const month = months[date.getMonth()]
		return `${dayOfWeek}, ${day} ${month} ${time}`
	}
}
