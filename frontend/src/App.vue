<script lang="ts">
export const description = 'A simple login form.'
</script>

<script setup lang="ts">
import Toggle from './components/Toggle.vue';

</script>

<template>
  <div class="relative min-h-screen w-full">
    <div class="absolute top-4 right-4">
      <Toggle />
    </div>
    
    <div class="flex min-h-screen w-full items-center justify-center px-4">
       <router-view />
    </div>
  </div>
</template>