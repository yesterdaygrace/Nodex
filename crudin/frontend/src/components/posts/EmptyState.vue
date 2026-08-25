<!-- EmptyState.vue — variants per Design §19. -->
<template>
  <div class="py-16 text-center">
    <p class="text-[22px] font-semibold text-[#18181B]">{{ title }}</p>
    <p v-if="line1" class="mt-2 text-[#71717A]">{{ line1 }}</p>
    <p v-if="line2" class="text-[#71717A]">{{ line2 }}</p>
    <p v-if="sub" class="mt-1 text-[#71717A]">{{ sub }}</p>

    <button
      v-if="showButton"
      type="button"
      @click="emit('new-post')"
      class="mt-6 inline-flex items-center justify-center gap-2 rounded-[8px] bg-[#4F46E5] px-4 py-2 text-[15px] font-medium text-[#FFFFFF] hover:bg-[#4338CA] focus:outline-none focus:ring-2 focus:ring-[#4F46E5] focus:ring-offset-2 transition-colors"
    >
      <Icon name="plus" size="16" class="icon-pressable" />
      + New Post
    </button>
  </div>
</template>

<script setup>
import Icon from "../ui/Icon.vue"

const props = defineProps({
  title: { type: String, required: true },
  line1: { type: String, default: "" },
  line2: { type: String, default: "" },
  sub: { type: String, default: "" },
  showButton: { type: Boolean, default: true },
})
defineEmits(["new-post"])
</script>
