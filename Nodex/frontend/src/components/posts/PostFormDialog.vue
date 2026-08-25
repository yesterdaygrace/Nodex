<!-- PostFormDialog.vue — Design §6.2 / §26. Modal create/edit form for posts.
     Teleport to body, overlay, focus trap, Esc + click-outside cancel.
     Raw Tailwind only, purple #4F46E5 accent, 12px dialog radius. -->
<template>
  <teleport to="body">
    <transition
      appear
      enter-active-class="transition ease-out duration-150"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition ease-in duration-150"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="show"
        ref="dialog"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4 outline-none"
        role="dialog"
        aria-modal="true"
        :aria-labelledby="titleId"
        @click.self="onCancel"
        @keydown.tab="trapFocus"
>
        <transition
          enter-active-class="transition ease-out duration-150"
          enter-from-class="opacity-0 scale-95"
          enter-to-class="opacity-100 scale-100"
          leave-active-class="transition ease-in duration-150"
          leave-from-class="opacity-100 scale-100"
          leave-to-class="opacity-0 scale-95"
        >
          <div
            class="w-full max-w-lg overflow-hidden rounded-[12px] border border-[#E4E4E7] bg-[#FFFFFF] p-6 shadow-xl"
          >
            <h2 :id="titleId" class="text-[20px] font-semibold text-[#18181B]">
              {{ isEdit ? "Edit Post" : "New Post" }}
            </h2>

            <form @submit.prevent="onSubmit" class="mt-5 flex flex-col gap-4">
              <!-- Title -->
              <div class="flex flex-col gap-1.5">
                <label for="post-title" class="text-[13px] font-medium text-[#18181B]">Title</label>
                <input
                  id="post-title"
                  ref="titleInput"
                  v-model="title"
                  type="text"
                  placeholder="Write a title..."
                  autocomplete="off"
                  class="w-full rounded-[8px] border border-[#E4E4E7] bg-[#FFFFFF] px-3 py-2 text-[15px] text-[#18181B] outline-none transition-colors focus:border-[#4F46E5] focus:ring-2 focus:ring-[#4F46E5]/30"
                />
                <p v-if="touched && !title.trim()" class="text-[12px] text-[#dc2626]">Title is required</p>
              </div>

              <!-- Content -->
              <div class="flex flex-col gap-1.5">
                <label for="post-content" class="text-[13px] font-medium text-[#18181B]">Content</label>
                <textarea
                  id="post-content"
                  v-model="content"
                  placeholder="Write the post content..."
                  rows="5"
                  class="w-full rounded-[8px] border border-[#E4E4E7] bg-[#FFFFFF] px-3 py-2 text-[15px] text-[#18181B] outline-none transition-colors resize-y focus:border-[#4F46E5] focus:ring-2 focus:ring-[#4F46E5]/30"
                ></textarea>
                <p v-if="touched && !content.trim()" class="text-[12px] text-[#dc2626]">Content is required</p>
              </div>

              <!-- Status select (Design §25) -->
              <div class="flex flex-col gap-1.5">
                <label class="text-[13px] font-medium text-[#18181B]">Status</label>
                <div class="relative" ref="statusRef">
                  <button
                    type="button"
                    @click="statusOpen = !statusOpen"
                    :aria-expanded="statusOpen"
                    class="inline-flex w-full items-center justify-between gap-2 rounded-[8px] border border-[#E4E4E7] bg-[#FFFFFF] px-3 py-2 text-[15px] text-[#18181B] outline-none transition-colors hover:bg-[#F5F5F6] focus:ring-2 focus:ring-[#4F46E5]/30"
                  >
                    <span>{{ statusLabel }}</span>
                    <Icon name="chevron-down" size="16" class="text-[#71717A] icon-pressable" />
                  </button>
                  <transition
                    enter-active-class="transition ease-out duration-150"
                    enter-from-class="opacity-0 scale-95"
                    enter-to-class="opacity-100 scale-100"
                    leave-active-class="transition ease-in duration-150"
                    leave-from-class="opacity-100 scale-95"
                    leave-to-class="opacity-0 scale-95"
                  >
                    <div
                      v-show="statusOpen"
                      class="absolute z-10 mt-1 w-full overflow-hidden rounded-[8px] border border-[#E4E4E7] bg-[#FFFFFF] shadow-lg"
                    >
                      <ul role="listbox" aria-label="Post status" class="py-1 text-[15px]">
                        <li
                          v-for="opt in statusOptions"
                          :key="opt.value"
                          :aria-selected="status === opt.value"
                          @click="setStatus(opt.value)"
                          class="cursor-pointer px-3 py-1.5 hover:bg-[#F5F3FF] hover:text-[#4F46E5] first:rounded-t-[8px] last:rounded-b-[8px]"
                          :class="status === opt.value ? 'bg-[#EFF0FF] text-[#4F46E5]' : 'text-[#18181B]'"
                        >
                          {{ opt.label }}
                        </li>
                      </ul>
                    </div>
                  </transition>
                </div>
              </div>

              <p v-if="touched && !statusOpen" class="text-[12px] text-[#71717A]">
                Saving as <span class="font-medium">{{ statusLabel }}</span>.
              </p>
            </form>

            <!-- Actions -->
            <div class="mt-6 flex justify-end gap-3">
              <button
                type="button"
                @click="onCancel"
                class="rounded-[8px] border border-[#E4E4E7] bg-[#FFFFFF] px-4 py-2 text-[15px] font-medium text-[#18181B] transition-colors hover:bg-[#F5F5F6] focus:outline-none focus:ring-2 focus:ring-[#4F46E5]/30 focus:ring-offset-2"
              >
                Cancel
              </button>
              <button
                type="button"
                @click="onSubmit"
                class="rounded-[8px] border border-transparent bg-[#4F46E5] px-4 py-2 text-[15px] font-medium text-[#FFFFFF] transition-colors hover:bg-[#4338CA] disabled:cursor-not-allowed disabled:opacity-60 focus:outline-none focus:ring-2 focus:ring-[#4F46E5] focus:ring-offset-2"
              >
                {{ isEdit ? "Save" : "Create" }}
              </button>
            </div>
          </div>
        </transition>
      </div>
    </transition>
  </teleport>
