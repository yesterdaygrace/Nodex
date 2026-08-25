<!-- FavoriteButton.vue — Design §10. Star in top-right corner of the card (postioned by PostCard).
     Toggles localStorage `fav_<id>` and emits `favorite-toggled`. -->
<template>
  <button
    type="button"
    :aria-label="favorited ? 'Remove from favorites' : 'Add to favorites'"
    :title="favorited ? 'Remove from favorites' : 'Add to favorites'"
    @click="onToggle"
    class="rounded-full p-1.5 text-[#71717A] hover:text-[#4F46E5] hover:bg-[#EFF0FF] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#4F46E5] focus-visible:ring-offset-2 transition-colors duration-150"
  >
    <Icon name="star" size="17" :filled="favorited" class="icon-pressable" :class="favorited ? 'text-[#4F46E5]' : 'text-[#71717A]'" />
  </button>
</template>

<script setup>
import Icon from "../ui/Icon.vue"
import { useToast } from "../../composables/useToast.js"

const props = defineProps({
  postId: { type: [Number, String], required: true },
  favorited: { type: Boolean, default: false },
})
const emit = defineEmits(["favorite-toggled"])

const toast = useToast()
const KEY = `fav_${props.postId}`

function onToggle() {
  const next = !props.favorited
  if (next) localStorage.setItem(KEY, "1")
  else localStorage.removeItem(KEY)
  emit("favorite-toggled", { postId: props.postId, favorited: next })
  toast(next ? "Added to favorites" : "Removed from favorites")
}
</script>
