// Lightweight toast helper. App.vue owns the toast stack and provides it.
import { inject } from "vue"

export const toastSymbol = "postManagerToast"

export function useToast() {
  const t = inject(toastSymbol)
  if (!t) throw new Error("useToast() must be used within App.vue")
  return t
}
