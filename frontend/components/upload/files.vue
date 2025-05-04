<script setup lang="ts">
import { toast } from 'vue-sonner';

const { formatBytes } = useUtils();
const files = useState<any[]>('files', () => []);
const isDragging = ref(false);
const isUploading = ref(false);
const uploadInput = useTemplateRef('uploadInput');
const fileTotalBytes = useState<number>('fileTotalBytes', () => 0);

function upload() {
  if (files.value.length === 0) {
    toast('No Files Selected', {
      description: 'Please select one or more files to proceed',
    })

    return
  } else if (isUploading.value) {
    toast('Upload in Progress', {
      description: 'Please wait until the current upload is complete',
    })

    return
  }

  isUploading.value = true;

  setTimeout(() => {
    files.value.length = 0;
    isUploading.value = false;
  }, 3000);
}

function filesAdd(flist: FileList | null | undefined) {
  if (!flist) {
    return;
  }

  for (const file of flist) {
    if (files.value.find((item) => item.name == file.name)) {
      continue;
    }

    files.value = [...files.value, file];
  }

  fileTotalBytes.value = 0;
  for (const file of files.value) {
    fileTotalBytes.value += file.size;
  }
}
function filesRm(index: number) {
  fileTotalBytes.value -= files.value[index].size;
  files.value.splice(index, 1);
}

function drop(event: DragEvent) {
  event.preventDefault();
  isDragging.value = false;
  filesAdd(event.dataTransfer?.files);
}
function dragover(event: DragEvent) {
  event.preventDefault();
  isDragging.value = true;
}
function dragleave() {
  isDragging.value = false;
}
function addInput(event: Event) {
  const el = event.target as HTMLInputElement;
  filesAdd(el.files);
}
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
      <div v-if="isUploading" class="h-56 border-2 rounded-lg">
      </div>
      <div v-else @click="uploadInput?.click()" @dragover="dragover" @dragleave="dragleave" @drop="drop" class="border-2 border-dashed h-56 rounded-lg sm:hover:bg-accent flex items-center cursor-pointer">
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
          {{ formatBytes(fileTotalBytes)}} in total
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
