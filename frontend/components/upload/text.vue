<script setup lang="ts">
import Textarea from '~/components/ui/textarea/Textarea.vue';
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Text</CardTitle>
      <CardDescription>
        Your comfy pastebin, type your message here.
      </CardDescription>
    </CardHeader>

    <CardContent class="space-y-2">
      <Textarea class="h-64"/>
    </CardContent>
    <CardFooter class="flex justify-end">
      <Button>Save</Button>
    </CardFooter>
  </Card>
</template>
