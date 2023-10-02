import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useRSAStore = defineStore('RSAkey', {state: () => {
  const key = ref<string>("")
  

  return { key }
}})
