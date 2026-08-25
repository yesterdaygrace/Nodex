<!-- AppSidebar.vue — Design §4. Persistent desktop, drawer on mobile (§23).
     Inert placeholder nav items except Posts (active). -->
<template>
  <!-- Mobile backdrop -->
  <div
    v-if="open"
    class="fixed inset-0 z-20 bg-black/40 lg:hidden"
    @click="emit('close')"
    aria-hidden="true"
  ></div>

  <aside
    :class="[
      'fixed inset-y-0 left-0 z-30 flex w-64 flex-col overflow-y-auto border-r border-[#E4E4E7] bg-[#FFFFFF] transition-transform duration-200',
      open ? 'translate-x-0' : '-translate-x-full lg:translate-x-0',
    ]"
  >
    <!-- Brand — Nodex §6 -->
    <div class="flex items-center gap-2 px-6 py-4 text-[22px] font-bold text-[#111827]">
      <span class="flex h-7 w-7 items-center justify-center rounded-[6px] bg-[#6D28D9] text-[12px] text-white">◇</span>
      <span>Nodex</span>
    </div>

    <!-- Nav -->
    <nav class="flex-1 px-3" aria-label="Primary">
      <div class="space-y-1">
        <template v-for="item in primaryItems" :key="item.name">
          <RouterLink
            v-if="!item.inert"
            :to="item.to"
            :aria-current="item.active ? 'page' : null"
            @click="$emit('navigate', item)"
            class="flex w-full items-center gap-2.5 rounded-[8px] px-2.5 py-2 text-left text-[15px] font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-[#4F46E5]/30 focus:ring-offset-2"
            :class="navClass(item)"
          >
            <Icon :name="item.icon" size="18" class="icon-pressable" />
            <span>{{ item.name }}</span>
          </RouterLink>
          <button
            v-else
            type="button"
            :aria-current="item.active ? 'page' : null"
            :disabled="true"
            class="flex w-full items-center gap-2.5 rounded-[8px] px-2.5 py-2 text-left text-[15px] font-medium transition-colors"
            :class="navClass(item)"
          >
            <Icon :name="item.icon" size="18" class="icon-pressable" />
            <span>{{ item.name }}</span>
          </button>
        </template>
      </div>

      <!-- Secondary nav -->
      <div class="mt-6 space-y-1">
        <div class="px-2 py-1">
          <div class="h-px w-full bg-[#E4E4E7]"></div>
        </div>
        <template v-for="item in secondaryItems" :key="item.name">
          <RouterLink
            v-if="!item.inert"
            :to="item.to"
            :aria-current="item.active ? 'page' : null"
            @click="$emit('navigate', item)"
            class="flex w-full items-center gap-2.5 rounded-[8px] px-2.5 py-2 text-left text-[15px] font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-[#4F46E5]/30 focus:ring-offset-2"
            :class="navClass(item)"
          >
            <Icon :name="item.icon" size="18" class="icon-pressable" />
            <span>{{ item.name }}</span>
          </RouterLink>
          <button
            v-else
            type="button"
            :aria-current="item.active ? 'page' : null"
            :disabled="true"
            class="flex w-full items-center gap-2.5 rounded-[8px] px-2.5 py-2 text-left text-[15px] font-medium transition-colors"
            :class="navClass(item)"
          >
            <Icon :name="item.icon" size="18" class="icon-pressable" />
            <span>{{ item.name }}</span>
          </button>
        </template>
      </div>

      <!-- Folders — Nodex §11 (static mock, + toasts Soon) -->
      <div class="mt-6 px-3">
        <div class="flex items-center justify-between px-2 py-1 text-[11px] font-semibold tracking-widest text-[#64748B]">FOLDERS <button type="button" @click="onComingSoon" class="rounded px-1 text-[14px] leading-none hover:bg-[#EFF0FF] hover:text-[#6D28D9]">+</button></div>
        <div class="mt-1 space-y-0.5">
          <button v-for="f in folders" :key="f.name" @click="filterFolder(f.name)" class="flex w-full items-center gap-2 rounded-[8px] px-2 py-1.5 text-left text-[14px] transition-colors" :class="route.query.folder === f.name ? 'bg-[#EFF0FF] text-[#6D28D9]' : 'text-[#334155] hover:bg-[#F8FAFC]'"><Icon name="folder" size="16" /> <span class="flex-1">{{ f.name }}</span> <span class="text-[12px] text-[#94A3B8]">{{ f.count }}</span></button>
        </div>
      </div>

      <!-- Tags — Nodex §12 -->
      <div class="mt-4 px-3">
        <div class="flex items-center justify-between px-2 py-1 text-[11px] font-semibold tracking-widest text-[#64748B]">TAGS <button type="button" @click="onComingSoon" class="rounded px-1 text-[14px] leading-none hover:bg-[#EFF0FF] hover:text-[#6D28D9]">+</button></div>
        <div class="mt-1 space-y-0.5">
          <button v-for="t in tags" :key="t.name" @click="filterTag(t.name)" class="flex w-full items-center gap-2 rounded-[8px] px-2 py-1.5 text-left text-[14px] transition-colors" :class="route.query.tags === t.name ? 'bg-[#EFF0FF] text-[#6D28D9]' : 'text-[#334155] hover:bg-[#F8FAFC]'"><span class="h-2 w-2 rounded-full" :style="{background:t.color}"></span> <span class="flex-1">#{{ t.name }}</span> <span class="text-[12px] text-[#94A3B8]">{{ t.count }}</span></button>
        </div>
      </div>
    </nav>

    <!-- User area (Nodex §13) -->
    <div
      class="flex items-center gap-2.5 border-t border-[#E4E4E7] px-6 py-3 text-left"
      role="presentation"
    >
      <span class="flex h-7 w-7 items-center justify-center rounded-full bg-[#4F46E5] text-[13px] font-semibold text-[#FFFFFF]">
        {{ userInitial }}
      </span>
      <div class="leading-4">
        <div class="text-[15px] font-medium text-[#18181B]">{{ user.name }}</div>
        <div class="text-[12px] text-[#71717A]">{{ user.role }}</div>
      </div>
      <Icon name="chevron-down" size="16" class="ml-auto text-[#71717A]" />
    </div>
  </aside>
