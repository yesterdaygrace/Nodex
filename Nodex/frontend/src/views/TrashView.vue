<!-- TrashView.vue — Phase-2. Design §17/§19/§20/§21. Lists soft-deleted posts.
     Fetch: GET /api/posts/trashed?page=&limit=. Footer per card: Restore (neutral)
     + Delete Permanently (red, confirmation via PermanentDeleteDialog). -->
<template>
  <div class="flex flex-col gap-6">
    <!-- Page header (Design §17) -->
    <div class="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
      <div>
        <h1 class="text-[32px] font-semibold text-[#18181B]">Trash</h1>
        <p class="mt-1 text-[#71717A]">Posts you've moved to trash.</p>
      </div>
    </div>

    <!-- Search (Design §7) -->
    <div class="relative w-full sm:max-w-xl">
      <Icon
        name="search"
        size="16"
        class="absolute left-3 top-1/2 -translate-y-1/2 text-[#71717A]"
      />
      <input
        id="trash-search"
        type="search"
        :value="searchQuery"
        @input="searchQuery = $event.target.value"
        placeholder="Search trashed posts..."
        aria-label="Search trashed posts"
        class="w-full rounded-[8px] border border-[#E4E4E7] bg-[#FFFFFF] pl-10 pr-3 py-2 text-[15px] text-[#18181B] placeholder-[#71717A] outline-none transition-colors focus:border-[#4F46E5] focus:ring-2 focus:ring-[#4F46E5]/30"
      />
    </div>

    <!-- Error state (Design §21) -->
    <ErrorState v-if="error" :message="error" @retry="loadTrashed" />

    <!-- Loading skeletons (Design §20) -->
    <div
      v-else-if="loading"
      class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3 lg:gap-6"
    >
      <PostCardSkeleton v-for="i in skeletonCount" :key="`skel-${i}`" />
    </div>

    <!-- Empty trash (Design §19) -->
    <EmptyState
      v-else-if="filteredPosts.length === 0 && !searchQuery"
      title="Trash is empty"
      line1="Deleted posts will appear here."
      :show-button="false"
    />

    <!-- Empty: search no results -->
    <EmptyState
      v-else-if="filteredPosts.length === 0"
      title="No posts found"
      :sub="false"
      :show-button="false"
      line1="Try another search term."
      line2=""
    />

    <!-- Grid of trashed cards -->
    <div
      v-else
      class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3 lg:gap-6"
    >
      <TrashCard
        v-for="post in filteredPosts"
        :key="post.id"
        :post="post"
        @restore="restore"
        @delete-permanent="confirmPermanentDelete"
      />
    </div>

    <!-- Pagination (Design §18) -->
    <Pagination
      v-if="!loading && !error && total > 0"
      :page="page"
      :limit="LIMIT"
      :total="total"
      :total-pages="totalPages"
      @page="gotoPage"
    />

    <!-- Permanent-delete confirmation (Design §16/§17) -->
    <PermanentDeleteDialog
      v-if="permaPost"
      :post="permaPost"
      :show="!!permaPost"
      @cancel="cancelPermanentDelete"
      @confirm="confirmPermanentDeleteNow"
    />
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from "vue"
import Icon from "../components/ui/Icon.vue"
import PostCardSkeleton from "../components/posts/PostCardSkeleton.vue"
import Pagination from "../components/posts/Pagination.vue"
import TrashCard from "../components/trash/TrashCard.vue"
import PermanentDeleteDialog from "../components/trash/PermanentDeleteDialog.vue"
import EmptyState from "../components/posts/EmptyState.vue"
import ErrorState from "../components/posts/ErrorState.vue"
import { fetchTrashedPosts, restorePost, deletePermanent } from "../services/api.js"
import { useToast } from "../composables/useToast.js"

const LIMIT = 20
const posts = ref([])
const loading = ref(false)
const error = ref(null)
const page = ref(1)
const total = ref(0)
const totalPages = ref(1)
const searchQuery = ref("")
const permaPost = ref(null)

const toast = useToast()

const enrich = (items) =>
  items.map((p) => ({
    ...p,
    isFavorite: localStorage.getItem(`fav_${p.id}`) === "1",
  }))

async function loadTrashed() {
  loading.value = true
  error.value = null
  try {
    const res = await fetchTrashedPosts({ page: page.value, limit: LIMIT })
    const body = res.data
    if (!body.success && body.success !== undefined) {
      throw new Error(body.message || "Failed to load trashed posts")
    }
    posts.value = enrich(body.data || [])
    total.value = body.total || 0
    totalPages.value = body.total_pages || 1
  } catch (err) {
    error.value =
      err.response?.data?.message || err.message || "Failed to load trashed posts"
  } finally {
    loading.value = false
  }
}

const filteredPosts = computed(() => {
  const q = searchQuery.value?.trim().toLowerCase()
  if (!q) return posts.value
  return posts.value.filter(
    (p) =>
      (p.title || "").toLowerCase().includes(q) ||
      (p.content || "").toLowerCase().includes(q)
  )
})

function gotoPage(p) {
  const target = Math.max(1, Math.min(p, totalPages.value))
  if (target !== page.value) page.value = target
}

async function restore(post) {
  try {
    await restorePost(post.id)
    posts.value = posts.value.filter((p) => p.id !== post.id)
    total.value = Math.max(0, total.value - 1)
    toast("Post restored")
  } catch (err) {
    toast(err.response?.data?.message || err.message || "Failed to restore post", {
      error: true,
    })
  }
}

function confirmPermanentDelete(post) {
  permaPost.value = post
}
function cancelPermanentDelete() {
  permaPost.value = null
}
async function confirmPermanentDeleteNow() {
  if (!permaPost.value) return
  const post = permaPost.value
  permaPost.value = null
  try {
    await deletePermanent(post.id)
    posts.value = posts.value.filter((p) => p.id !== post.id)
    total.value = Math.max(0, total.value - 1)
    toast("Post permanently deleted")
  } catch (err) {
    toast(
      err.response?.data?.message || err.message || "Failed to permanently delete post",
      { error: true }
    )
  }
}

const skeletonCount = 6

watch(page, loadTrashed)
onMounted(loadTrashed)
</script>
