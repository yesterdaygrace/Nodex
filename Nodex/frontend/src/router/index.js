// src/router/index.js — Phase-2 routing. Raw vue-router@4 (no shadcn).
// Routes map to the Phase-2 views. Categories/Tags/Settings remain inert nav.
import { createRouter, createWebHistory } from "vue-router"
import PostsView from "../views/PostsView.vue"
import DraftsView from "../views/DraftsView.vue"
import TrashView from "../views/TrashView.vue"
import FavoritesView from "../views/FavoritesView.vue"
import ArchiveView from "../views/ArchiveView.vue"
import NoteDetail from "../views/NoteDetail.vue"

export const routes = [
  { path: "/", component: PostsView, name: "all-notes" },
  { path: "/favorites", component: FavoritesView, name: "favorites" },
  { path: "/archive", component: ArchiveView, name: "archive" },
  { path: "/trash", component: TrashView, name: "trash" },
  { path: "/notes/:id", component: NoteDetail, name: "note-detail" },
  { path: "/drafts", redirect: "/archive", name: "drafts-redirect" },
]

export default createRouter({
  history: createWebHistory(),
  routes,
})
