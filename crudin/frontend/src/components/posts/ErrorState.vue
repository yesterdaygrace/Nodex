<!-- ErrorState.vue — Design §21. Non-technical, retryable. -->
<template>
  <div class="py-16 text-center">
    <p class="text-[22px] font-semibold text-[#18181B]">Something went wrong</p>
    <p class="mt-2 text-[#71717A]">We couldn't load your posts.</p>
    <button
      type="button"
      @click="emit('retry')"
      class="mt-6 inline-flex items-center justify-center gap-2 rounded-[8px] bg-[#4F46E5] px-4 py-2 text-[15px] font-medium text-[#FFFFFF] hover:bg-[#4338CA] focus:outline-none focus:ring-2 focus:ring-[#4F46E5] focus:ring-offset-2 transition-colors"
    >
      <Icon name="check" size="16" />
      Try Again
    </button>
  </div>
</template>

<script setup>
import Icon from "../ui/Icon.vue"
defineEmits(["retry"])
</script>
