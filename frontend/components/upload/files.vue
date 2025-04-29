<script setup lang="ts">
import { formatBytes } from '~/lib/utils';
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Files</CardTitle>
      <CardDescription>
        Your frenly neighbourhood file sharing website.
      </CardDescription>
    </CardHeader>
    <CardContent class="space-y-2">
      <input ref="uploadInput" type="file" multiple="true" class="hidden" @change="addInput" />
      <div v-if="isUploading" class="h-56 border-2">
      </div>
      <div v-else @click="clickUploadInput" @dragover="dragover" @dragleave="dragleave" @drop="drop" class="border-2 border-dashed h-56 rounded-lg sm:hover:bg-accent flex items-center cursor-pointer">
        <div class="mx-auto">
          <div class="w-min mx-auto">
            <Icon v-if="isDragging" name="mdi:add" class="text-7xl text-muted-foreground"/>
            <Icon v-else name="mdi:cloud-upload-outline" class="text-7xl text-muted-foreground"/>
          </div>
          <p v-if="isDragging" class="text-sm text-muted-foreground text-center">Drop & and I'll catch</p>
          <p v-else class="text-sm text-muted-foreground text-center">Drag & drop files here, or click to select files</p>
        </div>
      </div>

      <div v-if="files.length > 0" class="h-4" />

      <div v-if="files.length > 1" class="flex justify-between px-2">
        <div>
          {{ files.length }} files selected
        </div>
        <div>
          {{ formatBytes(totalBytes)}} in total
        </div>
      </div>
      <div v-for="(file, index) in files"q class="border rounded-lg p-2 flex justify-between gap-2">
        <div class="flex gap-2 truncate">
          <Icon name="uil:file"  class="text-4xl my-auto"/>
          <div class="truncate">
            <div class="font-bold my-auto truncate">
              {{ file.name }}
            </div>
            <div class="text-muted-foreground text-sm">
              {{ formatBytes(file.size) }}
            </div>
          </div>
        </div>
        <Button v-if="!isUploading" variant="ghost" class="my-auto" @click="filesRm(index)">
          <Icon name="mdi:close" />
        </Button>
        <Icon v-else class="my-auto px-6" name="svg-spinners:dot-revolve" />
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
      isUploading: false,
      totalBytes: 0,
      files: new Array(),
    };
  },
  methods: {
    upload() {
      this.isUploading = true;

      setTimeout(() => {
        this.files = new Array();
        this.isUploading = false;
      }, 3000);
    },
    filesAdd(files: FileList | null | undefined) {
      if (!files) {
        return;
      }

      for (const file of files) {
        if (this.files.find((item) => item.name == file.name)) {
          continue;
        }

        this.files = [...this.files, file];
      }

      this.totalBytes = 0;
      for (const file of this.files) {
        this.totalBytes += file.size;
      }
    },
    filesRm(index: number) {
      this.totalBytes -= this.files[index].size;
      this.files.splice(index, 1);
    },
    drop(event: DragEvent) {
      event.preventDefault();
      this.isDragging = false;
      this.filesAdd(event.dataTransfer?.files);
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
      this.filesAdd(el.files);
    },
    clickUploadInput() {
      const el = this.$refs.uploadInput as HTMLInputElement;
      el.click();
    }
  }
}
</script>
