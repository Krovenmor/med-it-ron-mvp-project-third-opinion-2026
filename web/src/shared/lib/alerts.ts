export function playAlertSound(): void {
  const context = new AudioContext()
  const tones = [880, 660, 880]
  tones.forEach((frequency, index) => {
    const oscillator = context.createOscillator()
    const gain = context.createGain()
    const start = context.currentTime + index * 0.22
    oscillator.frequency.value = frequency
    gain.gain.setValueAtTime(0.18, start)
    gain.gain.exponentialRampToValueAtTime(0.001, start + 0.2)
    oscillator.connect(gain).connect(context.destination)
    oscillator.start(start)
    oscillator.stop(start + 0.2)
  })
}

export function showBrowserNotification(title: string, body: string): void {
  if ('Notification' in window && Notification.permission === 'granted') {
    new Notification(title, { body })
  }
}

export function canAskNotificationPermission(): boolean {
  return 'Notification' in window && Notification.permission === 'default'
}
