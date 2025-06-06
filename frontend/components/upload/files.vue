<script setup lang="ts">
import { toast } from 'vue-sonner';

const { formatData, filesZip, chunksFromBlob: chunksFromFile, formatDuration } = useUtils();
const files = useState<File[]>('files', () => []);
const isDragging = ref(false);
const filesIsUploading = useState('filesIsUploading', () => false);
const textIsUploading = useState('textIsUploading', () => false);
const isZipping = useState('filesIsZipping', () => false);
const isPaused = useState('filesIsPaused', () => false);
const fileUploadETA = useState<string>('fileUploadETA', () => "Infinity");
const bytesUploadedPerSecond = useState<number>('bytesUploadedPerSecond', () => 0);
const uploadLink = useState<string>('uploadLink');
const fileUploadDialog = useState<boolean>('fileUploadDialog');
const uploadInput = useTemplateRef('fileUploadInput');
const uploadedChunkCount = useState<number>('uploadedChunkCount', () => 0);
const fileUploadProgress = useState<number>('fileUploadProgress', () => 0);
const fileTotalBytes = useState<number>('fileTotalBytes', () => 0);
const serverConfig = useServerConfig();
const fileXhrReq = useState<XMLHttpRequest>('fileXhrReq');
const fileSpeedInterval = useState<number>('fileSpeedInterval');
const fileChunkToken = useState<string>('fileChunkToken');
const serverUrl = useState('serverUrl', () => window.location.origin)


type ChunkPostResp = {
  chunk_token: string
}

type ChunkPostReq = {
  file_size: number
  name: string
}

type ChunkCancelReq = {
  chunk_token: string
}

async function upload() {
  if (files.value.length === 0) {
    toast('No Files Selected', {
      description: 'Please select one or more files to proceed',
    });

    return
  } else if ((filesIsUploading.value && !isPaused.value) || textIsUploading.value) {
    toast('Upload in Progress', {
      description: 'Please wait until the current upload is complete',
    });

    return
  } else if (fileTotalBytes.value > serverConfig.value.file_size_limit) {
    toast('Files Too Large', {
      description: `These files exceeds the upload size limit of ${formatData(serverConfig.value.file_size_limit)}`,
    });

    return
  }
  filesIsUploading.value = true;
  isPaused.value = false;

  let file: Blob | File
  let body: ChunkPostReq
  if (files.value.length > 1) {
    isZipping.value = true;
    const zip = await filesZip(files.value as any);
    if (!isZipping.value) return;
    isZipping.value = false;

    if (zip.size > serverConfig.value.file_size_limit) {
        toast('Zip file Too Large', {
          description: `These files after zipping exceeds the upload size limit of ${formatData(serverConfig.value.file_size_limit)}`,
        });

        filesIsUploading.value = true;
        return
    }

    file = zip;
    body = {
      name: "zip.zip",
      file_size: file.size
    }
  } else {
    file = files.value[0]
    body = {
      name: files.value[0].name,
      file_size: file.size
    }
  }

  if (uploadedChunkCount.value == 0) {
    const req = await useFetch(`${serverUrl.value}/_alina/upload/chunked`, {
      method: "POST",
      body: body
    })
    if (req.status.value == "error") {
      toast(req.error.value?.name ?? "Error", {
        description: req.error.value?.message,
      });

      filesIsUploading.value = false;
      return
    }
    const { chunk_token } = req.data.value as ChunkPostResp
    fileChunkToken.value = chunk_token
  }

  let totalUploaded = 0;
  let prevTotalUploaded = 0;
  fileSpeedInterval.value = setInterval(() => {
    bytesUploadedPerSecond.value = totalUploaded - prevTotalUploaded;
    prevTotalUploaded = totalUploaded;
    fileUploadETA.value = formatDuration(bytesUploadedPerSecond.value, body.file_size - totalUploaded);
  }, 1000) as any;

  let responseText: string | undefined
  const chunks = chunksFromFile(file, serverConfig.value.chunk_size)
  for (let i = uploadedChunkCount.value; i < chunks.length; i++) {
    const data = new FormData();
    data.append("chunk", chunks[i])
    data.append("chunk_token", fileChunkToken.value)
    data.append("chunk_index", `${i+1}`)

    fileXhrReq.value = new XMLHttpRequest();
    fileXhrReq.value.open('PATCH', `${serverUrl.value}/_alina/upload/chunked`)
    let isError = false;
    for (let retry = 0, retries = 3; retry < retries; retry++) {
      await new Promise<string>((resolve, reject) => {
        fileXhrReq.value.onload = () => {
          resolve(fileXhrReq.value.responseText)
        }
        fileXhrReq.value.onerror = () => {
          reject();
        }
        fileXhrReq.value.upload.onprogress = (event) => {
          const eventUploaded = event.loaded > chunks[i].size ? chunks[i].size : event.loaded
          totalUploaded = uploadedChunkCount.value * serverConfig.value.chunk_size + eventUploaded
          fileUploadProgress.value = ((totalUploaded/file.size) * 100)
        }

        fileXhrReq.value.send(data)
      }).then((data) => {
        isError = false;
        uploadedChunkCount.value += 1;
        responseText = data
        retry = retries;
      }).catch(() => {
        isError = true;
        if (isPaused.value) {
          retry = retries;
        } else {
          retry += 1;
        }
      })
    }

    if (isError || (i+1 == chunks.length && !responseText)) {
      toast("Upload Failed", {
        description: "Please check your internet connection and try again",
      });

      pause();
      return
    }
  }

  clearInterval(fileSpeedInterval.value)
  uploadedChunkCount.value = 0
  uploadLink.value = `${serverUrl.value}/${responseText}`;
  fileUploadDialog.value = true;
  files.value.length = 0;
  filesIsUploading.value = false;
}

