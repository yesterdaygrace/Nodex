<!-- App.vue — Phase-2 layout shell (Design §4/§5). Owns sidebar drawer state + global toasts.
     Active nav is driven by the current route via vue-router. Post list state lives
     in the routed views (PostsView/DraftsView/TrashView), not here. -->
<template>
  <div class="min-h-screen bg-[#FAFAFA] font-sans text-[#18181B]">
    <!-- Sidebar (fixed on desktop, drawer on mobile) -->
    <AppSidebar
      :nav-items="navItems"
      :open="sidebarOpen"
      :user="user"
      @close="sidebarOpen = false"
      @navigate="onNavigate"
    />

    <!-- Main content offset for desktop sidebar -->
    <div class="flex flex-col lg:ml-64">
      <AppHeader :sidebar-open="sidebarOpen" :title="pageTitle" :user="user" @menu="sidebarOpen = true" />

      <main class="flex-1 overflow-y-auto">
        <div class="mx-auto max-w-5xl p-6 lg:p-8">
          <RouterView />
        </div>
      </main>
    </div>

    <!-- Toast stack (Design §22) -->
    <transition-group
      tag="div"
      class="fixed top-4 right-4 z-50 flex flex-col gap-2"
      enter-active-class="transition ease-out duration-200"
      enter-from-class="opacity-0 translate-x-4"
      enter-to-class="opacity-100 translate-x-0"
      leave-active-class="transition ease-in duration-150"
      leave-from-class="opacity-100 translate-x-0"
      leave-to-class="opacity-0 translate-x-4"
    >
      <Toast
        v-for="t in toasts"
        :key="t.id"
        :message="t.message"
        :error="t.error"
        class="w-72"
      />
    </transition-group>
  </div>
</template>

<script setup>
import { ref, computed, provide } from "vue"
import { useRoute } from "vue-router"
import AppSidebar from "./components/layout/AppSidebar.vue"
import AppHeader from "./components/layout/AppHeader.vue"
import Toast from "./components/ui/Toast.vue"
import { toastSymbol } from "./composables/useToast.js"

const sidebarOpen = ref(false)
const toasts = ref([])
let toastCounter = 0

const route = useRoute()

// Route-aware navigation (Nodex §8). All Notes / Favorites / Archive / Trash.
const navItems = [
  { name: "All Notes", icon: "file-text", to: "/", active: route.path === "/" },
  { name: "Favorites", icon: "star", to: "/favorites", active: route.path === "/favorites" },
  { name: "Archive", icon: "archive", to: "/archive", active: route.path === "/archive" },
  { name: "Trash", icon: "trash-2", to: "/trash", secondary: true, active: route.path === "/trash" },
]

const user = { name: "Kevin", role: "Free Plan", initials: "K" }

const toast = (message, opts = {}) => {
  const id = ++toastCounter
  toasts.value.push({ id, message, error: !!opts.error })
  setTimeout(() => {
    const idx = toasts.value.findIndex((t) => t.id === id)
    if (idx !== -1) toasts.value.splice(idx, 1)
  }, 2400)
}

provide(toastSymbol, toast)

// Close mobile drawer after navigating to a live route.
function onNavigate(item) {
  if (!item.inert) sidebarOpen.value = false
}

const pageTitle = computed(() => {
  const map = { "/": "All Notes", "/favorites": "Favorites", "/archive": "Archive", "/trash": "Trash", "/drafts": "Archive" }
  return map[route.path] || "All Notes"
})
</script>
