<script setup lang="ts">
import { toast } from 'vue-sonner';

const fileUploadDialog = useState<boolean>('fileUploadDialog', () => false)
const uploadLink = useState<string>('uploadLink');

async function dialogLinkToClipBoard() {
  if (!navigator.clipboard) {
      toast("Clipboard Access Failed", {
        description: "For security reasons clipboard is disabled",
      });

      return;
  }

  fileUploadDialog.value = false;
  navigator.clipboard.writeText(uploadLink.value)
}

async function dialogCancel() {
  fileUploadDialog.value = false;
}
</script>

<template>
  <AlertDialog v-bind:open="fileUploadDialog">
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>Data Uploaded Successfully</AlertDialogTitle>
        <AlertDialogDescription>
          Your data have been uploaded successfully and is ready to be shared or downloaded as
          <a :href="uploadLink" target="_blank" class="text-black underline">{{uploadLink}}</a>
        </AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel :onclick="dialogCancel">Cancel</AlertDialogCancel>
        <AlertDialogAction :onclick="dialogLinkToClipBoard">Copy Link</AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>

  <Tabs default-value="files" class="w-full sm:w-3xl mx-auto p-4">

    <TabsList class="grid w-full grid-cols-2">
      <TabsTrigger value="files">
        <Icon name="lucide:files"/>
        Files
      </TabsTrigger>
      <TabsTrigger value="text">
        <Icon name="material-symbols:markdown-outline"/>
        Text
      </TabsTrigger>
    </TabsList>

    <TabsContent value="files">
      <UploadFiles />
    </TabsContent>

    <TabsContent value="text">
      <UploadText />
    </TabsContent>
  </Tabs>
</template>
