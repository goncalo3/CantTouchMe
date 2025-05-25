<template>
  <div class="fixed inset-0 z-50 bg-background/80 flex items-center justify-center">
    <div class="bg-background rounded-lg p-6 w-full max-w-md shadow-xl border">
      <!-- Title slot -->
      <div v-if="$slots.title" class="text-lg font-semibold mb-2">
        <slot name="title" />
      </div>

      <div v-if="$slots.description" class="text-sm text-muted-foreground mb-4">
        <slot name="description" />
      </div>

      <div class="mb-4">
        <slot />
      </div>

      <div v-if="$slots.footer" class="flex justify-end gap-2 mt-4">
        <slot name="footer" />
      </div>
    </div>
  </div>
</template>

<script setup>
defineEmits(['close'])

document.addEventListener('keydown', (e) => {
  if (e.key === 'Escape') {
    emit('close')
  }
})
</script>