</template>

<script setup>
import { computed, ref, onMounted, watch } from "vue"
import { useRouter, useRoute } from "vue-router"
import Icon from "../ui/Icon.vue"
import { useToast } from "../../composables/useToast.js"
import { fetchPosts } from "../../services/api.js"

const props = defineProps({
  navItems: { type: Array, default: () => [] },
  open: { type: Boolean, default: true },
  user: {
    type: Object,
    default: () => ({ name: "Kevin", role: "Admin", initials: "K" }),
  },
})
const emit = defineEmits(["close", "navigate"])

const userInitial = computed(() => props.user.initials || props.user.name?.[0] || "U")

const primaryItems = computed(() => props.navItems.filter((i) => !i.secondary))
const secondaryItems = computed(() => props.navItems.filter((i) => i.secondary))

function navClass(item) {
  if (item.inert) return "text-[#71717A] opacity-60"
  if (item.active)
    return "bg-[#F5F3FF] text-[#6D28D9] hover:bg-[#EFF0FF]"
  return "text-[#71717A] hover:bg-[#F5F5F6] hover:text-[#18181B]"
}

const toast = useToast()
const router = useRouter()
const route = useRoute()

const folders = ref([
  { name: "Work", count: 0 },
  { name: "Study", count: 0 },
  { name: "Personal", count: 0 },
  { name: "Ideas", count: 0 },
  { name: "Projects", count: 0 },
])
const tags = ref([
  { name: "go", count: 0, color: "#86EFAC" },
  { name: "backend", count: 0, color: "#93C5FD" },
  { name: "vue", count: 0, color: "#F9A8D4" },
  { name: "database", count: 0, color: "#FDE68A" },
  { name: "life", count: 0, color: "#C4B5FD" },
])

async function refreshCounts(){
  try{
    const res = await fetchPosts({ limit: 100 })
    const all = res.data.data || []
    const folderMap = {}
    const tagMap = {}
    for(const p of all){
      if(p.folder) folderMap[p.folder] = (folderMap[p.folder]||0)+1
      if(p.tags) for(const t of p.tags.split(",").map(s=>s.trim()).filter(Boolean)) tagMap[t]=(tagMap[t]||0)+1
    }
    for(const f of folders.value) f.count = folderMap[f.name]||0
    for(const t of tags.value) t.count = tagMap[t.name]||0
  }catch{}
}
onMounted(refreshCounts)
watch(() => route.query, refreshCounts, { deep: true })
defineExpose({ refreshCounts })

function onComingSoon(){ toast("Coming soon — create folders/tags in editor") }
function filterFolder(name){
  const q = { ...route.query }
  if (q.folder === name) delete q.folder
  else q.folder = name
  router.push({ path: "/", query: q })
}
function filterTag(name){
  const q = { ...route.query }
  if (q.tags === name) delete q.tags
  else q.tags = name
  router.push({ path: "/", query: q })
}
</script>

<style scoped>
/* Keep inert buttons from showing focus ring on keyboard; still focusable=false via disabled. */
button:disabled {
  cursor: not-allowed;
}
</style>
