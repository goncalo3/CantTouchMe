<template>
  <div class="fixed inset-0 flex items-center justify-center z-50">
    <div class="fixed inset-0 bg-black/50" @click="$emit('cancel')" />
    <Card class="relative w-96 bg-background border shadow-lg">
      <CardHeader>
        <CardTitle>Confirm Action</CardTitle>
        <CardDescription>{{ message }}</CardDescription>
      </CardHeader>
      <CardFooter class="flex justify-end gap-2">
        <Button variant="outline" @click="$emit('cancel')">Cancel</Button>
        <Button variant="destructive" @click="$emit('confirm')">Confirm</Button>
      </CardFooter>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { Card, CardHeader, CardTitle, CardDescription, CardFooter } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

defineProps({
  message: {
    type: String,
    required: true,
  },
});
</script>
