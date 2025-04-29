<script setup lang="ts">
import { formatBytes } from '~/lib/utils';
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Files</CardTitle>
      <CardDescription>
        Your awesome frenly neighbourhood file sharing website.
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

      <div v-if="haveAtleastTwoFiles" class="flex justify-between px-2 font-bold">
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
        <Button variant="ghost" class="my-auto" @click="filesRm(index)">
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
      haveAtleastTwoFiles: false,
      totalBytes: 0,
      files: new Array(),
    };
  },
  methods: {
    upload() {
      this.files = new Array();
      this.haveAtleastTwoFiles = false;
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

      if (this.files.length > 1) {
        this.haveAtleastTwoFiles = true;
      }
    },
    filesRm(index: number) {
      this.files.splice(index, 1);
    },
    drop(event: DragEvent) {
      event.preventDefault();
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