function pause() {
  isPaused.value = true;
  fileXhrReq.value.abort();
  clearInterval(fileSpeedInterval.value);
}
async function cancel() {
  fileXhrReq.value.abort();
  clearInterval(fileSpeedInterval.value)
  filesIsUploading.value = false;
  isZipping.value = false;
  uploadedChunkCount.value = 0;

  const body: ChunkCancelReq = {
    chunk_token: fileChunkToken.value
  }

  const req = await useFetch(`${serverUrl.value}/_alina/upload/chunked`, {
    method: "DELETE",
    body: body
  })
  if (req.status.value == "error") {
    toast("Cancel Failed", {
      description: req.error.value?.message,
    });
  }
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
      <input ref="fileUploadInput" type="file" multiple="true" class="hidden" @change="addInput" />

      <div v-if="filesIsUploading" class="h-56 border-2 rounded-lg p-6 space-y-4 flex flex-col justify-between">
        <div class="space-y-1.5 m-auto">
          <div v-if="!isZipping" class="text-4xl font-bold">
            {{ formatData(bytesUploadedPerSecond * 8, true)}}ps
          </div>
          <div v-else class="text-4xl font-bold">
            {{ formatData(fileTotalBytes)}}
          </div>

          <div v-if="isPaused" class="flex space-x-1 mx-auto w-min">
            <Icon name="lucide:circle-pause" class="my-auto" />
            <p class="my-auto">Paused</p>
          </div>
          <div v-else-if="isZipping" class="flex space-x-1 mx-auto w-min">
            <Icon name="svg-spinners:blocks-scale" class="my-auto" />
            <p class="my-auto">Zipping</p>
          </div>
          <div v-else class="flex space-x-1 mx-auto w-min">
            <Icon name="line-md:uploading-loop" class="my-auto" />
            <p class="my-auto">Uploading</p>
          </div>
        </div>

        <div v-if="!isZipping">
          <div class="flex justify-between">
            <div class="text-muted-foreground text-sm">
              {{fileUploadETA}}
            </div>
            <div class="flex">
              <div class="text-muted-foreground text-sm font-mono my-auto">
                {{fileUploadProgress.toFixed(2)}}
              </div>
              <div class="text-muted-foreground text-sm my-auto">%</div>
            </div>
          </div>
          <Progress :model-value="fileUploadProgress" />
        </div>
        <div v-else>
          <div class="flex justify-between">
            <div class="text-muted-foreground text-sm">
              Fast as possible
            </div>
            <div class="text-muted-foreground text-sm">
              In progress
            </div>
          </div>
          <Progress :model-value="0" class="bg-gradient-to-l from-black/60 to-black animate-pulse"/>
        </div>
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
        <div v-if="fileTotalBytes <= serverConfig.file_size_limit">
          {{ formatData(fileTotalBytes)}} in total
        </div>
        <div v-else class="text-red-700">
          {{ formatData(fileTotalBytes)}} in total
        </div>
      </div>
      <div v-for="(file, index) in files" class="border rounded-lg p-2 flex justify-between gap-2">
        <div class="flex gap-2 truncate">
          <Icon name="lucide:file"  class="text-4xl my-auto"/>
          <div class="truncate">
            <div class="font-bold my-auto truncate">
              {{ file.name }}
            </div>
            <div v-if="file.size <= serverConfig.file_size_limit" class="text-muted-foreground text-sm">
              {{ formatData(file.size) }}
            </div>
            <div v-else class="text-sm text-red-400">
              {{ formatData(file.size) }}
            </div>
          </div>
        </div>
        <Button v-if="!filesIsUploading" variant="ghost" class="my-auto" @click="filesRm(index)">
          <Icon name="lucide:x" />
        </Button>
      </div>
    </CardContent>
    <CardFooter>
      <div v-if="isZipping" class="flex justify-end w-full">
        <Button @click="cancel" class="right-0">Cancel</Button>
      </div>
      <div v-else-if="filesIsUploading" class="flex justify-between w-full space-x-2">
        <div v-if="!isZipping">
          <Button v-if="isPaused"  @click="upload">Resume</Button>
          <Button v-else  @click="pause">Pause</Button>
        </div>
        <Button @click="cancel" class="right-0">Cancel</Button>
      </div>
      <div v-else class="flex justify-end w-full">
        <Button @click="upload">Upload</Button>
      </div>
    </CardFooter>
  </Card>
</template>
