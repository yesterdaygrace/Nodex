<!-- ArchiveView.vue — Nodex §35. Archive: Keep old notes without deleting them. -->
<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
      <div>
        <h1 class="text-[32px] font-bold text-[#111827]">Archive</h1>
        <p class="mt-1 text-[15px] text-[#64748B]">Keep old notes without deleting them.</p>
      </div>
    </div>

    <PostFilters v-model:modelValue="searchQuery" v-model:statusFilter="statusFilter" v-model:view="view" />

    <ErrorState v-if="error" :message="error" @retry="loadPosts" />
    <PostGrid v-else-if="loading" :loading="true" />
    <EmptyState
      v-else-if="filteredPosts.length === 0 && !searchQuery"
      title="No archived notes"
      line1="Archived notes will appear here."
      :show-button="false"
    />
    <EmptyState
      v-else-if="filteredPosts.length === 0"
      title="No notes found"
      :sub="false"
      :show-button="false"
      line1="Try another search term or"
      line2="change your filters."
    />
    <PostGrid v-else :posts="filteredPosts" :on-edit="openEdit" :on-delete="handleDelete" :on-favorite="handleFavoriteToggled" :on-archive="handleUnarchive" />

    <Pagination v-if="!loading && !error && total>0" :page="page" :limit="LIMIT" :total="total" :total-pages="totalPages" @page="gotoPage" />

    <DeletePostConfirmDialog v-if="confirmPost" :post="confirmPost" :show="!!confirmPost" @cancel="cancelDelete" @confirm="confirmDelete" />
    <PostFormDialog v-if="showForm" :show="showForm" :post="editingPost" @submit="submitForm" @cancel="closeForm" />
  </div>
</template>

<script setup>
import PostFilters from "../components/posts/PostFilters.vue"
import PostGrid from "../components/posts/PostGrid.vue"
import Pagination from "../components/posts/Pagination.vue"
import PostFormDialog from "../components/posts/PostFormDialog.vue"
import DeletePostConfirmDialog from "../components/posts/DeletePostConfirmDialog.vue"
import EmptyState from "../components/posts/EmptyState.vue"
import ErrorState from "../components/posts/ErrorState.vue"
import { usePostsList } from "../composables/usePostsList.js"
import { useToast } from "../composables/useToast.js"
import { unarchivePost } from "../services/api.js"

const {
  searchQuery, statusFilter, view, filteredPosts, loading, error, page, total, totalPages, LIMIT,
  confirmPost, showForm, editingPost, openCreate, openEdit, closeForm, submitForm,
  handleDelete, cancelDelete, confirmDelete, handleFavoriteToggled, gotoPage, loadPosts, posts
} = usePostsList("archived")

const toast = useToast()
async function handleUnarchive(post){
  try{ await unarchivePost(post.id); toast("Note restored"); loadPosts() }catch(e){ toast(e.response?.data?.message||e.message,{error:true}) }
}
</script>
