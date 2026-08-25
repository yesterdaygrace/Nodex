<!-- NoteList.vue — Nodex §20 List view. Single column, denser. -->
<template>
  <div class="flex flex-col gap-3">
    <article
      v-for="post in posts"
      :key="post.id"
      class="flex flex-col gap-2 rounded-[12px] bg-white p-4 ring-1 ring-[#E2E8F0] hover:shadow-sm transition-shadow"
    >
      <div class="flex items-start justify-between gap-3">
        <h3 class="text-[16px] font-semibold text-[#111827] line-clamp-1">{{ post.title }}</h3>
        <FavoriteButton :post-id="post.id" :favorited="post.isFavorite" @favorite-toggled="onFavorite" />
      </div>
      <p class="text-[14px] text-[#64748B] line-clamp-2">{{ post.content }}</p>
      <div v-if="post.tags" class="flex flex-wrap gap-1.5">
        <span v-for="tag in post.tags.split(',').map(t=>t.trim()).filter(Boolean)" :key="tag" class="rounded-full bg-[#EFF6FF] px-2 py-0.5 text-[11px] text-[#1D4ED8]">#{{ tag }}</span>
      </div>
      <div class="flex items-center justify-between pt-2 border-t border-[#F1F5F9]">
        <span class="text-[12px] text-[#94A3B8]">{{ post.updated_at || post.created_at }}</span>
        <PostActions :post="post" @edit="onEdit" @delete="onDelete" @duplicate="onDuplicate" @archive="onArchive" @unarchive="onUnarchive" />
      </div>
    </article>
  </div>
</template>

<script setup>
import FavoriteButton from "./FavoriteButton.vue"
import PostActions from "./PostActions.vue"
defineProps({
  posts: { type: Array, default: () => [] },
  onEdit: { type: Function, default: () => {} },
  onDelete: { type: Function, default: () => {} },
  onFavorite: { type: Function, default: () => {} },
  onDuplicate: { type: Function, default: () => {} },
  onArchive: { type: Function, default: () => {} },
  onUnarchive: { type: Function, default: () => {} },
})
</script>
