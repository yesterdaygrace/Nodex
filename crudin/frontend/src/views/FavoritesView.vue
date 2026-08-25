<!-- FavoritesView.vue — Nodex §36. Favorites: Your most important notes. -->
<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
      <div>
        <h1 class="text-[32px] font-bold text-[#111827]">Favorites</h1>
        <p class="mt-1 text-[15px] text-[#64748B]">Your most important notes.</p>
      </div>
    </div>

    <div v-if="favoritesBanner" class="flex items-center justify-between rounded-[8px] border border-[#E2E8F0] bg-[#EFF6FF] px-4 py-3 text-[14px] text-[#1E40AF]">
      <span>★ Favorites are automatically pinned to the top.</span>
      <button type="button" @click="dismissBanner" class="ml-4 text-[#64748B] hover:text-[#111827]">×</button>
    </div>

    <PostFilters v-model:modelValue="searchQuery" v-model:statusFilter="statusFilter" v-model:view="view" />

    <ErrorState v-if="error" :message="error" @retry="loadPosts" />
    <PostGrid v-else-if="loading" :loading="true" />
    <EmptyState
      v-else-if="filteredFavorites.length === 0 && !searchQuery"
      title="No favorite notes"
      line1="Star important notes to keep them at the top."
      :show-button="false"
    />
    <EmptyState
      v-else-if="filteredFavorites.length === 0"
      title="No notes found"
      :sub="false"
      :show-button="false"
      line1="Try another search term or"
      line2="change your filters."
    />
    <PostGrid v-else :posts="filteredFavorites" :on-edit="openEdit" :on-delete="handleDelete" :on-favorite="handleFavoriteToggled" />

    <Pagination v-if="!loading && !error && filteredFavorites.length>0" :page="page" :limit="LIMIT" :total="filteredFavorites.length" :total-pages="1" @page="gotoPage" />

    <DeletePostConfirmDialog v-if="confirmPost" :post="confirmPost" :show="!!confirmPost" @cancel="cancelDelete" @confirm="confirmDelete" />
    <PostFormDialog v-if="showForm" :show="showForm" :post="editingPost" @submit="submitForm" @cancel="closeForm" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue"
import PostFilters from "../components/posts/PostFilters.vue"
import PostGrid from "../components/posts/PostGrid.vue"
import Pagination from "../components/posts/Pagination.vue"
import PostFormDialog from "../components/posts/PostFormDialog.vue"
import DeletePostConfirmDialog from "../components/posts/DeletePostConfirmDialog.vue"
import EmptyState from "../components/posts/EmptyState.vue"
import ErrorState from "../components/posts/ErrorState.vue"
import { usePostsList } from "../composables/usePostsList.js"

const {
  searchQuery, statusFilter, view, filteredPosts, loading, error, page, total, totalPages, LIMIT,
  confirmPost, showForm, editingPost, openCreate, openEdit, closeForm, submitForm,
  handleDelete, cancelDelete, confirmDelete, handleFavoriteToggled, gotoPage, loadPosts,
} = usePostsList("")

const favoritesBanner = ref(localStorage.getItem("nodex_banner_dismissed") !== "1" && filteredPosts.value.some(p=>p.isFavorite))
function dismissBanner(){ favoritesBanner.value=false; localStorage.setItem("nodex_banner_dismissed","1") }

const filteredFavorites = computed(() => {
  const q = searchQuery.value?.trim().toLowerCase()
  let list = filteredPosts.value.filter(p=>p.isFavorite)
  if(q) list = list.filter(p=> (p.title||"").toLowerCase().includes(q) || (p.content||"").toLowerCase().includes(q))
  return list
})
</script>
