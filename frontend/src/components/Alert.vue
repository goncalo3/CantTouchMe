<template>
  <div v-if="alert" class="fixed inset-0 flex items-center justify-center z-50">
    <div class="fixed inset-0 bg-black/50" @click="$emit('close')" />
    <Card class="relative w-96 bg-background border shadow-lg">
      <CardHeader>
        <div class="flex items-center gap-2">
          <AlertCircle v-if="alert.type === 'error'" class="h-5 w-5 text-destructive" />
          <Info v-if="alert.type === 'info'" class="h-5 w-5 text-primary" />
          <CardTitle :class="{
            'text-destructive': alert.type === 'error',
            'text-primary': alert.type === 'info',
          }">
            {{ alert.type === 'error' ? 'Error' : 'Success' }}
          </CardTitle>
        </div>
      </CardHeader>
      <CardContent>
        <p class="text-foreground text-center">
          {{ formatMessage(alert.message) }}
        </p>
      </CardContent>
      <CardFooter class="flex justify-end">
        <Button variant="ghost" @click="$emit('close')">Close</Button>
      </CardFooter>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { Card, CardHeader, CardTitle, CardContent, CardFooter } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { AlertCircle, Info } from 'lucide-vue-next';
import type { Alert as AlertType } from '@/store/notifications';

defineProps({
  alert: {
    type: Object as () => AlertType | null,
    required: true,
  },
});

function formatMessage(message: string): string {
  // Remove common prefixes that make messages look unprofessional
  message = message.replace(/^(Login failed:|Registration failed:|Failed to delete account:)/, '');
  
  // Capitalize first letter
  message = message.charAt(0).toUpperCase() + message.slice(1);
  
  // Add period if missing
  if (!message.endsWith('.')) {
    message += '.';
  }
  
  return message;
}
</script>
