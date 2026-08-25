<!-- NoteDetail.vue — Nodex §31 detail + §32 editor shell. Simple detail view. -->
<template>
  <div class="flex flex-col gap-6">
    <button type="button" @click="$router.back()" class="self-start text-[14px] text-[#64748B] hover:text-[#111827]">← Back to Notes</button>
    <div v-if="loading" class="text-[#64748B]">Loading...</div>
    <div v-else-if="error" class="text-red-600">{{ error }} <button @click="load" class="underline">Retry</button></div>
    <div v-else-if="note" class="rounded-[12px] bg-white p-6 ring-1 ring-[#E2E8F0]">
      <h1 class="text-[24px] font-bold text-[#111827]">{{ note.title }}</h1>
      <div v-if="note.tags" class="mt-2 flex flex-wrap gap-1.5">
        <span v-for="t in note.tags.split(',').map(s=>s.trim()).filter(Boolean)" :key="t" class="rounded-full bg-[#EFF6FF] px-2 py-0.5 text-[12px] text-[#1D4ED8]">#{{t}}</span>
      </div>
      <div class="mt-4 text-[15px] leading-relaxed text-[#334155] whitespace-pre-wrap">{{ note.content }}</div>
      <div class="mt-6 flex items-center justify-between border-t border-[#F1F5F9] pt-4">
        <span class="text-[12px] text-[#94A3B8]">{{ note.updated_at || note.created_at }}</span>
        <div class="flex gap-2">
          <button @click="edit" class="rounded-[8px] bg-[#6D28D9] px-3 py-1.5 text-white text-[14px]">Edit</button>
          <button @click="archive" class="rounded-[8px] border px-3 py-1.5 text-[14px]">Archive</button>
          <button @click="trash" class="rounded-[8px] bg-red-50 text-red-600 px-3 py-1.5 text-[14px]">Move to Trash</button>
        </div>
      </div>
    </div>
    <PostFormDialog v-if="showForm" :show="showForm" :post="note" @submit="save" @cancel="showForm=false" />
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue"
import { useRoute, useRouter } from "vue-router"
import { getPost, updatePost, deletePost, archivePost } from "../services/api.js"
import { useToast } from "../composables/useToast.js"
import PostFormDialog from "../components/posts/PostFormDialog.vue"

const route = useRoute()
const router = useRouter()
const toast = useToast()
const note = ref(null)
const loading = ref(true)
const error = ref(null)
const showForm = ref(false)

async function load(){
  loading.value=true; error.value=null
  try{ const r=await getPost(route.params.id); note.value=r.data.data }catch(e){ error.value=e.response?.data?.message||e.message }
  finally{ loading.value=false }
}
function edit(){ showForm.value=true }
async function save(payload){
  try{ await updatePost(note.value.id, payload); toast("Note updated"); showForm.value=false; load() }catch(e){ toast(e.message,{error:true}) }
}
async function archive(){ try{ await archivePost(note.value.id); toast("Note archived"); router.push("/archive") }catch(e){ toast(e.message,{error:true}) } }
async function trash(){ try{ await deletePost(note.value.id); toast("Moved to trash"); router.push("/trash") }catch(e){ toast(e.message,{error:true}) } }
onMounted(load)
</script>
