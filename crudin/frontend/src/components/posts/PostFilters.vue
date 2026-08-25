<!-- PostFilters.vue — Design §7. Search + Status dropdown + View toggle. -->
<template>
  <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
    <!-- Search -->
    <div class="relative w-full sm:max-w-xl">
      <Icon
        name="search"
        size="16"
        class="absolute left-3 top-1/2 -translate-y-1/2 text-[#71717A]"
      />
      <input
        id="post-search"
        type="search"
        :value="modelValue"
        @input="onSearch"
        placeholder="Search posts..."
        aria-label="Search posts"
        class="w-full rounded-[8px] border border-[#E4E4E7] bg-[#FFFFFF] pl-10 pr-3 py-2 text-[15px] text-[#18181B] placeholder-[#71717A] outline-none transition-colors focus:border-[#4F46E5] focus:ring-2 focus:ring-[#4F46E5]/30"
      />
    </div>

    <div class="flex flex-wrap items-center gap-2">
      <!-- Status dropdown -->
      <div class="relative" ref="statusRef">
        <button
          type="button"
          id="status-filter"
          aria-haspopup="listbox"
          :aria-expanded="statusOpen"
          @click="statusOpen = !statusOpen"
          class="inline-flex w-40 items-center justify-between gap-2 rounded-[8px] border border-[#E4E4E7] bg-[#FFFFFF] px-3 py-2 text-[15px] text-[#18181B] outline-none transition-colors hover:bg-[#F5F5F6] focus:ring-2 focus:ring-[#4F46E5]/30"
        >
          <span>{{ statusLabel }}</span>
          <Icon name="chevron-down" size="16" class="text-[#71717A] icon-pressable" />
        </button>

        <transition
          enter-active-class="transition ease-out duration-150"
          enter-from-class="opacity-0 scale-95"
          enter-to-class="opacity-100 scale-100"
          leave-active-class="transition ease-in duration-150"
          leave-from-class="opacity-100 scale-100"
          leave-to-class="opacity-0 scale-95"
        >
          <div
            v-show="statusOpen"
            class="absolute z-10 mt-1 w-40 overflow-hidden rounded-[8px] border border-[#E4E4E7] bg-[#FFFFFF] shadow-lg"
          >
            <ul role="listbox" aria-label="Status filter" class="py-1 text-[15px]">
              <li
                v-for="opt in statusOptions"
                :key="opt.value"
                :aria-selected="statusFilter === opt.value"
                @click="setStatus(opt.value)"
                class="cursor-pointer px-3 py-1.5 hover:bg-[#F5F3FF] hover:text-[#4F46E5] first:rounded-t-[8px] last:rounded-b-[8px]"
                :class="
                  statusFilter === opt.value
                    ? 'bg-[#EFF0FF] text-[#4F46E5]'
                    : 'text-[#18181B]'
                "
              >
                {{ opt.label }}
              </li>
            </ul>
          </div>
        </transition>
      </div>

      <!-- View toggle. Grid is default; list is inert in Phase-1. -->
      <div role="group" class="inline-flex items-center rounded-[8px] bg-[#F5F5F6] p-1 ring-1 ring-[#E4E4E7]">
        <button
          type="button"
          @click="setView('grid')"
          :aria-pressed="view === 'grid'"
          class="rounded-[6px] p-1.5 transition-colors"
          :class="view === 'grid' ? 'bg-[#4F46E5] text-[#FFFFFF]' : 'text-[#71717A] hover:text-[#18181B]'"
          aria-label="Grid view"
          title="Grid view"
        >
          <Icon name="layout-grid" size="16" class="icon-pressable" />
        </button>
        <button
          type="button"
          aria-label="List view (coming soon)"
          title="List view (coming soon)"
          class="rounded-[6px] p-1.5 text-[#71717A] opacity-60"
          aria-disabled="true"
        >
          <Icon name="list" size="16" class="icon-pressable" />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, computed } from "vue"
import Icon from "../ui/Icon.vue"

const props = defineProps({
  modelValue: { type: String, default: "" },
  statusFilter: { type: String, default: "" },
  view: { type: String, default: "grid" },
  statusOptions: {
    type: Array,
    default: () => [
      { label: "All Status", value: "" },
      { label: "Published", value: "published" },
      { label: "Draft", value: "draft" },
    ],
  },
})
const emit = defineEmits(["update:modelValue", "update:statusFilter", "update:view"])

const statusOpen = ref(false)
const statusRef = ref(null)

const statusLabel = computed(() => {
  const found = props.statusOptions.find((o) => o.value === props.statusFilter)
  return found ? found.label : "All Status"
})

function onSearch(e) {
  emit("update:modelValue", e.target.value)
}
function setStatus(value) {
  emit("update:statusFilter", value)
  statusOpen.value = false
}
function setView(v) {
  emit("update:view", v)
}

function onOutside(e) {
  if (statusOpen.value && statusRef.value && !statusRef.value.contains(e.target)) {
    statusOpen.value = false
  }
}
onMounted(() => {
  document.addEventListener("click", onOutside)
  document.addEventListener("keydown", (e) => {
    if (e.key === "Escape") statusOpen.value = false
  })
})
onBeforeUnmount(() => {
  document.removeEventListener("click", onOutside)
})
</script>
