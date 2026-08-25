<!-- Icon.vue — thin Vue 3 wrapper around HugeIcons (hugeicons-vue free).
     Stable public API: name, size, filled, className.
     Renders stroke-rounded SVGs at 18–20 px with Emil interaction polish. -->
<template>
  <component
    v-if="resolvedComponent"
    :is="resolvedComponent"
    :size="size"
    :class="resolvedClasses"
    :fill="filled ? 'currentColor' : 'none'"
    v-bind="$attrs"
  />
</template>

<script setup>
import { computed } from "vue"
import {
  Menu01Icon,
  Search01Icon,
  Notification01Icon,
  StarIcon,
  PencilEdit02Icon,
  Delete02Icon,
  Clock01Icon,
  PlusSignIcon,
  GridViewIcon,
  ListViewIcon,
  ArrowDown01Icon,
  File01Icon,
  File02Icon,
  Folder01Icon,
  Tag01Icon,
  Settings01Icon,
  Tick01Icon,
  Cancel01Icon,
  Archive01Icon,
  Copy01Icon,
} from "hugeicons-vue"

const props = defineProps({
  name: { type: String, required: true },
  size: { type: Number, default: 19 },
  filled: { type: Boolean, default: false },
  className: { type: String, default: "" },
})

/* ── Name → HugeIcons component registry (stroke-rounded) ────────────── */
const ICON_MAP = {
  menu: Menu01Icon,
  search: Search01Icon,
  bell: Notification01Icon,
  star: StarIcon,
  pencil: PencilEdit02Icon,
  "trash-2": Delete02Icon,
  clock: Clock01Icon,
  plus: PlusSignIcon,
  "layout-grid": GridViewIcon,
  list: ListViewIcon,
  "chevron-down": ArrowDown01Icon,
  file: File01Icon,
  "file-text": File02Icon,
  folder: Folder01Icon,
  tag: Tag01Icon,
  settings: Settings01Icon,
  check: Tick01Icon,
  x: Cancel01Icon,
  archive: Archive01Icon,
  copy: Copy01Icon,
}

const resolvedComponent = computed(() => ICON_MAP[props.name])

/* Merge consumer className + filled marker for CSS override */
const resolvedClasses = computed(() => {
  const cls = [props.className]
  if (props.filled) cls.push("icon-filled")
  return cls.filter(Boolean).join(" ")
})
</script>

<style>
/* ══════════════════════════════════════════════════════════════════════
   Emil interaction polish — composable utility classes
   GPU-friendly: only transform + opacity are animated.
   ══════════════════════════════════════════════════════════════════════ */

/* Base: smooth ease-out transition on every HugeIcons SVG */
.hugeicons {
  transition: transform 0.18s cubic-bezier(0.23, 1, 0.32, 1),
              opacity 0.18s cubic-bezier(0.23, 1, 0.32, 1);
}

/* ── Press: 3 % scale on active ─────────────────────────────────────── */
.icon-pressable {
  cursor: pointer;
}
.icon-pressable:active {
  transform: scale(0.97);
}

/* ── Hover: subtle opacity fade (only on devices that support hover) ─ */
@media (hover: hover) {
  .icon-pressable:hover {
    opacity: 0.85;
  }
}

/* ── Focus ring (keyboard accessibility) ───────────────────────────── */
.icon-pressable:focus-visible {
  outline: 2px solid currentColor;
  outline-offset: 2px;
}

/* ── Filled: strip stroke so fill-only shapes stay crisp ───────────── */
.icon-filled path {
  stroke: transparent;
}

/* ── Reduced motion: kill all transitions & transforms ─────────────── */
@media (prefers-reduced-motion: reduce) {
  .hugeicons {
    transition: none;
  }
  .icon-pressable:active {
    transform: none;
  }
  .icon-grid [data-icon] {
    opacity: 1;
    animation: none;
  }
}

/* ── Grid stagger reveal: 40 ms per item (30–50 ms range) ─────────── */
.icon-grid [data-icon] {
  opacity: 0;
  animation: icon-stagger 0.35s cubic-bezier(0.23, 1, 0.32, 1) both;
}
.icon-grid [data-icon]:nth-child(1)  { animation-delay:   0ms; }
.icon-grid [data-icon]:nth-child(2)  { animation-delay:  40ms; }
.icon-grid [data-icon]:nth-child(3)  { animation-delay:  80ms; }
.icon-grid [data-icon]:nth-child(4)  { animation-delay: 120ms; }
.icon-grid [data-icon]:nth-child(5)  { animation-delay: 160ms; }
.icon-grid [data-icon]:nth-child(6)  { animation-delay: 200ms; }
.icon-grid [data-icon]:nth-child(7)  { animation-delay: 240ms; }
.icon-grid [data-icon]:nth-child(8)  { animation-delay: 280ms; }
.icon-grid [data-icon]:nth-child(9)  { animation-delay: 320ms; }
.icon-grid [data-icon]:nth-child(n+10) { animation-delay: 360ms; }

@keyframes icon-stagger {
  to { opacity: 1; }
}
</style>
