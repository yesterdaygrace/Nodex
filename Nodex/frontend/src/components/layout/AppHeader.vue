<!-- AppHeader.vue — Design §5. Top header; mobile menu + right-side controls. -->
<template>
  <header class="flex h-14 items-center justify-between border-b border-[#E4E4E7] bg-[#FFFFFF] px-4 sm:px-6 lg:px-8">
    <div class="flex items-center gap-2">
      <!-- Mobile menu toggle -->
      <button
        type="button"
        @click="emit('menu')"
        class="lg:hidden inline-flex h-8 w-8 items-center justify-center rounded-[8px] text-[#71717A] hover:bg-[#F5F5F6] hover:text-[#18181B] focus:outline-none focus:ring-2 focus:ring-[#4F46E5]/30"
        aria-label="Open menu"
        aria-controls="app-sidebar"
        :aria-expanded="sidebarOpen"
      >
        <Icon name="menu" size="18" class="icon-pressable" />
      </button>
      <h1 class="text-[17px] font-medium text-[#18181B]">{{ title }}</h1>
    </div>

    <nav class="flex items-center gap-1">
      <!-- Inert controls (Phase-1): search bell are decorative; user dropdown is live. -->
      <span aria-hidden="true" title="Search" class="inline-flex h-8 w-8 items-center justify-center rounded-[8px] text-[#71717A]">
        <Icon name="search" size="18" />
      </span>
      <span aria-hidden="true" title="Notifications" class="inline-flex h-8 w-8 items-center justify-center rounded-[8px] text-[#71717A]">
        <Icon name="bell" size="18" />
      </span>
      <UserMenu :user="user" />
    </nav>
  </header>
</template>

<script setup>
import Icon from "../ui/Icon.vue"
import UserMenu from "./UserMenu.vue"

defineProps({
  sidebarOpen: { type: Boolean, default: false },
  title: { type: String, default: "Posts" },
  user: {
    type: Object,
    default: () => ({ name: "Kevin", role: "Admin", initials: "K" }),
  },
})
defineEmits(["menu"])
</script>
