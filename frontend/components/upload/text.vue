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
const serverUrl = useState('serverUrl', () => window.location.origin)

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
  const blob = new Blob([textContentTrimmed], { type: 'text/plain'});
  body.append("file", blob, 'note.text');

  await $fetch(`${serverUrl.value}/_alina/uplcvoad/simple`, {
    method: "POST",
    body: body,
    server: false,
    headers: {"cache-control": "no-cache"},
  }).then((d) => {
    if (textMarkdownSwitch.value) {
      const  data =  d as string
      const fileId = data.split('/').slice(-1)[0].replace('.txt', '')
      uploadLink.value = `${serverUrl.value}/notes/${fileId}`
    } else {
      uploadLink.value = d as string;
    }

    fileUploadDialog.value = true;
  }).catch(() => {
    toast("Upload Failed", {
      description: "Failed upload notes, try again later",
    });
  }).finally(() => {
    textIsUploading.value = false;
  })
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
