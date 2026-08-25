import axios from "axios"

const api = axios.create({
  baseURL: "/api",
  headers: {
    "Content-Type": "application/json",
  },
})

/**
 * Fetch a paginated list of posts.
 * `status`/`folder`/`tags` are forwarded when provided (API supports ?status=&?folder=&?tags=).
 * @param {{ page?: number, limit?: number, status?: string, folder?: string, tags?: string }} [opts={}]
 */
export const fetchPosts = (opts = {}) => {
  const params = { page: opts.page, limit: opts.limit }
  if (opts.status) params.status = opts.status
  if (opts.folder) params.folder = opts.folder
  if (opts.tags) params.tags = opts.tags
  return api.get("/posts", { params })
}

/** GET /api/posts/trashed → paginated list of soft-deleted posts (Design §17). */
export const fetchTrashedPosts = (opts = {}) => {
  const params = { page: opts.page, limit: opts.limit }
  if (opts.status) params.status = opts.status
  if (opts.folder) params.folder = opts.folder
  if (opts.tags) params.tags = opts.tags
  return api.get("/posts/trashed", { params })
}

/** GET /api/posts/:id → 200 { data } */
export const getPost = (id) => api.get(`/posts/${id}`)

/**
 * POST /api/posts {title, content, status?}.
 * `status` is forwarded when present (create default = published).
 */
export const createPost = (data) => api.post("/posts", data)

/**
 * PUT /api/posts/:id {title, content, status?}.
 * `status` is forwarded when present so edits preserve Published/Draft.
 */
export const updatePost = (id, data) => api.put(`/posts/${id}`, data)

/** DELETE /api/posts/:id → soft-delete → 200 */
export const deletePost = (id) => api.delete(`/posts/${id}`)

/** PUT /api/posts/:id/restore → un-soft-delete (Design §17) */
export const restorePost = (id) => api.put(`/posts/${id}/restore`)

/** PUT /api/posts/:id/archive → status=archived (Nodex §35) */
export const archivePost = (id) => api.put(`/posts/${id}/archive`)

/** PUT /api/posts/:id/unarchive → status=published (Nodex §35) */
export const unarchivePost = (id) => api.put(`/posts/${id}/unarchive`)

/** DELETE /api/posts/:id/permanent → hard-delete a trashed post (Design §17) */
export const deletePermanent = (id) => api.delete(`/posts/${id}/permanent`)

export default api
