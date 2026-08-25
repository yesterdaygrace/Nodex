<!-- PostCard.vue — Design §9/§11/§12/§14. Bento card. -->
<template>
  <article
    class="relative flex h-full flex-col gap-3 rounded-[16px] bg-[#FFFFFF] p-6 pb-4 ring-1 ring-[#E4E4E7] shadow-sm transition-shadow duration-200 hover:shadow-md"
  >
    <!-- Favorite star: top-right corner (only once per card) -->
    <FavoriteButton
      :post-id="post.id"
      :favorited="post.isFavorite"
      class="absolute top-4 right-4"
      @favorite-toggled="onFavorite"
    />

    <header>
      <RouterLink :to="`/notes/${post.id}`" class="block hover:underline">
        <h3
          class="text-[20px] font-semibold text-[#18181B] hover:text-[#6D28D9]"
          style="display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden"
          :title="post.title"
        >
          {{ post.title }}
        </h3>
      </RouterLink>
    </header>

    <p
      class="text-[15px] leading-relaxed text-[#71717A]"
      style="display: -webkit-box; -webkit-line-clamp: 3; -webkit-box-orient: vertical; overflow: hidden"
      :title="post.content"
    >
      {{ description }}
    </p>

    <!-- Tags — Nodex §28 -->
    <div v-if="tagList.length" class="flex flex-wrap gap-1.5">
      <span v-for="tag in tagList" :key="tag" class="rounded-full bg-[#EFF6FF] px-2 py-0.5 text-[12px] font-medium text-[#1D4ED8]">#{{ tag }}</span>
    </div>

    <PostStatusBadge :status="post.status" class="self-start" />

    <footer class="mt-auto pt-3 ring-t-1 ring-[#E4E4E7] flex items-center justify-between">
      <div class="flex items-center gap-1.5 text-[13px] text-[#71717A]">
        <Icon name="clock" size="13" />
        <time :datetime="post.updated_at || post.created_at">{{ timeAgo(post.updated_at || post.created_at) }}</time>
      </div>
      <PostActions :post="post" @edit="emitEdit" @delete="emitDelete" @duplicate="emitDuplicate" @archive="emitArchive" @unarchive="emitUnarchive" />
    </footer>
  </article>
</template>

<script setup>
import { computed } from "vue"
import Icon from "../ui/Icon.vue"
import { timeAgo } from "../../utils/timeAgo.js"
import PostStatusBadge from "./PostStatusBadge.vue"
import PostActions from "./PostActions.vue"
import FavoriteButton from "./FavoriteButton.vue"

const props = defineProps({
  post: { type: Object, required: true },
})
const emit = defineEmits(["edit", "delete", "favorite-toggled", "duplicate", "archive", "unarchive"])

const description = computed(() => {
  const c = props.post.content || ""
  return c.length <= 160 ? c : c.slice(0, 160).replace(/\s+$/, "") + "..."
})

const tagList = computed(() => {
  if (!props.post.tags) return []
  return props.post.tags.split(",").map(t=>t.trim()).filter(Boolean)
})

function emitEdit() { emit("edit", props.post) }
function emitDelete() { emit("delete", props.post) }
function emitDuplicate() { emit("duplicate", props.post) }
function emitArchive() { emit("archive", props.post) }
function emitUnarchive() { emit("unarchive", props.post) }
function onFavorite() { emit("favorite-toggled", props.post.id) }
</script>
