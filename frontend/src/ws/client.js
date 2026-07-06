class WebSocketClient {
  constructor() {
    this.ws = null
    this.url = ''
    this.listeners = {}
    this.reconnectTimer = null
    this.connected = false
    this.pendingMessages = []
    this.manualClose = false
  }

  connect() {
    const token = localStorage.getItem('token')
    if (!token) return
    if (this.ws && [WebSocket.CONNECTING, WebSocket.OPEN].includes(this.ws.readyState)) return

    this.manualClose = false
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    this.url = `${protocol}//${window.location.host}/ws`

    this.ws = new WebSocket(`${this.url}?token=${encodeURIComponent(token)}`)

    this.ws.onopen = () => {
      this.connected = true
      this.flushPendingMessages()
      this.emit('connected')
    }

    this.ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data)
        this.emit(msg.type, msg.payload, msg)
        this.emit('message', msg)
      } catch (e) {
        console.error('WS parse error:', e)
      }
    }

    this.ws.onclose = () => {
      this.connected = false
      this.ws = null
      this.emit('disconnected')
      if (!this.manualClose && localStorage.getItem('token')) {
        this.scheduleReconnect()
      }
    }

    this.ws.onerror = (err) => {
      console.error('WS error:', err)
    }
  }

  scheduleReconnect() {
    if (this.reconnectTimer) return
    this.reconnectTimer = setTimeout(() => {
      this.reconnectTimer = null
      this.connect()
    }, 3000)
  }

  send(type, payload = {}) {
    const message = JSON.stringify({ type, payload })
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(message)
      return
    }
    this.pendingMessages.push(message)
    this.connect()
  }

  flushPendingMessages() {
    while (this.pendingMessages.length > 0 && this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(this.pendingMessages.shift())
    }
  }

  on(event, callback) {
    if (!this.listeners[event]) {
      this.listeners[event] = []
    }
    this.listeners[event].push(callback)
  }

  off(event, callback) {
    if (!this.listeners[event]) return
    this.listeners[event] = this.listeners[event].filter(cb => cb !== callback)
  }

  emit(event, ...args) {
    if (!this.listeners[event]) return
    this.listeners[event].forEach(cb => cb(...args))
  }

  disconnect() {
    this.manualClose = true
    this.pendingMessages = []
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
    this.connected = false
  }
}

export const wsClient = new WebSocketClient()
export default wsClient
