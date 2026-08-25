<!-- Pagination.vue — Design §18. "Showing 1 to 6 of 12 posts | Previous 1 2 3 Next" -->
<template>
  <div
    class="flex flex-col items-center justify-between gap-3 border-t border-[#E4E4E7] pt-4 text-[15px] sm:flex-row"
  >
    <span class="text-[#71717A]">
      Showing {{ from }} to {{ to }} of {{ total }} posts
    </span>

    <nav class="flex items-center gap-1" aria-label="Pagination">
      <button
        type="button"
        @click="emit('page', page - 1)"
        :disabled="page <= 1 || loading"
        class="rounded-[8px] border border-[#E4E4E7] bg-[#FFFFFF] px-2.5 py-1 text-[#71717A] hover:text-[#18181B] hover:bg-[#F5F5F6] disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-[#FFFFFF] focus:outline-none focus:ring-2 focus:ring-[#4F46E5]/30"
      >
        Previous
      </button>

      <template v-for="p in pageNumbers" :key="p.label">
        <button
          v-if="p.type === 'page'"
          type="button"
          @click="emit('page', p.value)"
          :aria-current="p.value === page ? 'page' : null"
          class="rounded-[8px] px-2.5 py-1 text-[15px] font-medium transition-colors"
          :class="
            p.value === page
              ? 'border border-[#4F46E5] bg-[#4F46E5]/10 text-[#4F46E5]'
              : 'text-[#71717A] hover:text-[#18181B] hover:bg-[#F5F5F6]'
          "
        >
          {{ p.value }}
        </button>
        <span
          v-else
          class="px-1 text-[#71717A]"
          aria-hidden="true"
        >
          {{ p.label }}
        </span>
      </template>

      <button
        type="button"
        @click="emit('page', page + 1)"
        :disabled="page >= totalPages || loading"
        class="rounded-[8px] border border-[#E4E4E7] bg-[#FFFFFF] px-2.5 py-1 text-[#71717A] hover:text-[#18181B] hover:bg-[#F5F5F6] disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-[#FFFFFF] focus:outline-none focus:ring-2 focus:ring-[#4F46E5]/30"
      >
        Next
      </button>
    </nav>
  </div>
</template>

<script setup>
import { computed } from "vue"

const props = defineProps({
  page: { type: Number, default: 1 },
  limit: { type: Number, default: 20 },
  total: { type: Number, default: 0 },
  totalPages: { type: Number, default: 1 },
  loading: { type: Boolean, default: false },
})
const emit = defineEmits(["page"])

const from = computed(() => (props.total === 0 ? 0 : (props.page - 1) * props.limit + 1))
const to = computed(() => (props.total === 0 ? 0 : Math.min(props.page * props.limit, props.total)))

// Build a windowed page-number list (Design §18: Previous 1 2 3 Next).
const pageNumbers = computed(() => {
  const n = props.totalPages
  const items = []
  const window = 1
  const current = props.page
  for (let i = 1; i <= n; i++) {
    if (n <= 5 || Math.abs(i - current) <= window || i === 1 || i === n) {
      items.push({ type: "page", value: i })
    } else if (
      items[items.length - 1] &&
      items[items.length - 1].type === "page"
    ) {
      items.push({ type: "ellipsis", label: "…" })
    }
  }
  return items
})
</script>
