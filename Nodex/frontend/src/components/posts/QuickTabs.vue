<!-- QuickTabs.vue — Nodex §21. All | Favorites | Archive | Trash with live counts. -->
<template>
  <div class="flex flex-wrap gap-2">
    <RouterLink
      v-for="tab in tabs"
      :key="tab.to"
      :to="tab.to"
      class="inline-flex items-center gap-2 rounded-full px-4 py-1.5 text-[14px] font-medium transition-colors border"
      :class="tab.active ? 'bg-[#EFF0FF] border-[#6D28D9] text-[#6D28D9]' : 'bg-white border-[#E2E8F0] text-[#64748B] hover:bg-[#F8FAFC] hover:text-[#111827]'"
    >
      {{ tab.label }} <span class="rounded-full bg-white px-1.5 py-0.5 text-[12px]" :class="tab.active ? 'text-[#6D28D9]' : 'text-[#94A3B8]'">{{ tab.count }}</span>
    </RouterLink>
  </div>
</template>

<script setup>
import { computed } from "vue"
import { useRoute } from "vue-router"
const props = defineProps({
  counts: { type: Object, default: () => ({ all: 0, favorites: 0, archive: 0, trash: 0 }) }
})
const route = useRoute()
const tabs = computed(() => [
  { label: "All", to: "/", count: props.counts.all, active: route.path === "/" },
  { label: "Favorites", to: "/favorites", count: props.counts.favorites, active: route.path === "/favorites" },
  { label: "Archive", to: "/archive", count: props.counts.archive, active: route.path === "/archive" },
  { label: "Trash", to: "/trash", count: props.counts.trash, active: route.path === "/trash" },
])
</script>
