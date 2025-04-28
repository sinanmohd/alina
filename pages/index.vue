<script setup lang="ts">
import Textarea from '~/components/ui/textarea/Textarea.vue';
</script>

<template>
  <Tabs default-value="files" class="w-full sm:w-3xl mx-auto p-4">

    <TabsList class="grid w-full grid-cols-3">
      <TabsTrigger value="files">
        <Icon name="mdi:files"/>
        Files
      </TabsTrigger>
      <TabsTrigger value="text">
        <Icon name="mdi:text"/>
        Text
      </TabsTrigger>
      <TabsTrigger value="terminal">
        <Icon name="mdi:terminal"/>
        Terminal
      </TabsTrigger>
    </TabsList>

    <TabsContent value="files">
      <Card>
        <CardHeader>
          <CardTitle>Files</CardTitle>
          <CardDescription>
            Your frenly neighbourhood file sharing web thing.
          </CardDescription>
        </CardHeader>
        <CardContent class="space-y-2">
          <div @click="clickUploadInput" class="border-2 border-dashed p-16 rounded-lg">
            <div class="w-min mx-auto">
              <Icon name="mdi:cloud-upload-outline" class="text-7xl text-muted-foreground"/>
            </div>
            <p class="text-sm text-muted-foreground text-center">Drag & drop files here, or click to select files</p>
          </div>
          <input ref="uploadInput" type="file" multiple="true" class="invisible" >
        </CardContent>
        <CardFooter class="flex justify-end">
          <Button>Upload</Button>
        </CardFooter>
      </Card>
    </TabsContent>

    <TabsContent value="text">
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
    </TabsContent>

    <TabsContent value="terminal">
      <Card>
        <CardHeader>
          <CardTitle>Terminal</CardTitle>
          <CardDescription>
            Upload your files from the terminal using curl, the link will be printed to stdout after completion.
          </CardDescription>
        </CardHeader>

        <CardContent class="space-y-2">
        <code class="text-sm inline-flex bg-gray-50 rounded-lg p-2">
            <span class="flex gap-4">
                <span class="shrink-0 text-gray-500 my-auto">
                    $
                </span>

                <span class="flex-1 wrap-anywhere my-auto">
                  curl -F "file=@/path/to/your/file.jpg" https://share.sinanmohd.com/
                </span>
                <Button variant="outline" class="my-auto">
                  <Icon name="mdi:content-copy"/>
                </Button>
            </span>
        </code>
        </CardContent>
      </Card>
    </TabsContent>
  </Tabs>
</template>

<script lang="ts">
export default {
  methods: {
    clickUploadInput() {
      this.$refs.uploadInput.click()
    }
  }
}
</script>
