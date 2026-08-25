<!-- PostStatusBadge.vue — Design §13. Compact, 6px radius. -->
<template>
  <span
    class="inline-flex items-center rounded-[6px] px-2.5 py-0.5 text-[12px] font-medium capitalize"
    :class="classes"
  >
    {{ statusLabel }}
  </span>
</template>

<script setup>
import { computed } from "vue"

const props = defineProps({
  status: { type: String, default: "published" },
})

// Defensive: Phase-1 backend may omit `status` on list rows; default to published.
const statusLabel = computed(() => props.status?.trim() || "Published")
const classes = computed(() => {
  const s = String(props.status || "published").toLowerCase()
  if (s === "published") return "text-[#15803d] bg-[#DCFCE7] border border-[#86EFAC]"
  if (s === "draft") return "text-[#92400E] bg-[#FEF9C3] border border-[#FDE64B]"
  return "text-[#71717A] bg-[#F5F5F6] border border-[#E4E4E7]"
})
</script>
