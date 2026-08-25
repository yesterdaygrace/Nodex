<!-- PostCardSkeleton.vue — Design §20. Matches card structure while loading. -->
<template>
  <article
    class="relative flex h-full flex-col gap-3 rounded-[16px] bg-[#FFFFFF] p-6 pb-4 ring-1 ring-[#E4E4E7] shadow-sm"
  >
    <div class="absolute top-4 right-4 h-4 w-4 rounded-full bg-[#E4E4E7] animate-pulse"></div>
    <div class="h-5 w-3/4 rounded bg-[#E4E4E7] animate-pulse"></div>
    <div class="mt-1 space-y-1.5">
      <div class="h-4 w-full rounded bg-[#E4E4E7] animate-pulse"></div>
      <div class="h-4 w-5/6 rounded bg-[#E4E4E7] animate-pulse"></div>
      <div class="h-4 w-2/3 rounded bg-[#E4E4E7] animate-pulse"></div>
    </div>
    <div class="h-5 w-16 rounded-[6px] bg-[#E4E4E7] animate-pulse"></div>
    <footer class="mt-auto pt-3 ring-t-1 ring-[#E4E4E7] flex items-center justify-between">
      <div class="flex items-center gap-1.5">
        <div class="h-3 w-3 rounded bg-[#E4E4E7] animate-pulse"></div>
        <div class="h-3 w-16 rounded bg-[#E4E4E7] animate-pulse"></div>
      </div>
      <div class="flex items-center gap-1">
        <div class="h-7 w-7 rounded-[8px] bg-[#E4E4E7] animate-pulse"></div>
        <div class="h-7 w-7 rounded-[8px] bg-[#E4E4E7] animate-pulse"></div>
      </div>
    </footer>
  </article>
</template>

<script setup>
// purely presentational; no script needed.
</script>
