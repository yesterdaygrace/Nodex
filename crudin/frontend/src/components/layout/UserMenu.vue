<!-- UserMenu.vue — Design §4.19 / §5 user dropdown.
     Dropdown items are inert placeholders in Phase-1 (no auth/router). -->
<template>
  <div class="relative" ref="menuRef">
    <button
      type="button"
      @click="open = !open"
      :aria-expanded="open"
      class="inline-flex items-center gap-1.5 rounded-[8px] px-2 py-1 text-[15px] text-[#18181B] hover:bg-[#F5F5F6] focus:outline-none focus:ring-2 focus:ring-[#4F46E5]/30"
      aria-label="Open user menu"
    >
      <span class="flex h-6 w-6 items-center justify-center rounded-full bg-[#4F46E5] text-[11px] font-semibold text-[#FFFFFF]">
        {{ userInitial }}
      </span>
      <span class="hidden sm:inline">{{ user.name }}</span>
      <Icon name="chevron-down" size="14" class="text-[#71717A] icon-pressable" />
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
        v-show="open"
        class="absolute right-0 mt-1 w-40 overflow-hidden rounded-[8px] border border-[#E4E4E7] bg-[#FFFFFF] shadow-lg"
      >
        <ul class="py-1 text-[15px]">
          <li>
            <a
              href="#"
              class="block px-3 py-1.5 text-[#18181B] hover:bg-[#F5F3FF] hover:text-[#4F46E5]"
              @click="noop"
              >Profile</a
            >
          </li>
          <li>
            <a
              href="#"
              class="block px-3 py-1.5 text-[#18181B] hover:bg-[#F5F3FF] hover:text-[#4F46E5]"
              @click="noop"
              >Settings</a
            >
          </li>
          <li>
            <a
              href="#"
              class="block px-3 py-1.5 text-[#18181B] hover:bg-[#F5F3FF] hover:text-[#4F46E5]"
              @click="noop"
              >Logout</a
            >
          </li>
        </ul>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from "vue"
import Icon from "../ui/Icon.vue"

const props = defineProps({
  user: {
    type: Object,
    default: () => ({ name: "Kevin", role: "Admin", initials: "K" }),
  },
})

const open = ref(false)
const menuRef = ref(null)
const userInitial = computed(() => props.user.initials || props.user.name?.[0] || "U")

function noop(e) {
  // Phase-1: auth/router not wired.
  e.preventDefault()
}

function onOutside(e) {
  if (open.value && menuRef.value && !menuRef.value.contains(e.target)) {
    open.value = false
  }
}
onMounted(() => {
  document.addEventListener("click", onOutside)
  document.addEventListener("keydown", (e) => {
    if (e.key === "Escape") open.value = false
  })
})
onBeforeUnmount(() => {
  document.removeEventListener("click", onOutside)
})
</script>
