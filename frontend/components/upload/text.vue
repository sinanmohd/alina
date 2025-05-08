<script setup lang="ts">
import { toast } from 'vue-sonner';
import Textarea from '~/components/ui/textarea/Textarea.vue';

const { formatBytes } = useUtils();
const textContent = useState<string>('textContent', () => "")
const uploadLink = useState<string>('uploadLink');
const fileUploadDialog = useState<boolean>('fileUploadDialog');
const textMarkdownSwitch = useState<boolean>('textMarkdownSwitch', () => true);
const textIsUploading = useState('textIsUploading', () => false);
const filesIsUploading = useState('filesIsUploading', () => false);
const serverConfig = useServerConfig();
const appConfig = useAppConfig();

async function upload() {
  const textContentTrimmed = textContent.value.trimEnd();

  if (textContentTrimmed.length === 0) {
    toast('Empty Text', {
      description: 'Please write something to proceed',
    });

    return
  } else if (filesIsUploading.value || textIsUploading.value) {
    toast('Upload in Progress', {
      description: 'Please wait until the current upload is complete',
    });

    return
  } else if (textContentTrimmed.length > serverConfig.value.file_size_limit) {
    toast('Files Too Large', {
      description: `These files exceeds the upload size limit of ${formatBytes(serverConfig.value.file_size_limit)}`,
    });

    return
  }
  const body = new FormData();
  body.append("file", new Blob([textContentTrimmed], { type: 'text/plain'}))

  const req = await useFetch(`${appConfig.serverUrl}/_alina/upload/simple`, {
    method: "POST",
    body: body
  })
  if (req.status.value == "error" || !req.data.value) {
    toast("Upload Failed", {
      description: req.error.value?.message,
    });

    textIsUploading.value = false;
    return
  }

  if (textMarkdownSwitch.value) {
    const  data =  req.data.value as string
    const fileId = data.split('/').slice(-1)[0].replace('.txt', '')
    uploadLink.value = `${appConfig.serverUrl}/notes/${fileId}`
  } else {
    uploadLink.value = req.data.value as string;
  }

  fileUploadDialog.value = true;
}
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
      <Textarea class="h-56" v-model="textContent"/>
    </CardContent>
    <CardFooter class="flex justify-between">
      <div class="flex items-center space-x-2">
        <Switch v-model:model-value="textMarkdownSwitch"/>
        <Label>Markdown</Label>
      </div>

      <Button v-if="!textIsUploading" :onclick="upload">Save</Button>
      <Button v-else :onclick="upload">
        <Icon name="svg-spinners:gooey-balls-1"/>
        Saving
      </Button>
    </CardFooter>
  </Card>
</template>
