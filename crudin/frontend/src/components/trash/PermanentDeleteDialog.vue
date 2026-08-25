<!-- PermanentDeleteDialog.vue — Design §16/§17. Destructive confirmation for hard-delete.
     Modeled on DeletePostConfirmDialog, red destructive action, "This cannot be undone." -->
<template>
  <teleport to="body">
    <transition
      enter-active-class="transition ease-out duration-150"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition ease-in duration-150"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="show"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4 outline-none"
        role="dialog"
        aria-modal="true"
        :aria-labelledby="titleId"
        @click.self="onCancel"
        @keydown.esc.prevent.stop="onCancel"
        tabindex="-1"
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
            class="w-full max-w-md overflow-hidden rounded-[12px] border border-[#E4E4E7] bg-[#FFFFFF] p-6 shadow-xl"
          >
            <h2 :id="titleId" class="text-[20px] font-semibold text-[#18181B]">
              Delete permanently?
            </h2>
            <p class="mt-2 text-[#71717A]">
              Are you sure you want to permanently delete
              <span class="font-medium text-[#18181B]">
                "{{ post?.title }}"
              </span>
              ? This cannot be undone.
            </p>

            <div class="mt-6 flex justify-end gap-3">
              <button
                type="button"
                @click="onCancel"
                ref="cancelBtn"
                class="rounded-[8px] border border-[#E4E4E7] bg-[#FFFFFF] px-4 py-2 text-[15px] font-medium text-[#18181B] transition-colors hover:bg-[#F5F5F6] focus:outline-none focus:ring-2 focus:ring-[#4F46E5]/30 focus:ring-offset-2"
              >
                Cancel
              </button>
              <button
                type="button"
                @click="onConfirm"
                ref="confirmBtn"
                class="rounded-[8px] border border-transparent bg-red-600 px-4 py-2 text-[15px] font-medium text-[#FFFFFF] transition-colors hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500/30 focus:ring-offset-2"
              >
                Delete Permanently
              </button>
            </div>
          </div>
        </transition>
      </div>
    </transition>
  </teleport>
</template>

<script setup>
import { ref, nextTick, watch } from "vue"

const props = defineProps({
  post: { type: Object, default: null },
  show: { type: Boolean, default: false },
})
const emit = defineEmits(["cancel", "confirm"])

const titleId = "delete-permanent-title"
const cancelBtn = ref(null)

function onCancel() { emit("cancel") }
function onConfirm() { emit("confirm") }

// Move focus into the dialog on open for keyboard accessibility (§29).
watch(
  () => props.show,
  (on) => {
    if (on) nextTick(() => cancelBtn.value?.focus())
  }
)
</script>
