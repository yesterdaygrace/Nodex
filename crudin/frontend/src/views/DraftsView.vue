<!-- DraftsView.vue — Phase-2. Dedicated Drafts page (Design §7/§8/§9/§18/§19).
     Reuses the shared post components; fetches GET /api/posts?status=draft. -->
<template>
  <div class="flex flex-col gap-6">
    <!-- Page header (Design §6.1) -->
    <div class="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
      <div>
        <h1 class="text-[32px] font-semibold text-[#18181B]">Drafts</h1>
        <p class="mt-1 text-[#71717A]">Drafts you've saved but not published yet.</p>
      </div>
      <button
        type="button"
        @click="openCreate"
        class="inline-flex items-center justify-center gap-2 rounded-[8px] bg-[#4F46E5] px-4 py-2 text-[15px] font-medium text-[#FFFFFF] hover:bg-[#4338CA] focus:outline-none focus:ring-2 focus:ring-[#4F46E5] focus:ring-offset-2 transition-colors"
      >
        <Icon name="plus" size="16" class="icon-pressable" />
        + New Post
      </button>
    </div>

    <!-- Search + Status + View (reuse PostFilters) -->
    <PostFilters
      v-model:modelValue="searchQuery"
      v-model:statusFilter="statusFilter"
      v-model:view="view"
    />

    <!-- Error state (Design §21) -->
    <ErrorState v-if="error" :message="error" @retry="loadPosts" />

    <!-- Loading skeletons (Design §20) -->
    <PostGrid v-else-if="loading" :loading="true" />

    <!-- Empty: no drafts -->
    <EmptyState
      v-else-if="filteredPosts.length === 0 && !searchQuery"
      title="No drafts yet"
      line1="Start creating"
      line2="content and save it as a draft."
      @new-post="openCreate"
    />

    <!-- Empty: search/filter no results -->
    <EmptyState
      v-else-if="filteredPosts.length === 0"
      title="No posts found"
      :sub="false"
      :show-button="false"
      line1="Try another search term or"
      line2="change your filters."
    />

    <!-- Grid -->
    <PostGrid
      v-else
      :posts="filteredPosts"
      :on-edit="openEdit"
      :on-delete="handleDelete"
      :on-favorite="handleFavoriteToggled"
    />

    <!-- Pagination (Design §18) -->
    <Pagination
      v-if="!loading && !error && total > 0"
      :page="page"
      :limit="LIMIT"
      :total="total"
      :total-pages="totalPages"
      @page="gotoPage"
    />

    <!-- Delete confirmation (Design §16) -->
    <DeletePostConfirmDialog
      v-if="confirmPost"
      :post="confirmPost"
      :show="!!confirmPost"
      @cancel="cancelDelete"
      @confirm="confirmDelete"
    />

    <!-- Create / Edit form (Design §6.2) -->
    <PostFormDialog
      v-if="showForm"
      :show="showForm"
      :post="editingPost"
      @submit="submitForm"
      @cancel="closeForm"
    />
  </div>
</template>

<script setup>
import Icon from "../components/ui/Icon.vue"
import PostFilters from "../components/posts/PostFilters.vue"
import PostGrid from "../components/posts/PostGrid.vue"
import Pagination from "../components/posts/Pagination.vue"
import PostFormDialog from "../components/posts/PostFormDialog.vue"
import DeletePostConfirmDialog from "../components/posts/DeletePostConfirmDialog.vue"
import EmptyState from "../components/posts/EmptyState.vue"
import ErrorState from "../components/posts/ErrorState.vue"
import { usePostsList } from "../composables/usePostsList.js"

// Drafts page defaults to the draft status (Design §32/§33).
const {
  searchQuery,
  statusFilter,
  view,
  filteredPosts,
  loading,
  error,
  page,
  total,
  totalPages,
  LIMIT,
  confirmPost,
  showForm,
  editingPost,
  openCreate,
  openEdit,
  closeForm,
  submitForm,
  handleDelete,
  cancelDelete,
  confirmDelete,
  handleFavoriteToggled,
  gotoPage,
  loadPosts,
} = usePostsList("draft")
</script>
