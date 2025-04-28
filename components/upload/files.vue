<template>
  <Card>
    <CardHeader>
      <CardTitle>Files</CardTitle>
      <CardDescription>
        Your frenly neighbourhood file sharing web thing.
      </CardDescription>
    </CardHeader>
    <CardContent class="space-y-2">
      <div @click="clickUploadInput" @dragover="dragover" @dragleave="dragleave" @drop="drop" class="border-2 border-dashed p-16 rounded-lg sm:hover:bg-accent">
        <div class="w-min mx-auto">
          <Icon v-if="isDragging" name="mdi:add" class="text-7xl text-muted-foreground"/>
          <Icon v-else name="mdi:cloud-upload-outline" class="text-7xl text-muted-foreground"/>
        </div>
        <p v-if="isDragging" class="text-sm text-muted-foreground text-center">Drop & and I'll catch</p>
        <p v-else class="text-sm text-muted-foreground text-center">Drag & drop files here, or click to select files</p>
      </div>
      <input ref="uploadInput" type="file" multiple="true" class="invisible" @change="addInput" >

      <div v-for="file in files"q class="border rounded-lg p-2 flex justify-between gap-2">
        <div class="flex gap-2 truncate">
          <Icon name="uil:file"  class="text-4xl my-auto"/>
          <div>
            <div class="font-bold my-auto truncate">
              {{ file.name }}
            </div>
            <div class="text-muted-foreground text-sm">
              {{ file.size }}
            </div>
          </div>
        </div>
        <Button variant="ghost" class="my-auto">
          <Icon name="mdi:close" />
        </Button>
      </div>
    </CardContent>
    <CardFooter class="flex justify-end">
      <Button @click="upload">Upload</Button>
    </CardFooter>
  </Card>
</template>

<script lang="ts">
export default {
  data() {
    return {
      isDragging: false,
      files: new Array(),
    };
  },
  methods: {
    upload() {
        console.log(this.files)
      this.files = new Array();
    },
    drop(event: DragEvent) {
      event.preventDefault();

      if (!event.dataTransfer) {
        return
      }
      this.files = [...this.files, ...event.dataTransfer.files];
    },
    dragover(event: DragEvent) {
      event.preventDefault();
      this.isDragging = true;
    },
    dragleave() {
      this.isDragging = false;
    },
    addInput(event: Event) {
      const el = event.target as HTMLInputElement;

      if (!el.files) {
        return
      }
      this.files = [...this.files, ...el.files];
    },
    clickUploadInput() {
      const el = this.$refs.uploadInput as HTMLInputElement;
      el.click();
    }
  }
}
</script>
