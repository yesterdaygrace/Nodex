<!-- PostActions.vue — Design §15/§16 + Nodex §29. Edit + Duplicate + Archive + Delete -->
<template>
  <div class="flex items-center gap-1">
    <button
      type="button"
      aria-label="Edit post"
      title="Edit post"
      @click="emit('edit', post)"
      class="rounded-[8px] p-1.5 text-[#71717A] hover:text-[#18181B] hover:bg-[#F5F5F6] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#4F46E5] focus-visible:ring-offset-2 transition-colors duration-150"
    >
      <Icon name="pencil" size="16" class="icon-pressable" />
    </button>
    <button
      type="button"
      aria-label="Duplicate post"
      title="Duplicate post"
      @click="emit('duplicate', post)"
      class="rounded-[8px] p-1.5 text-[#71717A] hover:text-[#18181B] hover:bg-[#F5F5F6] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#4F46E5] focus-visible:ring-offset-2 transition-colors duration-150"
    >
      <Icon name="copy" size="16" class="icon-pressable" />
    </button>
    <button
      v-if="post.status === 'archived'"
      type="button"
      aria-label="Unarchive post"
      title="Unarchive post"
      @click="emit('unarchive', post)"
      class="rounded-[8px] p-1.5 text-[#71717A] hover:text-[#6D28D9] hover:bg-[#F5F3FF] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#6D28D9] focus-visible:ring-offset-2 transition-colors duration-150"
    >
      <Icon name="archive" size="16" class="icon-pressable" />
    </button>
    <button
      v-else
      type="button"
      aria-label="Archive post"
      title="Archive post"
      @click="emit('archive', post)"
      class="rounded-[8px] p-1.5 text-[#71717A] hover:text-[#6D28D9] hover:bg-[#F5F3FF] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#6D28D9] focus-visible:ring-offset-2 transition-colors duration-150"
    >
      <Icon name="archive" size="16" class="icon-pressable" />
    </button>
    <button
      type="button"
      aria-label="Move post to trash"
      title="Move post to trash"
      @click="emit('delete', post)"
      class="rounded-[8px] p-1.5 text-[#71717A] hover:text-red-600 hover:bg-red-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-red-500 focus-visible:ring-offset-2 transition-colors duration-150"
    >
      <Icon name="trash-2" size="16" class="icon-pressable" />
    </button>
  </div>
</template>

<script setup>
import Icon from "../ui/Icon.vue"
const props = defineProps({
  post: { type: Object, required: true },
})
const emit = defineEmits(["edit", "delete", "duplicate", "archive", "unarchive"])
</script>