</template>

<script setup>
import { ref, computed, watch, nextTick, onMounted, onBeforeUnmount } from "vue"
import Icon from "../ui/Icon.vue"

const props = defineProps({
  show: { type: Boolean, default: false },
  post: { type: Object, default: null },
})
const emit = defineEmits(["submit", "cancel"])

const title = ref("")
const content = ref("")
const status = ref("published")
const touched = ref(false)
const statusOpen = ref(false)
const titleInput = ref(null)
const dialog = ref(null)
const statusRef = ref(null)

const titleId = "post-form-title"

const isEdit = computed(() => !!props.post)

const statusOptions = [
  { label: "Published", value: "published" },
  { label: "Draft", value: "draft" },
]

const statusLabel = computed(() => {
  const found = statusOptions.find((o) => o.value === status.value)
  return found ? found.label : "Published"
})

watch(
  () => props.show,
  (visible) => {
    if (visible) {
      resetForm()
      nextTick(() => titleInput.value?.focus())
    }
  }
)

function resetForm() {
  if (props.post) {
    title.value = props.post.title || ""
    content.value = props.post.content || ""
    status.value = props.post.status || "published"
  } else {
    title.value = ""
    content.value = ""
    status.value = "published"
  }
  touched.value = false
  statusOpen.value = false
}

function setStatus(val) {
  status.value = val
  statusOpen.value = false
}

function onCancel() {
  emit("cancel")
}

function onSubmit() {
  touched.value = true
  if (!title.value.trim() || !content.value.trim()) return
  emit("submit", {
    title: title.value.trim(),
    content: content.value.trim(),
    status: status.value,
  })
}

// Lightweight focus trap: keep Tab within the dialog.
function focusableNodes() {
  if (!dialog.value) return []
  return dialog.value.querySelectorAll(
    "button, input, select, textarea, [tabindex]:not([tabindex='-1'])"
  )
}
function trapFocus(e) {
  const nodes = Array.from(focusableNodes()).filter(
    (n) => !n.hasAttribute("disabled") && n.offsetParent !== null
  )
  if (nodes.length === 0) return
  const first = nodes[0]
  const last = nodes[nodes.length - 1]
  if (e.shiftKey && document.activeElement === first) {
    e.preventDefault()
    last.focus()
  } else if (!e.shiftKey && document.activeElement === last) {
    e.preventDefault()
    first.focus()
  }
}

function onOutsideStatus(e) {
  if (statusOpen.value && statusRef.value && !statusRef.value.contains(e.target)) {
    statusOpen.value = false
  }
}
function onKeyDown(e) {
  if (e.key === "Escape") {
    // Close the status dropdown first; otherwise cancel the dialog.
    if (statusOpen.value) {
      statusOpen.value = false
    } else {
      onCancel()
    }
  }
}
onMounted(() => {
  document.addEventListener("click", onOutsideStatus)
  document.addEventListener("keydown", onKeyDown)
})
onBeforeUnmount(() => {
  document.removeEventListener("click", onOutsideStatus)
  document.removeEventListener("keydown", onKeyDown)
})
</script>
