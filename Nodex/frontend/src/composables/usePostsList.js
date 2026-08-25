// src/composables/usePostsList.js — shared data layer for PostsView & DraftsView.
// Encapsulates: fetch (status-filtered), search, pagination, favorite toggle,
// soft-delete (trash) flow, and create/edit form state. Views stay presentational.
import { ref, computed, watch, onMounted } from "vue"
import { fetchPosts, deletePost, createPost, updatePost, archivePost, unarchivePost } from "../services/api.js"
import { useToast } from "./useToast.js"

export const LIMIT = 20

export function usePostsList(defaultStatus = "") {
  const posts = ref([])
  const loading = ref(false)
  const error = ref(null)
  const page = ref(1)
  const total = ref(0)
  const totalPages = ref(1)
  const searchQuery = ref("")
  const statusFilter = ref(defaultStatus || "")
  const folderFilter = ref("")
  const tagFilter = ref("")
  const view = ref("grid")
  const confirmPost = ref(null)

  // Create/edit form state.
  const showForm = ref(false)
  const editingPost = ref(null)

  const toast = useToast()

  function enrich(items) {
    return items.map((p) => ({
      ...p,
      isFavorite: localStorage.getItem(`fav_${p.id}`) === "1",
    }))
  }

  async function loadPosts() {
    loading.value = true
    error.value = null
    try {
      const res = await fetchPosts({
        page: page.value,
        limit: LIMIT,
        status: statusFilter.value || undefined,
        folder: folderFilter.value || undefined,
        tags: tagFilter.value || undefined,
      })
      const body = res.data
      if (!body.success && body.success !== undefined) {
        throw new Error(body.message || "Failed to load posts")
      }
      posts.value = enrich(body.data || [])
      total.value = body.total || 0
      totalPages.value = body.total_pages || 1
    } catch (err) {
      error.value =
        err.response?.data?.message || err.message || "Failed to load posts"
    } finally {
      loading.value = false
    }
  }

  // Refetch when the status filter changes (reset to first page).
  watch(statusFilter, () => {
    page.value = 1
    loadPosts()
  })
  watch([folderFilter, tagFilter], () => {
    page.value = 1
    loadPosts()
  })

  const filteredPosts = computed(() => {
    const q = searchQuery.value?.trim().toLowerCase()
    let list = posts.value
    if (q) {
      list = list.filter(
        (p) =>
          (p.title || "").toLowerCase().includes(q) ||
          (p.content || "").toLowerCase().includes(q) ||
          (p.tags || "").toLowerCase().includes(q)
      )
    }
    // Pin starred notes first — favorites stay at top, stable order otherwise
    return [...list].sort((a, b) => {
      if (a.isFavorite && !b.isFavorite) return -1
      if (!a.isFavorite && b.isFavorite) return 1
      return 0
    })
  })

  // Client-side favorite toggle (localStorage-backed, Design §10).
  function handleFavoriteToggled(id) {
    const p = posts.value.find((x) => x.id === id)
    if (p) p.isFavorite = !p.isFavorite
  }

  async function handleDuplicate(post) {
    try {
      await createPost({ title: post.title + " (copy)", content: post.content, status: post.status, folder: post.folder, tags: post.tags })
      toast("Note duplicated")
      page.value = 1
      loadPosts()
    } catch (err) {
      toast(err.response?.data?.message || err.message || "Failed to duplicate", { error: true })
    }
  }
  async function handleArchive(post) {
    try {
      await archivePost(post.id)
      posts.value = posts.value.filter((p) => p.id !== post.id)
      total.value = Math.max(0, total.value - 1)
      toast("Note archived")
    } catch (err) {
      toast(err.response?.data?.message || err.message || "Failed to archive", { error: true })
    }
  }
  async function handleUnarchive(post) {
    try {
      await unarchivePost(post.id)
      posts.value = posts.value.filter((p) => p.id !== post.id)
      total.value = Math.max(0, total.value - 1)
      toast("Note restored")
    } catch (err) {
      toast(err.response?.data?.message || err.message || "Failed to restore", { error: true })
    }
  }

  // Soft-delete → trash flow (Design §16).
  function handleDelete(post) {
    confirmPost.value = post
  }
  function cancelDelete() {
    confirmPost.value = null
  }
  async function confirmDelete() {
    if (!confirmPost.value) return
    const post = confirmPost.value
    confirmPost.value = null
    try {
      await deletePost(post.id)
      posts.value = posts.value.filter((p) => p.id !== post.id)
      total.value = Math.max(0, total.value - 1)
      toast("Post moved to trash")
      if (posts.value.length === 0 && page.value < totalPages.value) {
        page.value -= 1
      }
    } catch (err) {
      toast(err.response?.data?.message || err.message || "Failed to move to trash", {
        error: true,
      })
    }
  }

  // Create/edit form helpers.
  function openCreate() {
    editingPost.value = null
    showForm.value = true
  }
  function openEdit(post) {
    editingPost.value = post
    showForm.value = true
  }
  function closeForm() {
    showForm.value = false
    editingPost.value = null
  }

  // Submit handler mirrors backend DTO: { title, content, status? } (Design §22).
  async function submitForm(payload) {
    try {
      if (editingPost.value) {
        await updatePost(editingPost.value.id, payload)
        toast("Post updated successfully")
      } else {
        await createPost(payload)
        toast("Post created successfully")
      }
      closeForm()
      // Reset to first page so freshly created posts are visible, then refetch.
      page.value = 1
      loadPosts()
    } catch (err) {
      toast(
        err.response?.data?.message || err.message || "Failed to save post",
        { error: true }
      )
    }
  }

  onMounted(loadPosts)

  return {
    // data
    posts,
    loading,
    error,
    page,
    total,
    totalPages,
    searchQuery,
    statusFilter,
    folderFilter,
    tagFilter,
    view,
    filteredPosts,
    LIMIT,
    // delete flow
    confirmPost,
    handleDelete,
    cancelDelete,
    confirmDelete,
    // favorite
    handleFavoriteToggled,
    handleDuplicate,
    handleArchive,
    handleUnarchive,
    // form
    showForm,
    editingPost,
    openCreate,
    openEdit,
    closeForm,
    submitForm,
    // nav
    gotoPage: (p) => {
      const target = Math.max(1, Math.min(p, totalPages.value))
      if (target !== page.value) page.value = target
      loadPosts()
    },
    // external reload
    loadPosts,
  }
}
