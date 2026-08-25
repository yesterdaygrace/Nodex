<!-- PostsView.vue — Design §6/§7/§8/§18/§19/§20/§21/§22. Phase-2 live view.
     Fetches GET /api/posts?status=published&page=&limit=. Owns search/pagination/
     delete flow + create/edit via PostFormDialog. Reuses shared post components. -->
<template>
  <div class="flex flex-col gap-6">
    <!-- Page header (Design §6.1) -->
    <div class="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
      <div>
        <h1 class="text-[32px] font-bold text-[#111827]">All Notes</h1>
        <p class="mt-1 text-[15px] text-[#64748B]">Your ideas, organized and always with you.</p>
      </div>
      <button
        type="button"
        @click="openCreate"
        class="inline-flex items-center justify-center gap-2 rounded-[8px] bg-[#6D28D9] px-4 py-2 text-[15px] font-medium text-[#FFFFFF] hover:bg-[#5B21B6] focus:outline-none focus:ring-2 focus:ring-[#6D28D9] focus:ring-offset-2 transition-colors"
      >
        <Icon name="plus" size="16" class="icon-pressable" />
        + New Note
      </button>
    </div>

    <div v-if="favoritesBanner" class="flex items-center justify-between rounded-[8px] border border-[#E2E8F0] bg-[#EFF6FF] px-4 py-3 text-[14px] text-[#1E40AF]">
      <span>★ Favorites are automatically pinned to the top.</span>
      <button type="button" @click="dismissBanner" class="ml-4 text-[#64748B] hover:text-[#111827]">×</button>
    </div>

    <!-- Search + Status + View (Design §7) -->
    <PostFilters
      v-model:modelValue="searchQuery"
      v-model:statusFilter="statusFilter"
      v-model:view="view"
    />

    <QuickTabs :counts="tabCounts" />

    <!-- Error state (Design §21) -->
    <ErrorState v-if="error" :message="error" @retry="loadPosts" />

    <!-- Loading skeletons (Design §20) -->
    <PostGrid v-else-if="loading" :loading="true" />

    <!-- Empty: no posts -->
    <EmptyState
      v-else-if="filteredPosts.length === 0 && !searchQuery"
      title="No posts yet"
      line1="Create your first post"
      line2="to get started."
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

    <!-- Grid / List -->
    <PostGrid
      v-else-if="view === 'grid'"
      :posts="filteredPosts"
      :on-edit="openEdit"
      :on-delete="handleDelete"
      :on-favorite="handleFavoriteToggled"
      :on-duplicate="handleDuplicate"
      :on-archive="handleArchive"
      :on-unarchive="handleUnarchive"
    />
    <NoteList
      v-else
      :posts="filteredPosts"
      :on-edit="openEdit"
      :on-delete="handleDelete"
      :on-favorite="handleFavoriteToggled"
      :on-duplicate="handleDuplicate"
      :on-archive="handleArchive"
      :on-unarchive="handleUnarchive"
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
import { ref, watch, computed, onMounted } from "vue"
import { useRoute } from "vue-router"
import { usePostsList } from "../composables/usePostsList.js"
import QuickTabs from "../components/posts/QuickTabs.vue"
import NoteList from "../components/posts/NoteList.vue"
import { fetchPosts } from "../services/api.js"

// All Notes shows every active note (published + draft), archived excluded by backend when no status.
const route = useRoute()
const {
  posts,
  searchQuery,
  statusFilter,
  folderFilter,
  tagFilter,
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
  handleDuplicate,
  handleArchive,
  handleUnarchive,
  gotoPage,
  loadPosts,
} = usePostsList("")

watch(() => route.query.folder, v => { folderFilter.value = v || ""; }, { immediate: true })
watch(() => route.query.tags, v => { tagFilter.value = v || ""; }, { immediate: true })

const favoritesBanner = ref(localStorage.getItem("nodex_banner_dismissed") !== "1")
watch(filteredPosts, (list) => {
  if (!list.some(p=>p.isFavorite)) favoritesBanner.value = false
}, { immediate: true })
function dismissBanner(){ favoritesBanner.value=false; localStorage.setItem("nodex_banner_dismissed","1") }

const viewPersistKey = "nodex_view"
if (localStorage.getItem(viewPersistKey)) view.value = localStorage.getItem(viewPersistKey)
watch(view, v=> localStorage.setItem(viewPersistKey, v))

const archiveCount = ref(0)
const trashCount = ref(0)
const tabCounts = computed(() => ({
  all: total.value,
  favorites: posts.value.filter(p=>p.isFavorite).length,
  archive: archiveCount.value,
  trash: trashCount.value,
}))
onMounted(async ()=>{
  try{ const r = await fetchPosts({ status:"archived", limit:1 }); archiveCount.value = r.data.total }catch{}
  try{ const { fetchTrashedPosts } = await import("../services/api.js"); const r2 = await fetchTrashedPosts({ limit:1 }); trashCount.value = r2.data.total }catch{}
})
</script>
