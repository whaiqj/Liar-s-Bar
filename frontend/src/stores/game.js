import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useGameStore = defineStore('game', () => {
  const gameState = ref(null)
  const chatMessages = ref([])
  const roomId = ref(null)
  const isConnected = ref(false)
  const isMatching = ref(false)

  function updateState(state) {
    gameState.value = state
  }

  function addChat(msg) {
    chatMessages.value.push(msg)
  }

  function clearChat() {
    chatMessages.value = []
  }

  function reset() {
    gameState.value = null
    chatMessages.value = []
    roomId.value = null
    isConnected.value = false
  }

  return {
    gameState, chatMessages, roomId, isConnected, isMatching,
    updateState, addChat, clearChat, reset,
  }
})
