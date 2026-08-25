<!-- TrashCard.vue — Design §17. Trashed post card: title, description, "Deleted X ago"
     + footer text actions Restore (neutral) / Delete Permanently (red). Reuses the
     bento card style for visual consistency (Design §34). No favorite star in trash;
     destructive actions use explicit text labels (§29). -->
<template>
  <article
    class="flex h-full flex-col gap-3 rounded-[16px] bg-[#FFFFFF] p-6 pb-4 ring-1 ring-[#E4E4E7] shadow-sm transition-shadow duration-200 hover:shadow-md"
  >
    <header>
      <h3
        class="text-[20px] font-semibold text-[#18181B]"
        style="display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden"
        :title="post.title"
      >
        {{ post.title }}
      </h3>
    </header>

    <p
      class="text-[15px] leading-relaxed text-[#71717A]"
      style="display: -webkit-box; -webkit-line-clamp: 3; -webkit-box-orient: vertical; overflow: hidden"
      :title="post.content"
    >
      {{ description }}
    </p>

    <footer class="mt-auto pt-3 ring-t-1 ring-[#E4E4E7] flex items-center justify-between">
      <div class="flex items-center gap-1.5 text-[13px] text-[#71717A]">
        <Icon name="trash-2" size="13" />
        <span>Deleted {{ deletedLabel }}</span>
      </div>
      <div class="flex items-center gap-1.5">
        <button
          type="button"
          aria-label="Restore post"
          title="Restore post"
          @click="$emit('restore', post)"
          class="rounded-[8px] border border-[#E4E4E7] bg-[#FFFFFF] px-2.5 py-1 text-[15px] font-medium text-[#18181B] transition-colors hover:bg-[#F5F5F6] focus:outline-none focus:ring-2 focus:ring-[#4F46E5]/30 focus:ring-offset-2"
        >
          Restore
        </button>
        <button
          type="button"
          aria-label="Delete permanently"
          title="Delete permanently"
          @click="$emit('delete-permanent', post)"
          class="rounded-[8px] border border-transparent bg-red-600 px-2.5 py-1 text-[15px] font-medium text-[#FFFFFF] transition-colors hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500/30 focus:ring-offset-2"
        >
          Delete Permanently
        </button>
      </div>
    </footer>
  </article>
</template>

<script setup>
import { computed } from "vue"
import Icon from "../ui/Icon.vue"
import { timeAgo } from "../../utils/timeAgo.js"

const props = defineProps({
  post: { type: Object, required: true },
})
defineEmits(["restore", "delete-permanent"])

const description = computed(() => {
  const c = props.post.content || ""
  return c.length <= 160 ? c : c.slice(0, 160).replace(/\s+$/, "") + "..."
})

// Trashed rows carry deleted_at (snake_case matches the posts JSON envelope).
const deletedLabel = computed(() =>
  timeAgo(props.post.deleted_at || props.post.deletedAt)
)
</script>
