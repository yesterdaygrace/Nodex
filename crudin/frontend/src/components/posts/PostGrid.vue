<!-- PostGrid.vue — Design §8. Responsive bento grid. Props/contract:
     PostGrid(posts, onEdit, onDelete, onFavorite). -->
<template>
  <div
    v-if="loading"
    class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3 lg:gap-6"
  >
    <PostCardSkeleton v-for="i in skeletonCount" :key="`skel-${i}`" />
  </div>

  <div
    v-else
    class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3 lg:gap-6"
  >
    <PostCard
      v-for="post in posts"
      :key="post.id"
      :post="post"
      @edit="onEdit"
      @delete="onDelete"
      @favorite-toggled="onFavorite"
      @duplicate="onDuplicate"
      @archive="onArchive"
      @unarchive="onUnarchive"
    />
  </div>
</template>

<script setup>
import PostCard from "./PostCard.vue"
import PostCardSkeleton from "./PostCardSkeleton.vue"

const props = defineProps({
  posts: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
  skeletonCount: { type: Number, default: 6 },
  onEdit: { type: Function, default: () => {} },
  onDelete: { type: Function, default: () => {} },
  onFavorite: { type: Function, default: () => {} },
  onDuplicate: { type: Function, default: () => {} },
  onArchive: { type: Function, default: () => {} },
  onUnarchive: { type: Function, default: () => {} },
})
</script>
