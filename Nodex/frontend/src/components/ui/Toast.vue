<!-- Toast.vue — compact toast (Design §22). Supports success (✓) and error (✕) variants. -->
<template>
  <div
    role="status"
    aria-live="polite"
    class="pointer-events-auto flex items-center gap-2 rounded-[8px] px-4 py-2.5 text-sm shadow-lg ring-1 transition-all duration-200"
    :class="toastClass"
  >
    <Icon :name="error ? 'x' : 'check'" :class="iconClass" size="16" />
    <span>{{ message }}</span>
  </div>
</template>

<script setup>
import { computed } from "vue"
import Icon from "./Icon.vue"

const props = defineProps({
  message: { type: String, required: true },
  error: { type: Boolean, default: false },
})

const toastClass = computed(() =>
  props.error
    ? "bg-[#FFFFFF] text-[#b91c1c] ring-[#fee2e2]"
    : "bg-[#FFFFFF] text-[#15803d] ring-[#E4E4E7]"
)
const iconClass = computed(() =>
  props.error ? "text-[#dc2626]" : "text-[#16A34A]"
)
</script>
